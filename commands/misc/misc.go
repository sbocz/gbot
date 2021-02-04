package misc

import (
	"gbot/database"
)

func Initialize(db *database.DB) {
	initYeller(database.NewBucket(db, database.ShoutPhrases))
}
