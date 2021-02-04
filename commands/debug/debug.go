package debug

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/bot/extras/middlewares"
	"github.com/diamondburned/arikawa/v2/gateway"
)

// Debug Flag for administrators only.
type Debug struct {
	Context  *bot.Context
	debugSub *bot.Subcommand
}

// Setup demonstrates the CanSetup interface. This function will never be parsed
// as a callback of any event.
func (d *Debug) Setup(sub *bot.Subcommand) {
	d.debugSub = sub
	// Set a custom command (e.g. "!go ..."):
	sub.Command = "go"
	// Set a custom description:
	sub.Description = "Print Go debugging variables"

	// Manually set the usage for each function.
	sub.ChangeCommandInfo("GOOS", "GOOS", "Prints the current operating system")
	sub.ChangeCommandInfo("GC", "GC", "Triggers the garbage collector")
	sub.ChangeCommandInfo("Goroutines", "", "Prints the current number of Goroutines")

	sub.Hide("Die")
	sub.AddMiddleware("Die", middlewares.AdminOnly(d.Context))
}

func (d *Debug) Help(*gateway.MessageCreateEvent) (string, error) {
	return d.debugSub.Help(), nil
}

// ~go goroutines
func (d *Debug) Goroutines(*gateway.MessageCreateEvent) (string, error) {
	return fmt.Sprintf(
		"goroutines: %d",
		runtime.NumGoroutine(),
	), nil
}

// ~go GOOS
func (d *Debug) GOOS(*gateway.MessageCreateEvent) (string, error) {
	return strings.Title(runtime.GOOS), nil
}

// ~go GC
func (d *Debug) GC(*gateway.MessageCreateEvent) (string, error) {
	runtime.GC()
	return "Done.", nil
}

// ~go die
// This command will be hidden from ~help by default.
func (d *Debug) Die(m *gateway.MessageCreateEvent) error {
	log.Fatalln("User", m.Author.Username, "killed the bot x_x")
	return nil
}
