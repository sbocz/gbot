package banking

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
)

const (
	INTEREST_VALUE = 10
)

type Account struct {
	UserId        discord.UserID
	Balance       int
	CreationDate  time.Time
	LastInterest  time.Time
	InterestValue int
}

func NewAccount(id discord.UserID, balance int, creationDate time.Time) *Account {
	var now = time.Now().UTC()
	if creationDate.After(now) {
		creationDate = now
	}
	return &Account{
		UserId:        discord.UserID(id),
		Balance:       balance,
		CreationDate:  creationDate,
		LastInterest:  creationDate,
		InterestValue: INTEREST_VALUE,
	}
}

func (a *Account) CalculateInterest(interestPeriod time.Duration) int {
	duration := time.Now().UTC().Sub(a.LastInterest)
	if (duration / interestPeriod) < 1 {
		return 0
	}
	return int(duration/interestPeriod) * a.InterestValue
}
