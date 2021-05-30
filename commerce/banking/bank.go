package banking

import (
	"encoding/json"
	"fmt"
	"gbot/database"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
)

const (
	INTEREST_PERIOD = time.Hour * 4
)

type Bank struct {
	accountBucket *database.Bucket
}

// Initializes a new Banck object using the database provided for state management
func NewBank(db *database.DB) (*Bank, error) {
	if db == nil {
		return nil, fmt.Errorf("'Bank' creation failed as DB was nil.")
	}

	bucket, err := database.NewBucket(db, database.BankAccounts)
	if err != nil {
		return nil, fmt.Errorf("could not initialize 'Bank': %v", err)
	}

	return &Bank{
		accountBucket: bucket,
	}, nil
}

// Retrieves the balance for the user's account.
// Returns the balance and an error if one occurred while retrieving the account.
func (b *Bank) Balance(userId discord.UserID) (int, error) {
	account, err := b.FetchAccount(userId)
	if err != nil {
		return 0, fmt.Errorf("Cannot show balance: %s", err)
	}
	return account.Balance, nil
}

// Deposits the provided value into the account for the user provided.
// Returns the value after the deposit and an error if one occurred during the operation.
func (b *Bank) Deposit(userId discord.UserID, value int) error {
	if value < 1 {
		return fmt.Errorf("Cannot deposit a value less than 1")
	}

	account, err := b.FetchAccount(userId)
	if err != nil {
		return fmt.Errorf("Cannot fetch account for deposit: %s", err)
	}

	account.Balance += value
	return b.SaveAccount(account)
}

// Withdraws the provided value from the user's account. Withdrawals are
// only allowed for positive values and cannot leave the account with a negative balance.
// Returns the value after the withdrawal and an error if one occurred during the operation.
func (b *Bank) Withdraw(userId discord.UserID, value int) error {
	if value < 1 {
		return fmt.Errorf("Cannot withdraw a value less than 1")
	}
	account, err := b.FetchAccount(userId)
	if err != nil {
		return fmt.Errorf("Cannot withdraw: %s", err)
	}
	if account.Balance-value < 0 {
		return fmt.Errorf("Cannot withraw to below 0. User's balance is %v", account.Balance)
	}
	account.Balance -= value

	return b.SaveAccount(account)
}

// Pays interest for the specified account
func (b *Bank) PayInterest(userId discord.UserID) error {
	account, err := b.FetchAccount(userId)
	if err != nil {
		return fmt.Errorf("could not pay interest: %v", err)
	}

	interest := account.CalculateInterest(INTEREST_PERIOD)
	account.Balance += interest
	account.LastInterest = time.Now().UTC()

	return b.SaveAccount(account)
}

// Fetch account from database. Create account if one does not exist.
func (b *Bank) FetchAccount(userId discord.UserID) (*Account, error) {
	rawValue, err := b.accountBucket.Get(fmt.Sprint(userId))
	if err != nil {
		return nil, err
	}

	// No result in database
	if rawValue == nil {
		account := NewAccount(userId, 0, time.Now().UTC())
		err := b.SaveAccount(account)
		if err != nil {
			return nil, fmt.Errorf("User specified does not have a bank account and account creation failed.")
		}
		return account, nil
	}

	var account *Account
	err = json.Unmarshal(rawValue, &account)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshall response %s: %s", rawValue, err)
	}
	return account, nil
}

func (b *Bank) SaveAccount(account *Account) error {
	return b.accountBucket.Put(fmt.Sprint(account.UserId), account)
}
