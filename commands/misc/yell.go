package misc

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

var bucket *database.Bucket

type YelledMessage struct {
	Author  discord.UserID
	Message string
}

func initYeller(b *database.Bucket) {
	bucket = b
}

func YellHandler(ctx *bot.Context) func(*gateway.MessageCreateEvent) {
	return func(e *gateway.MessageCreateEvent) {
		m, err := ctx.State.Message(e.ChannelID, e.ID)
		if err != nil {
			log.Println("Not found:", e.ID)
			return
		}
		if !m.Author.Bot && strings.ToUpper(m.Content) == m.Content {
			bucket.PutRandom(&YelledMessage{Author: m.Author.ID, Message: m.Content})
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
