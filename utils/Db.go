package utils

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func OpenDb() sqlx.DB {
	db, err := sqlx.Connect("mysql", "root@unix(/var/run/mysqld/mysqld.sock)/spaced_repetition")
	if err != nil {
		log.Fatal(err)
	}

	return *db
}
