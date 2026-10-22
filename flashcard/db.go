package flashcard

import (
	"context"
	"log"
	"sort"
	"spacedrepetitiongo/image"
	"spacedrepetitiongo/utils"
	"strings"

	"github.com/jmoiron/sqlx"
)

// All cards live in one table. A card's queue says what the user does with
// it next; position orders the cards of one box within a queue.
const (
	CARDS_TABLE = "cards"

	QueueRevise   = "revise"
	QueueMemorize = "memorize"
)

// Tables used before the merge into CARDS_TABLE. They are only read, once, by
// migrateLegacyTables, and are left in place so an older release still works
// after a rollback.
const (
	legacyReviseTable   = "flash_cards_to_revise"
	legacyMemorizeTable = "flash_cards_to_memorize"
	migrationLockName   = "spacedrepetitiongo_cards_migration"
)

type FlashcardEntity struct {
	Id          string  `db:"page_id"`
	Queue       string  `db:"queue"`
	Position    int     `db:"position"`
	BoxId       string  `db:"notion_data_base_id"`
	Image       *string `db:"image_url"`
	ImageSource *string `db:"image_source"`
	Name        string  `db:"name"`
	Example     *string `db:"example"`
	Explanation *string `db:"answer"`
	KnowLevel1  *bool   `db:"know_level_1"`
	KnowLevel2  *bool   `db:"know_level_2"`
	KnowLevel3  *bool   `db:"know_level_3"`
	KnowLevel4  *bool   `db:"know_level_4"`
	KnowLevel5  *bool   `db:"know_level_5"`
	KnowLevel6  *bool   `db:"know_level_6"`
	KnowLevel7  *bool   `db:"know_level_7"`
	KnowLevel8  *bool   `db:"know_level_8"`
	KnowLevel9  *bool   `db:"know_level_9"`
	KnowLevel10 *bool   `db:"know_level_10"`
	KnowLevel11 *bool   `db:"know_level_11"`
	KnowLevel12 *bool   `db:"know_level_12"`
	KnowLevel13 *bool   `db:"know_level_13"`
	KnowLevel14 *bool   `db:"know_level_14"`
}

func toEntity(flashcard Flashcard) FlashcardEntity {
	return FlashcardEntity{
		Id:          flashcard.Id,
		Queue:       flashcard.Queue,
		Position:    flashcard.Position,
		BoxId:       flashcard.BoxId,
		Image:       flashcard.Image,
		ImageSource: flashcard.ImageSource,
		Name:        flashcard.Name,
		Example:     flashcard.Example,
		Explanation: flashcard.Explanation,
		KnowLevel1:  flashcard.KnowLevels[1],
		KnowLevel2:  flashcard.KnowLevels[2],
		KnowLevel3:  flashcard.KnowLevels[3],
		KnowLevel4:  flashcard.KnowLevels[4],
		KnowLevel5:  flashcard.KnowLevels[5],
		KnowLevel6:  flashcard.KnowLevels[6],
		KnowLevel7:  flashcard.KnowLevels[7],
		KnowLevel8:  flashcard.KnowLevels[8],
		KnowLevel9:  flashcard.KnowLevels[9],
		KnowLevel10: flashcard.KnowLevels[10],
		KnowLevel11: flashcard.KnowLevels[11],
		KnowLevel12: flashcard.KnowLevels[12],
		KnowLevel13: flashcard.KnowLevels[13],
		KnowLevel14: flashcard.KnowLevels[14],
	}
}

func toFlashcard(entity FlashcardEntity) Flashcard {
	return Flashcard{
		Id:          entity.Id,
		Queue:       entity.Queue,
		Position:    entity.Position,
		BoxId:       entity.BoxId,
		Image:       entity.Image,
		ImageSource: entity.ImageSource,
		Name:        entity.Name,
		Example:     entity.Example,
		Explanation: entity.Explanation,
		KnowLevels: map[int]*bool{
			1:  entity.KnowLevel1,
			2:  entity.KnowLevel2,
			3:  entity.KnowLevel3,
			4:  entity.KnowLevel4,
			5:  entity.KnowLevel5,
			6:  entity.KnowLevel6,
			7:  entity.KnowLevel7,
			8:  entity.KnowLevel8,
			9:  entity.KnowLevel9,
			10: entity.KnowLevel10,
			11: entity.KnowLevel11,
			12: entity.KnowLevel12,
			13: entity.KnowLevel13,
			14: entity.KnowLevel14,
		},
	}
}

