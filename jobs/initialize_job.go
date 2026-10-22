package jobs

import (
	"spacedrepetitiongo/schema"
	"spacedrepetitiongo/utils"
)

func Initialize() {
	db := utils.OpenDb()
	defer db.Close()
	schema.Ensure(db)
}
