package db

import (
	"context"
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T, from Accounts, to Accounts) Transfers {
	arg := CreateTransferParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreatTransfer(t *testing.T) {
	from := CreateRandomAccount(t)
	to := CreateRandomAccount(t)
	CreateRandomTransfer(t, from, to)
}

func TestGetTransfer(t *testing.T) {
	// Create a random account
	from := CreateRandomAccount(t)
	to := CreateRandomAccount(t)
	transfer1 := CreateRandomTransfer(t, from, to)

	// Get the transfer
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)                           // Check if the ID is the same
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)     // Check if the account ID is the same
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)         // Check if the account ID is the same
	require.Equal(t, transfer1.Amount, transfer2.Amount)                   // Check if the amount is the same
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, 0) // Check if the created time is the same
}

func TestListTransfer(t *testing.T) {
	// Create a random account
	from := CreateRandomAccount(t)
	to := CreateRandomAccount(t)

	// Create 5 random transfers
	for range 5 {
		CreateRandomTransfer(t, from, to)
	}

	// List the transfers
	arg := ListTransfersParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)      // Check if there is no error
	require.Len(t, transfers, 5) // Check if the length of the transfers is 5

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer) // Check if the transfer is not empty
	}
}