// CreateTableAndMigrate creates CARDS_TABLE. When the table didn't exist yet
// (the first start after the merge), it also copies the cards from the old
// per-queue tables.
func CreateTableAndMigrate(db sqlx.DB) {
	existed := tableExists(db, CARDS_TABLE)
	createTable(db)
	if !existed {
		migrateLegacyTables(db)
	}
}

func tableExists(db sqlx.DB, table string) bool {
	var count int
	err := db.Get(&count, `SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = DATABASE() AND table_name = ?`, table)
	if err != nil {
		log.Panicln("Could not check for table "+table, err)
	}
	return count > 0
}

func createTable(db sqlx.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + CARDS_TABLE + ` (
			page_id VARCHAR(255) NOT NULL PRIMARY KEY,
			queue VARCHAR(16) NOT NULL,
			position INT NOT NULL,
			notion_data_base_id VARCHAR(255) NOT NULL,
			image_url TEXT,
			image_source TEXT,
			name TEXT,
			example TEXT,
			answer TEXT,
			know_level_1 BOOLEAN,
			know_level_2 BOOLEAN,
			know_level_3 BOOLEAN,
			know_level_4 BOOLEAN,
			know_level_5 BOOLEAN,
			know_level_6 BOOLEAN,
			know_level_7 BOOLEAN,
			know_level_8 BOOLEAN,
			know_level_9 BOOLEAN,
			know_level_10 BOOLEAN,
			know_level_11 BOOLEAN,
			know_level_12 BOOLEAN,
			know_level_13 BOOLEAN,
			know_level_14 BOOLEAN,
			INDEX idx_box_queue_position (notion_data_base_id, queue, position)
		);
	`)
	if err != nil {
		log.Panic("Could not create table "+CARDS_TABLE, err)
	}
}

// ---------- Reading ----------

// LoadAll returns every card, ordered by box and position.
func LoadAll(db sqlx.DB) []Flashcard {
	return selectCards(db, `SELECT * FROM `+CARDS_TABLE+` ORDER BY notion_data_base_id, queue, position`)
}

// LoadQueue returns the cards of one queue, ordered by box and position.
func LoadQueue(db sqlx.DB, queue string) []Flashcard {
	return selectCards(db, `SELECT * FROM `+CARDS_TABLE+` WHERE queue = ? ORDER BY notion_data_base_id, position`, queue)
}

// LoadById returns the card with its Previous/Next neighbours filled in, or
// nil if there is no such card (e.g. it was already answered).
func LoadById(db sqlx.DB, id string) *Flashcard {
	return getCard(db, `SELECT * FROM `+CARDS_TABLE+` WHERE page_id = ?`, id)
}

// FirstInBox returns the first card of a box's queue, or nil if it is empty.
func FirstInBox(db sqlx.DB, boxId string, queue string) *Flashcard {
	return getCard(db, `SELECT * FROM `+CARDS_TABLE+`
		WHERE notion_data_base_id = ? AND queue = ? ORDER BY position LIMIT 1`, boxId, queue)
}

// LastInBox returns the last card of a box's queue, or nil if it is empty.
func LastInBox(db sqlx.DB, boxId string, queue string) *Flashcard {
	return getCard(db, `SELECT * FROM `+CARDS_TABLE+`
		WHERE notion_data_base_id = ? AND queue = ? ORDER BY position DESC LIMIT 1`, boxId, queue)
}

func selectCards(db sqlx.DB, query string, args ...any) []Flashcard {
	entities := []FlashcardEntity{}
	if err := db.Select(&entities, query, args...); err != nil {
		log.Panicln("Could not read cards", err)
	}
	return utils.Map(entities, toFlashcard)
}

func getCard(db sqlx.DB, query string, args ...any) *Flashcard {
	entities := []FlashcardEntity{}
	if err := db.Select(&entities, query, args...); err != nil {
		log.Panicln("Could not read card", err)
	}
	if len(entities) == 0 {
		return nil
	}
	card := toFlashcard(entities[0])
	card.Previous = neighbourId(db, card, `position < ? ORDER BY position DESC`)
	card.Next = neighbourId(db, card, `position > ? ORDER BY position`)
	return &card
}

func neighbourId(db sqlx.DB, card Flashcard, condition string) *string {
	ids := []string{}
	err := db.Select(&ids, `SELECT page_id FROM `+CARDS_TABLE+`
		WHERE notion_data_base_id = ? AND queue = ? AND `+condition+` LIMIT 1`,
		card.BoxId, card.Queue, card.Position)
	if err != nil {
		log.Panicln("Could not read neighbouring card", err)
	}
	if len(ids) == 0 {
		return nil
	}
	return &ids[0]
}

