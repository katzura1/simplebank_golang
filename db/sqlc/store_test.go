package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	// Create two random accounts
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println("Before:", account1.Balance, account2.Balance)

	// run and concurrently transfer money between the two accounts
	n := 2
	amount := int64(10)

	errors := make(chan error)             // channel to collect errors
	results := make(chan TransferTxResult) // channel to collect results

	for k := range n {
		txName := fmt.Sprintf("Tx %d", k)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errors <- err
			results <- result
		}()
	}

	// check the results
	existed := make(map[int]bool)
	for range n {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testQueries.GetTransfer(context.Background(), transfer.ID) // check if the transfer is stored in the database
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testQueries.GetEntry(context.Background(), fromEntry.ID) // check if the entry is stored in the database
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testQueries.GetEntry(context.Background(), toEntry.ID) // check if the entry is stored in the database
		require.NoError(t, err)

		// Check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// Check account balances
		fmt.Println("Tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance // e.g. before: 100, after: 90, diff: 10 (100 - 90 = 10)
		diff2 := toAccount.Balance - account2.Balance   // e.g. before: 100, after: 110, diff: 10 (110 - 100 = 10)
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)         // check if the diff is positive
		require.True(t, diff1%amount == 0) // check if the diff is a multiple of the amount, 1 * amount, 2 * amount, 3 * amount, etc.

		//example
		// 1. account1: 863, account2: 140
		// simulate loopings
		// 1. diff1 = 863 - 853 = 10, diff2 = 150 - 140 = 10
		// 2. diff1 = 863 - 843 = 20, diff2 = 160 - 140 = 20
		// 3. diff1 = 863 - 833 = 30, diff2 = 170 - 140 = 30
		// 4. diff1 = 863 - 823 = 40, diff2 = 180 - 140 = 40
		// 5. diff1 = 863 - 813 = 50, diff2 = 180 - 140 = 50

		k := int(diff1 / amount)
		// simulate loopings
		// 1. k = 10 / 10 = 1
		// 2. k = 20 / 10 = 2
		// 3. k = 30 / 10 = 3
		// 4. k = 40 / 10 = 4
		// 5. k = 50 / 10 = 5

		require.True(t, k >= 1 && k <= n)  // check if k is between 1 and n
		require.NotContains(t, existed, k) // check if k is not in existed
		existed[k] = true
		// fmt.Println("K:", k, "Existed:", existed, "Diff1:", diff1, "Diff2:", diff2)
	}

	// check the final account balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	fmt.Println("After:", updatedAccount1.Balance, updatedAccount2.Balance)
}
