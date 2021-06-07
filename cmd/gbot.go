package cmd

import (
	"fmt"
	"math/rand"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

type GbotCmd struct {
	Ctx *bot.Context
}

func NewGbotCmd() *GbotCmd {
	return &GbotCmd{}
}

func (g *GbotCmd) Setup(sub *bot.Subcommand) {
	sub.ChangeCommandInfo("Add", "add", "adds some numbers")
	sub.ChangeCommandInfo("Ping", "ping", "check if bot is alive")
	sub.ChangeCommandInfo("Inspire", "inspire", "display an inspirational message")
	sub.ChangeCommandInfo("Choose", "choose", "choose from a list of options")
}

// Help prints the default help message.
func (g *GbotCmd) Help(*gateway.MessageCreateEvent) (string, error) {
	return g.Ctx.Help(), nil
}

func (g *GbotCmd) Add(_ *gateway.MessageCreateEvent, a, b int) (string, error) {
	return fmt.Sprintf("%d + %d = %d", a, b, a+b), nil
}

func (g *GbotCmd) Ping(*gateway.MessageCreateEvent) (string, error) {
	return "Pong!", nil
}

func (g *GbotCmd) Choose(_ *gateway.MessageCreateEvent, choices ...string) (string, error) {
	return choices[rand.Intn(len(choices))], nil
}
