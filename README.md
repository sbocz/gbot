Discord bot for a myriad of functional and entertaining tasks powered by the [arikawa](https://github.com/diamondburned/arikawa/tree/master) project.

# Using the bot

## Configuration
Configure the environment variables manually or use the `.env` file to configure the bot.
- `BOT_TOKEN`: Unique discord bot secret token. See [discord docs](https://discord.com/developers/docs/topics/oauth2) for more info
- `BOT_PREFIX`: Prefix to trigger commands. eg. if set to `bob ` then you can use `bob ping` to ping the bot.
- `DATABASE_FILE` BoltDB database file to use. This will be created if it does not already exist.

## Testing
Test all modules using `go test ./... -v`

## Manging dependencies
`go install`
`go mod tidy`

## Running
`go run gbot`
