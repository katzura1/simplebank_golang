package db

import (
	"context"
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T, account Accounts) Entries {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(), // Random amount
	}

	// Create a random entry
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	// Create a random account
	account := CreateRandomAccount(t)
	// Create a random entry
	CreateRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	// Create a random account
	account := CreateRandomAccount(t)
	// Create a random entry
	entry1 := CreateRandomEntry(t, account)

	// Get the entry
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)               // Check if the ID is the same
	require.Equal(t, entry1.AccountID, entry2.AccountID) // Check if the account ID is the same
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, 0)
}

func TestListEntry(t *testing.T) {
	// Create a random account
	account := CreateRandomAccount(t)

	// Create 5 random entries
	for range 5 {
		CreateRandomEntry(t, account)
	}

	// List the entries
	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
