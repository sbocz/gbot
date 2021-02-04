package bank

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
)

type Account struct {
	UserId        discord.UserID
	Balance       int
	Age           time.Time
	LastInterest  time.Time
	InterestValue int
}

func NewAccount(id discord.UserID, balance int) *Account {
	var now = time.Now().UTC()
	return &Account{
		UserId:        discord.UserID(id),
		Balance:       balance,
		Age:           now,
		LastInterest:  now,
		InterestValue: 10,
	}
}

func (a *Account) CalculateInterest(periodMillis int) int {
	duration := time.Now().UTC().Sub(a.LastInterest).Milliseconds()
	if (int(duration) / periodMillis) < 1 {
		return 0
	}
	return (int(duration) / periodMillis) * a.InterestValue
}
