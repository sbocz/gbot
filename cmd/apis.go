package cmd

import (
	"gbot/apis"

	"github.com/diamondburned/arikawa/v2/gateway"
)

func (bot *Bot) Inspire(m *gateway.MessageCreateEvent) (string, error) {
	resp, err := apis.GetInspirobotMessage()
	return resp, err
}
