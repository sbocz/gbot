package main

import (
	"gbot/commands"
	"gbot/commands/debug"
	"gbot/commands/misc"
	"gbot/database"
	"log"
	"os"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/utils/handler"
	"github.com/joho/godotenv"
)

// To run, do `BOT_TOKEN="TOKEN HERE" go run .`

func main() {
	godotenv.Load()

	var dbFile = os.Getenv("DATABASE_FILE")
	var db = database.NewDb(dbFile)
	defer db.Shutdown()

	var token = os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}
	var prefix = os.Getenv("BOT_PREFIX")
	if token == "" {
		prefix = "!"
		log.Println("No $BOT_PREFIX given. Defaulting to '!'")
	}

	misc.Initialize(db)

	commands := &commands.Bot{}

	wait, err := bot.Start(token, commands, func(ctx *bot.Context) error {
		ctx.HasPrefix = bot.NewPrefix(prefix)
		ctx.EditableCommands = true
		ctx.State.PreHandler = handler.New()
		ctx.State.PreHandler.Synchronous = true
		ctx.State.PreHandler.AddHandler(misc.YellHandler(ctx))

		// Subcommand demo, but this can be in another package.
		ctx.MustRegisterSubcommand(&debug.Debug{})

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Bot started")

	// As of this commit, wait() will block until SIGINT or fatal. The past
	// versions close on call, but this one will block.
	// If for some reason you want the Cancel() function, manually make a new
	// context.
	if err := wait(); err != nil {
		log.Fatalln("Gateway fatal error:", err)
	}
}