// ---------- Answer handling ----------

// Delete removes the card from the local cache.
func (flashcard Flashcard) Delete(db sqlx.DB) Flashcard {
	if _, err := db.Exec(`DELETE FROM `+CARDS_TABLE+` WHERE page_id = ?`, flashcard.Id); err != nil {
		log.Panicln("Could not delete card "+flashcard.Id, err)
	}
	return flashcard
}

// MoveToEndOf moves the card to the end of its box's queue in one transaction.
func (flashcard Flashcard) MoveToEndOf(db sqlx.DB, queue string) Flashcard {
	tx, err := db.Beginx()
	if err != nil {
		log.Panicln("Could not start transaction", err)
	}
	defer tx.Rollback()

	var maxPosition *int
	err = tx.Get(&maxPosition, `SELECT MAX(position) FROM `+CARDS_TABLE+`
		WHERE notion_data_base_id = ? AND queue = ?`, flashcard.BoxId, queue)
	if err != nil {
		log.Panicln("Could not read last position", err)
	}
	position := 0
	if maxPosition != nil {
		position = *maxPosition + 1
	}

	_, err = tx.Exec(`UPDATE `+CARDS_TABLE+` SET queue = ?, position = ? WHERE page_id = ?`,
		queue, position, flashcard.Id)
	if err != nil {
		log.Panicln("Could not move card "+flashcard.Id, err)
	}
	if err := tx.Commit(); err != nil {
		log.Panicln("Could not commit card move", err)
	}

	flashcard.Queue = queue
	flashcard.Position = position
	return flashcard
}

// ---------- Sync from Notion ----------

// Sync makes the cache match Notion in one transaction. Cards are given
// positions in the order passed, per box. Cards in keepAsIs (answers not yet
// in Notion) are neither updated nor deleted.
func Sync(db sqlx.DB, toRevise []Flashcard, toMemorize []Flashcard, keepAsIs map[string]struct{}) {
	cards := append(withPositions(toRevise, QueueRevise), withPositions(toMemorize, QueueMemorize)...)

	keepIds := make(map[string]struct{}, len(cards)+len(keepAsIs))
	for _, card := range cards {
		keepIds[card.Id] = struct{}{}
	}
	for id := range keepAsIs {
		keepIds[id] = struct{}{}
	}

	tx, err := db.Beginx()
	if err != nil {
		log.Panicln("Could not start transaction", err)
	}
	defer tx.Rollback()

	existingIds := []string{}
	if err := tx.Select(&existingIds, `SELECT page_id FROM `+CARDS_TABLE); err != nil {
		log.Panicln("Could not read card ids", err)
	}
	for _, id := range existingIds {
		if _, keep := keepIds[id]; keep {
			continue
		}
		if _, err := tx.Exec(`DELETE FROM `+CARDS_TABLE+` WHERE page_id = ?`, id); err != nil {
			log.Panicln("Could not delete card "+id, err)
		}
	}

	upsert(tx, cards)

	if err := tx.Commit(); err != nil {
		log.Panicln("Could not commit sync", err)
	}
}

func withPositions(cards []Flashcard, queue string) []Flashcard {
	nextPosition := map[string]int{}
	result := make([]Flashcard, len(cards))
	for i, card := range cards {
		card.Queue = queue
		card.Position = nextPosition[card.BoxId]
		nextPosition[card.BoxId]++
		result[i] = card
	}
	return result
}

