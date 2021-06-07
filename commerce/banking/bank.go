package banking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gbot/commerce"
	"gbot/database"
	"gbot/util"
	"text/tabwriter"
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
// Returns an error if one occurred during the operation.
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
// Returns an error if one occurred during the operation.
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

// Pays interest for all accounts
func (b *Bank) PayInterest() error {
	userIds, err := b.fetchUserIds()
	if err != nil {
		return err
	}
	for _, userId := range userIds {
		err = b.payAccountInterest(userId)
		if err != nil {
			return err
		}
	}
	return nil
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

// Prints a table of the balances for all accounts.
// Returns a string for the textual table and an error if one occurred while retrieving any account.
func (b *Bank) PrintBalances() (string, error) {
	balances, err := b.fetchBalances()
	if err != nil {
		return "", fmt.Errorf("Error retrieving bank balances: %v", err)
	}
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "User\tBalance\n")

	for userId, balance := range balances {
		fmt.Fprintf(w, "%s\t%v\n",
			userId.Mention(),
			commerce.Currency(balance))
	}
	w.Flush()
	return buffer.String(), nil
}

func (b *Bank) fetchUserIds() ([]discord.UserID, error) {
	keys, err := b.accountBucket.Keys()
	if err != nil {
		return nil, err
	}
	userIds, err := util.UserIdsFromBytes(keys)
	if err != nil {
		return nil, err
	}
	return userIds, nil
}

func (b *Bank) payAccountInterest(userId discord.UserID) error {
	account, err := b.FetchAccount(userId)
	if err != nil {
		return fmt.Errorf("Could not pay interest: %v", err)
	}

	interest := account.CalculateInterest(INTEREST_PERIOD)
	if interest < 1 {
		return nil
	}
	account.Balance += interest
	account.LastInterest = time.Now().UTC()

	return b.SaveAccount(account)
}

func (b *Bank) fetchBalances() (map[discord.UserID]int, error) {
	userIds, err := b.fetchUserIds()
	if err != nil {
		return nil, err
	}
	balances := make(map[discord.UserID]int)
	for _, userId := range userIds {
		balance, err := b.Balance(userId)
		if err != nil {
			return nil, err
		}
		balances[userId] = balance
	}
	return balances, nil
}
