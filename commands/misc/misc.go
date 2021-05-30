package misc

import (
	"fmt"
	"gbot/database"
)

func Initialize(db *database.DB, denyList []string) error {
	bucket, err := database.NewBucket(db, database.ShoutPhrases)
	if err != nil {
		return fmt.Errorf("could not initialize misc commands: %v", err)
	}

	initYeller(bucket, denyList)
	return nil
}
