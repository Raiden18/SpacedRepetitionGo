// Package outbox makes write-back of answers to Notion reliable.
//
// Every answer is first stored as a row in notion_outbox, then sent to Notion
// by Drain, which retries with backoff until Notion accepts it. While a page
// still has a pending row, the updater must not overwrite the local cache
// with Notion's (stale) copy of that page — see PendingPageIds.
package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"spacedrepetitiongo/notion"

	"github.com/jmoiron/sqlx"
	notionApi "github.com/jomei/notionapi"
)

const (
	TABLE_NAME = "notion_outbox"

	drainLockName = "spacedrepetitiongo_notion_outbox_drain"

	// Notion allows ~3 requests/second per integration.
	delayBetweenRequests = 350 * time.Millisecond

	baseBackoff = 10 * time.Second
	maxBackoff  = time.Hour

	// After this many failed attempts a row is reported loudly on every pass.
	alertAfterAttempts = 10

	workerInterval = 30 * time.Second
)

// Update is the set of Notion properties to write to one page.
type Update struct {
	Checkboxes map[string]bool    `json:"checkboxes,omitempty"`
	Numbers    map[string]float64 `json:"numbers,omitempty"`
}

func (update Update) toRequest() notionApi.PageUpdateRequest {
	properties := notionApi.Properties{}
	for name, value := range update.Checkboxes {
		properties[name] = notionApi.CheckboxProperty{Checkbox: value}
	}
	for name, value := range update.Numbers {
		properties[name] = notionApi.NumberProperty{Number: value}
	}
	return notion.NewUpdateRequest(properties)
}

type row struct {
	Id         int64  `db:"id"`
	PageId     string `db:"page_id"`
	Properties string `db:"properties"`
	Attempts   int    `db:"attempts"`
	Due        bool   `db:"due"`
}

func CreateTableIfNotExist(db sqlx.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + TABLE_NAME + ` (
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			page_id VARCHAR(255) NOT NULL,
			properties TEXT NOT NULL,
			attempts INT NOT NULL DEFAULT 0,
			last_error TEXT,
			next_attempt_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_page_id (page_id)
		);
	`)
	if err != nil {
		log.Panic("Could not create table "+TABLE_NAME, err)
	}
}

// Enqueue durably records that pageId must be updated in Notion and wakes the
// worker. Call it before changing the local cache, so an answer is never lost.
func Enqueue(db sqlx.DB, pageId string, update Update) error {
	properties, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("encode outbox update for %s: %w", pageId, err)
	}
	_, err = db.Exec(
		`INSERT INTO `+TABLE_NAME+` (page_id, properties) VALUES (?, ?)`,
		pageId, string(properties),
	)
	if err != nil {
		return fmt.Errorf("insert outbox row for %s: %w", pageId, err)
	}
	Kick()
	return nil
}

// PendingPageIds returns the pages whose answers have not reached Notion yet.
func PendingPageIds(db sqlx.DB) (map[string]struct{}, error) {
	ids := []string{}
	if err := db.Select(&ids, `SELECT DISTINCT page_id FROM `+TABLE_NAME); err != nil {
		return nil, fmt.Errorf("read pending outbox pages: %w", err)
	}
	pending := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		pending[id] = struct{}{}
	}
	return pending, nil
}

// Drain sends due outbox rows to Notion in insertion order. Rows for a page
// are applied strictly in order: once one fails or is waiting for its retry
// time, later rows for the same page wait too. Only one process drains at a
// time; if another holds the lock for longer than waitForLock, Drain returns
// without doing anything.
func Drain(db sqlx.DB, client notion.Client, waitForLock time.Duration) {
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Println("outbox: could not get DB connection:", err)
		return
	}
	defer conn.Close()

	var locked int
	if err := conn.QueryRowContext(ctx, `SELECT GET_LOCK(?, ?)`, drainLockName, int(waitForLock.Seconds())).Scan(&locked); err != nil {
		log.Println("outbox: could not take drain lock:", err)
		return
	}
	if locked != 1 {
		return
	}
	defer conn.ExecContext(ctx, `SELECT RELEASE_LOCK(?)`, drainLockName)

	rows := []row{}
	err = db.Select(&rows, `
		SELECT id, page_id, properties, attempts, next_attempt_at <= NOW() AS due
		FROM `+TABLE_NAME+`
		ORDER BY id`)
	if err != nil {
		log.Println("outbox: could not read rows:", err)
		return
	}

	blockedPages := map[string]bool{}
	sentAny := false
	for _, r := range rows {
		if blockedPages[r.PageId] {
			continue
		}
		if !r.Due {
			blockedPages[r.PageId] = true
			continue
		}

		var update Update
		if err := json.Unmarshal([]byte(r.Properties), &update); err != nil {
			// Undecodable rows can never succeed; drop them so they don't block the page forever.
			log.Printf("outbox: dropping undecodable row %d for page %s: %v", r.Id, r.PageId, err)
			deleteRow(db, r.Id)
			continue
		}

		if sentAny {
			time.Sleep(delayBetweenRequests)
		}
		sentAny = true

		request := update.toRequest()
		if err := client.UpdatePage(r.PageId, &request); err != nil {
			blockedPages[r.PageId] = true
			recordFailure(db, r, err)
			continue
		}
		deleteRow(db, r.Id)
	}
}

// StartWorker drains the outbox whenever Enqueue is called and every
// workerInterval, for the lifetime of the process.
func StartWorker(db sqlx.DB, client notion.Client) {
	go func() {
		for {
			Drain(db, client, 0)
			select {
			case <-kick:
			case <-time.After(workerInterval):
			}
		}
	}()
}

var kick = make(chan struct{}, 1)

// Kick wakes the worker, if one is running, without blocking.
func Kick() {
	select {
	case kick <- struct{}{}:
	default:
	}
}

func deleteRow(db sqlx.DB, id int64) {
	if _, err := db.Exec(`DELETE FROM `+TABLE_NAME+` WHERE id = ?`, id); err != nil {
		log.Printf("outbox: could not delete sent row %d: %v", id, err)
	}
}

func recordFailure(db sqlx.DB, r row, sendErr error) {
	attempts := r.Attempts + 1
	backoff := time.Duration(float64(baseBackoff) * math.Pow(2, float64(attempts-1)))
	if backoff > maxBackoff || backoff <= 0 {
		backoff = maxBackoff
	}

	if attempts >= alertAfterAttempts {
		log.Printf("outbox: ALERT page %s has failed %d times; answer not yet in Notion: %v", r.PageId, attempts, sendErr)
	} else {
		log.Printf("outbox: page %s failed (attempt %d), retrying in %s: %v", r.PageId, attempts, backoff, sendErr)
	}

	_, err := db.Exec(
		`UPDATE `+TABLE_NAME+`
		SET attempts = ?, last_error = ?, next_attempt_at = NOW() + INTERVAL ? SECOND
		WHERE id = ?`,
		attempts, sendErr.Error(), int(backoff.Seconds()), r.Id,
	)
	if err != nil {
		log.Printf("outbox: could not record failure for row %d: %v", r.Id, err)
	}
}
