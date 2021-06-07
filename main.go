package main

import (
	"gbot/cmd"
	"gbot/commerce/banking"
	"gbot/database"
	"gbot/task"
	"log"
	"os"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/utils/handler"
	"github.com/joho/godotenv"
)

// To run, do `BOT_TOKEN="TOKEN HERE" go run .`

func main() {
	godotenv.Load()

	dbFile := os.Getenv("DATABASE_FILE")
	db := database.NewDb(dbFile)
	defer db.Shutdown()

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}
	prefix := os.Getenv("BOT_PREFIX")
	if token == "" {
		prefix = "!"
		log.Println("No $BOT_PREFIX given. Defaulting to '!'")
	}
	denyListString := os.Getenv("YELL_DENYLIST")
	denyListString = strings.ToLower(denyListString)
	denyList := strings.Split(denyListString, ",")

	cmd.Initialize(db, denyList)
	bank, err := banking.NewBank(db)
	if err != nil {
		log.Fatalln(err)
	}

	// Start background tasks
	go task.PayInterest(24*time.Hour, bank)

	wait, err := bot.Start(token, gbot, func(ctx *bot.Context) error {
		ctx.HasPrefix = bot.NewPrefix(prefix)
		ctx.EditableCommands = true
		ctx.State.PreHandler = handler.New()
		ctx.State.PreHandler.Synchronous = true
		ctx.State.PreHandler.AddHandler(cmd.YellHandler(ctx))

		ctx.MustRegisterSubcommand(&cmd.Debug{})
		ctx.MustRegisterSubcommand(cmd.NewBankCmd(bank))

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
