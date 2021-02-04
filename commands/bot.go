package commands

import (
	"fmt"
	"gbot/inspiration"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

type Bot struct {
	Ctx *bot.Context
}

func (bot *Bot) Setup(sub *bot.Subcommand) {
}

// Help prints the default help message.
func (bot *Bot) Help(*gateway.MessageCreateEvent) (string, error) {
	return bot.Ctx.Help(), nil
}

// Add demonstrates the usage of typed arguments. Run it with "~add 1 2".
func (bot *Bot) Add(_ *gateway.MessageCreateEvent, a, b int) (string, error) {
	return fmt.Sprintf("%d + %d = %d", a, b, a+b), nil
}

func (bot *Bot) Ping(*gateway.MessageCreateEvent) (string, error) {
	return "Pong!", nil
}

func (bot *Bot) Inspire(m *gateway.MessageCreateEvent) (string, error) {
	resp, err := inspiration.NewInspirationalMessage()
	return resp, err
}
