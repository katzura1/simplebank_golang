package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "tom",
		Balance:  100,
		Currency: "USD",
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)      // Check if there is no error
	require.NotEmpty(t, account) // Check if the account is not empty

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	arg := UpdateAccountParams{
		ID:       1,
		Owner:    "budi",
		Balance:  200,
		Currency: "USD",
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.ID, account.ID)
	require.Equal(t, arg.Balance, account.Balance)
}

func TestDeleteAccount(t *testing.T) {
	err := testQueries.DeleteAccount(context.Background(), 3)
	require.NoError(t, err)
}