func upsert(tx *sqlx.Tx, cards []Flashcard) {
	if len(cards) == 0 {
		return
	}
	_, err := tx.NamedExec(`INSERT INTO `+CARDS_TABLE+` (
		page_id, queue, position, notion_data_base_id, image_url, image_source,
		name, example, answer,
		know_level_1, know_level_2, know_level_3, know_level_4, know_level_5,
		know_level_6, know_level_7, know_level_8, know_level_9, know_level_10,
		know_level_11, know_level_12, know_level_13, know_level_14
	) VALUES (
		:page_id, :queue, :position, :notion_data_base_id, :image_url, :image_source,
		:name, :example, :answer,
		:know_level_1, :know_level_2, :know_level_3, :know_level_4, :know_level_5,
		:know_level_6, :know_level_7, :know_level_8, :know_level_9, :know_level_10,
		:know_level_11, :know_level_12, :know_level_13, :know_level_14
	) ON DUPLICATE KEY UPDATE
		queue=VALUES(queue),
		position=VALUES(position),
		notion_data_base_id=VALUES(notion_data_base_id),
		image_url=VALUES(image_url),
		image_source=VALUES(image_source),
		name=VALUES(name),
		example=VALUES(example),
		answer=VALUES(answer),
		know_level_1=VALUES(know_level_1),
		know_level_2=VALUES(know_level_2),
		know_level_3=VALUES(know_level_3),
		know_level_4=VALUES(know_level_4),
		know_level_5=VALUES(know_level_5),
		know_level_6=VALUES(know_level_6),
		know_level_7=VALUES(know_level_7),
		know_level_8=VALUES(know_level_8),
		know_level_9=VALUES(know_level_9),
		know_level_10=VALUES(know_level_10),
		know_level_11=VALUES(know_level_11),
		know_level_12=VALUES(know_level_12),
		know_level_13=VALUES(know_level_13),
		know_level_14=VALUES(know_level_14);`,
		utils.Map(cards, toEntity),
	)
	if err != nil {
		log.Panicln("Could not save cards", err)
	}
}

// ---------- Images ----------

func UpdateImage(db sqlx.DB, id string, newImage string) {
	if _, err := db.Exec(`UPDATE `+CARDS_TABLE+` SET image_url = ? WHERE page_id = ?`, newImage, id); err != nil {
		log.Panicln("Could not update image of card "+id, err)
	}
}

type storedImage struct {
	Id     string  `db:"page_id"`
	Image  *string `db:"image_url"`
	Source *string `db:"image_source"`
}

// StoredImages returns the image columns of every card, by page id.
func StoredImages(db sqlx.DB) map[string]storedImage {
	rows := []storedImage{}
	if err := db.Select(&rows, `SELECT page_id, image_url, image_source FROM `+CARDS_TABLE); err != nil {
		log.Panicln("Could not read stored images", err)
	}
	byId := make(map[string]storedImage, len(rows))
	for _, row := range rows {
		byId[row.Id] = row
	}
	return byId
}

// KeepDownloadedImages keeps each card's already-downloaded local image when
// its cover in Notion hasn't changed. Otherwise (new card, changed cover, or
// not downloaded yet) the card gets the fresh remote URL, which the bot can
// send directly until the image job downloads it. Notion's file URLs expire,
// so a stale remote URL is never kept.
func KeepDownloadedImages(flashCards []Flashcard, stored map[string]storedImage) []Flashcard {
	result := make([]Flashcard, len(flashCards))
	for i, card := range flashCards {
		if previous, ok := stored[card.Id]; ok &&
			previous.Image != nil && strings.HasPrefix(*previous.Image, "/") &&
			previous.Source != nil && card.ImageSource != nil &&
			image.SameSource(*previous.Source, *card.ImageSource) {
			card.Image = previous.Image
		} else {
			card.Image = card.ImageSource
		}
		result[i] = card
	}
	return result
}

// ---------- One-time migration from the old per-queue tables ----------

type legacyEntity struct {
	Id          string  `db:"page_id"`
	Image       *string `db:"image_url"`
	ImageSource *string `db:"image_source"`
	BoxId       string  `db:"notion_data_base_id"`
	Name        string  `db:"name"`
	Example     *string `db:"example"`
	Explanation *string `db:"answer"`
	Next        *string `db:"next"`
	Previous    *string `db:"previous"`
	KnowLevel1  *bool   `db:"know_level_1"`
	KnowLevel2  *bool   `db:"know_level_2"`
	KnowLevel3  *bool   `db:"know_level_3"`
	KnowLevel4  *bool   `db:"know_level_4"`
	KnowLevel5  *bool   `db:"know_level_5"`
	KnowLevel6  *bool   `db:"know_level_6"`
	KnowLevel7  *bool   `db:"know_level_7"`
	KnowLevel8  *bool   `db:"know_level_8"`
	KnowLevel9  *bool   `db:"know_level_9"`
	KnowLevel10 *bool   `db:"know_level_10"`
	KnowLevel11 *bool   `db:"know_level_11"`
	KnowLevel12 *bool   `db:"know_level_12"`
	KnowLevel13 *bool   `db:"know_level_13"`
	KnowLevel14 *bool   `db:"know_level_14"`
}

