package bank

import (
	"encoding/json"
	"fmt"
	"gbot/database"

	"github.com/diamondburned/arikawa/v2/discord"
)

type Bank struct {
	accountBucket *database.Bucket
}

func NewBank(db *database.DB) *Bank {
	return &Bank{
		accountBucket: database.NewBucket(db, database.BankAccounts),
	}
}

func (b *Bank) Deposit(userId discord.UserID, value int) error {
	if value < 1 {
		return fmt.Errorf("Cannot deposit a value less than 1")
	}
	rawValue, err := b.accountBucket.Get(fmt.Sprint(userId))
	if err != nil {
		return err
	}
	var account *Account
	err = json.Unmarshal(rawValue, &account)
	if err != nil {
		return fmt.Errorf("Could not unmarshall response %s: %s", rawValue, err)
	}
	if account == nil {
		return fmt.Errorf("User specified does not have a bank account.")
	}
	account.Balance += value
	err = b.accountBucket.Put(fmt.Sprint(userId), account)
	return err
}

func (b *Bank) Withdraw(userId discord.UserID, value int) error {
	if value < 1 {
		return fmt.Errorf("Cannot deposit a value less than 1")
	}
	rawValue, err := b.accountBucket.Get(fmt.Sprint(userId))
	if err != nil {
		return err
	}
	var account *Account
	err = json.Unmarshal(rawValue, &account)
	if err != nil {
		return fmt.Errorf("Could not unmarshall response %s: %s", rawValue, err)
	}
	if account == nil {
		return fmt.Errorf("User specified does not have a bank account.")
	}
	if account.Balance-value < 0 {
		return fmt.Errorf("Cannot withraw to below 0.")
	}
	account.Balance += value
	err = b.accountBucket.Put(fmt.Sprint(userId), account)
	return err
}
