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

func main() {
	godotenv.Load()

	// Initialize all users of db
	dbFile := os.Getenv("DATABASE_FILE")
	db := database.NewDb(dbFile)
	defer db.Shutdown()

	bank, err := banking.NewBank(db)
	if err != nil {
		log.Fatalln(err)
	}

	// Start background tasks
	go task.PayInterest(24*time.Hour, bank)

	// Initialize and start the bot
	cmd.Initialize(db, fetchDenyList())
	gbot := cmd.NewGbotCmd()
	wait, err := bot.Start(fetchToken(), gbot, func(ctx *bot.Context) error {
		ctx.HasPrefix = bot.NewPrefix(fetchPrefix())
		ctx.EditableCommands = true
		ctx.State.PreHandler = handler.New()
		ctx.State.PreHandler.Synchronous = true
		ctx.State.PreHandler.AddHandler(cmd.YellHandler(ctx))

		ctx.MustRegisterSubcommand(cmd.NewDebugCmd())
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

func fetchDenyList() []string {
	denyListString := os.Getenv("YELL_DENYLIST")
	return strings.Split(strings.ToLower(denyListString), ",")
}

func fetchToken() string {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}
	return token
}

func fetchPrefix() string {
	prefix := os.Getenv("BOT_PREFIX")
	if prefix == "" {
		prefix = "!"
		log.Println("No $BOT_PREFIX given. Defaulting to '!'")
	}
	return prefix
}
