package cmd

import (
	"encoding/json"
	"fmt"
	"gbot/database"
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

const censor = "****"

var bucket *database.Bucket
var denyMap map[string]bool

type YelledMessage struct {
	Author  discord.UserID
	Message string
}

func Initialize(db *database.DB, denyList []string) error {
	bucket, err := database.NewBucket(db, database.ShoutPhrases)
	if err != nil {
		return fmt.Errorf("could not initialize misc commands: %v", err)
	}

	initYeller(bucket, denyList)
	return nil
}

func initYeller(b *database.Bucket, denyList []string) {
	bucket = b
	denyMap = make(map[string]bool)
	for _, deniedWord := range denyList {
		denyMap[strings.ToUpper(deniedWord)] = true
	}
}

func YellHandler(ctx *bot.Context) func(*gateway.MessageCreateEvent) {
	return func(e *gateway.MessageCreateEvent) {
		m, err := ctx.State.Message(e.ChannelID, e.ID)
		if err != nil {
			log.Println("Not found:", e.ID)
			return
		}
		if !m.Author.Bot && strings.ToUpper(m.Content) == m.Content {
			toSave := censoredString(m.Content)
			bucket.PutRandom(&YelledMessage{Author: m.Author.ID, Message: toSave})
			rawValue, _ := bucket.GetRandom()
			var response *YelledMessage
			err = json.Unmarshal(rawValue, &response)
			if err != nil {
				log.Println(fmt.Sprintf("Could not unmarshall response %s: %s", rawValue, err))
				return
			}
			ctx.SendMessage(m.ChannelID, response.Message, nil)
		}
	}
}

func censoredString(toCensor string) string {
	contentWords := strings.Split(toCensor, " ")
	for i, word := range contentWords {
		if denyMap[word] {
			contentWords[i] = censor
		}
	}
	return strings.Join(contentWords, " ")
}
