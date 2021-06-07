package cmd

import (
	"fmt"
	"gbot/commerce"
	"gbot/commerce/banking"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/bot/extras/arguments"
	"github.com/diamondburned/arikawa/v2/gateway"
)

type BankCmd struct {
	Context *bot.Context
	bankSub *bot.Subcommand
	bank    *banking.Bank
}

func NewBankCmd(bank *banking.Bank) *BankCmd {
	if bank == nil {
		panic(fmt.Errorf("Cannot construct command. Bank is nil"))
	}
	return &BankCmd{
		bank: bank,
	}
}

func (b *BankCmd) Setup(sub *bot.Subcommand) {
	b.bankSub = sub
	sub.Command = "bank"
	sub.Description = "Interact with the bank"

	sub.ChangeCommandInfo("Tip", "", "Tip a user a certain value. eg. 'bank tip @Esbee 100'")
	sub.ChangeCommandInfo("Balance", "", "Display your current balance")
	sub.ChangeCommandInfo("List", "", "List all bank accounts")
}

func (b *BankCmd) Help(*gateway.MessageCreateEvent) (string, error) {
	return b.bankSub.Help(), nil
}

func (b *BankCmd) Tip(m *gateway.MessageCreateEvent, userString string, value int) (string, error) {
	var tippedUser arguments.UserMention
	err := tippedUser.Parse(userString)
	if err != nil {
		return "", err
	}

	err = b.bank.Withdraw(m.Author.ID, value)
	if err != nil {
		return "", err
	}

	err = b.bank.Deposit(tippedUser.ID(), value)
	if err != nil {
		return "", err
	}

	return "donezo", nil
}

func (b *BankCmd) Balance(m *gateway.MessageCreateEvent) (string, error) {
	val, err := b.bank.Balance(m.Author.ID)
	if err != nil {
		return "", err
	}

	return commerce.Currency(val).String(), nil
}

func (b *BankCmd) List(m *gateway.MessageCreateEvent) (string, error) {
	return b.bank.PrintBalances()
}