func (e legacyEntity) toFlashcard() Flashcard {
	return Flashcard{
		Id: e.Id, Image: e.Image, ImageSource: e.ImageSource, BoxId: e.BoxId,
		Name: e.Name, Example: e.Example, Explanation: e.Explanation,
		KnowLevels: map[int]*bool{
			1: e.KnowLevel1, 2: e.KnowLevel2, 3: e.KnowLevel3, 4: e.KnowLevel4,
			5: e.KnowLevel5, 6: e.KnowLevel6, 7: e.KnowLevel7, 8: e.KnowLevel8,
			9: e.KnowLevel9, 10: e.KnowLevel10, 11: e.KnowLevel11, 12: e.KnowLevel12,
			13: e.KnowLevel13, 14: e.KnowLevel14,
		},
	}
}

// migrateLegacyTables copies cards from the old flash_cards_to_revise and
// flash_cards_to_memorize tables into CARDS_TABLE, keeping their order, so the
// bot has cards to show right after the first deploy instead of waiting for
// the next sync. If another process already copied them, it does nothing.
func migrateLegacyTables(db sqlx.DB) {
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Panicln("Could not get DB connection for migration", err)
	}
	defer conn.Close()

	// Every binary runs this on start; make sure only one copies at a time.
	var locked *int
	if err := conn.QueryRowContext(ctx, `SELECT GET_LOCK(?, 60)`, migrationLockName).Scan(&locked); err != nil || locked == nil || *locked != 1 {
		log.Panicln("Could not take migration lock", err)
	}
	defer conn.ExecContext(ctx, `SELECT RELEASE_LOCK(?)`, migrationLockName)

	var count int
	if err := db.Get(&count, `SELECT COUNT(*) FROM `+CARDS_TABLE); err != nil {
		log.Panicln("Could not count cards", err)
	}
	if count > 0 {
		return
	}

	toRevise := loadLegacyQueue(db, legacyReviseTable)
	toMemorize := loadLegacyQueue(db, legacyMemorizeTable)
	if len(toRevise)+len(toMemorize) == 0 {
		return
	}

	tx, err := db.Beginx()
	if err != nil {
		log.Panicln("Could not start migration transaction", err)
	}
	defer tx.Rollback()
	upsert(tx, append(withPositions(toRevise, QueueRevise), withPositions(toMemorize, QueueMemorize)...))
	if err := tx.Commit(); err != nil {
		log.Panicln("Could not commit migration", err)
	}
	log.Printf("Migrated %d cards to revise and %d to memorize into %s", len(toRevise), len(toMemorize), CARDS_TABLE)
}

// loadLegacyQueue returns a legacy table's cards in list order (following the
// next links from each box's head), or nil if the table doesn't exist.
func loadLegacyQueue(db sqlx.DB, table string) []Flashcard {
	if !tableExists(db, table) {
		return nil
	}

	rows := []legacyEntity{}
	if err := db.Unsafe().Select(&rows, `SELECT * FROM `+table); err != nil {
		log.Panicln("Could not read legacy table "+table, err)
	}
	return orderLegacy(rows)
}

func orderLegacy(rows []legacyEntity) []Flashcard {
	byId := map[string]legacyEntity{}
	boxes := []string{}
	seenBox := map[string]bool{}
	for _, row := range rows {
		byId[row.Id] = row
		if !seenBox[row.BoxId] {
			seenBox[row.BoxId] = true
			boxes = append(boxes, row.BoxId)
		}
	}
	sort.Strings(boxes)

	ordered := []Flashcard{}
	visited := map[string]bool{}
	for _, box := range boxes {
		boxRows := []legacyEntity{}
		for _, row := range rows {
			if row.BoxId == box {
				boxRows = append(boxRows, row)
			}
		}
		sort.Slice(boxRows, func(i, j int) bool { return boxRows[i].Id < boxRows[j].Id })

		// Walk from each head (no previous card in this box), then add any
		// cards the links didn't reach.
		walk := func(start legacyEntity) {
			for current, ok := start, true; ok && !visited[current.Id]; {
				visited[current.Id] = true
				ordered = append(ordered, current.toFlashcard())
				if current.Next == nil {
					break
				}
				current, ok = byId[*current.Next]
				ok = ok && current.BoxId == box
			}
		}
		for _, row := range boxRows {
			previous, hasPrevious := legacyEntity{}, false
			if row.Previous != nil {
				previous, hasPrevious = byId[*row.Previous]
			}
			if !hasPrevious || previous.BoxId != box {
				walk(row)
			}
		}
		for _, row := range boxRows {
			walk(row)
		}
	}
	return ordered
}
