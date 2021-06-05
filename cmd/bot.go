package cmd

import (
	"fmt"
	"math/rand"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

type Bot struct {
	Ctx *bot.Context
}

func NewBot() *Bot {
	return &Bot{}
}

func (bot *Bot) Setup(sub *bot.Subcommand) {
	sub.ChangeCommandInfo("Add", "add", "adds some numbers")
	sub.ChangeCommandInfo("Ping", "ping", "check if bot is alive")
	sub.ChangeCommandInfo("Inspire", "inspire", "display an inspirational message")
	sub.ChangeCommandInfo("Choose", "choose", "choose from a list of options")
}

// Help prints the default help message.
func (bot *Bot) Help(*gateway.MessageCreateEvent) (string, error) {
	return bot.Ctx.Help(), nil
}

func (bot *Bot) Add(_ *gateway.MessageCreateEvent, a, b int) (string, error) {
	return fmt.Sprintf("%d + %d = %d", a, b, a+b), nil
}

func (bot *Bot) Ping(*gateway.MessageCreateEvent) (string, error) {
	return "Pong!", nil
}

func (bot *Bot) Choose(_ *gateway.MessageCreateEvent, choices ...string) (string, error) {
	return choices[rand.Intn(len(choices))], nil
}
