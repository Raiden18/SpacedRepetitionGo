package bot

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var (
	BOT_STATE = "bot_state"
)

type State struct {
	Id string `db:"state"`
}

func IsAddGreekVocabularyState(db sqlx.DB) bool {
	currentState := GetCurrentState(db)
	return currentState.Id == NewAddGreekVocabularyState().Id
}

func NewAddGreekVocabularyState() State {
	return State{
		Id: "add_greek_vocabulary",
	}
}

func GetCurrentState(db sqlx.DB) State {
	query := "SELECT * FROM " + BOT_STATE + " LIMIT 1"
	var state State
	err := db.Get(&state, query)
	if err != nil {
		log.Fatal(err)
	}
	return state
}

func SaveState(db sqlx.DB, state State) {
	_, err := db.NamedExec(
		`INSERT INTO `+BOT_STATE+` (state) VALUES (:state);`,
		state,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTableIfNotExist(db sqlx.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS ` + BOT_STATE + ` (state VARCHAR(255) NOT NULL PRIMARY KEY);`)
	if err != nil {
		log.Fatal(err)
	}
}
