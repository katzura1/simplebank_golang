package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Accounts {
	user := CreateRandomUser(t) // Create a random user

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)      // Check if there is no error
	require.NotEmpty(t, account) // Check if the account is not empty

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	// Create a new account
	account1 := CreateRandomAccount(t)
	// Delete the account
	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	// Check if the account is deleted
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)      // Check if there is an error
	require.Empty(t, account2) // Check if the account is empty
}

func TestListAccount(t *testing.T) {
	// Create 5 random accounts
	for range 5 {
		CreateRandomAccount(t)
	}

	// List the accounts
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	// Check if the accounts are not empty
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)     // Check if there is no error
	require.Len(t, accounts, 5) // Check if the length of the accounts is 5

	for _, account := range accounts {
		require.NotEmpty(t, account) // Check if the account is not empty
	}
}
