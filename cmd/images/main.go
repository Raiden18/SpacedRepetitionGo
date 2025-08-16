package main

import (
	"spacedrepetitiongo/jobs"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	jobs.ReplaceImages()
}
