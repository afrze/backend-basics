package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/afrze/backend-basics/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	accountOne := createRandomAccount(t)
	accountTwo, err := testQueries.GetAccount(context.Background(), accountOne.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountTwo)

	require.Equal(t, accountOne.ID, accountTwo.ID)
	require.Equal(t, accountOne.Owner, accountTwo.Owner)
	require.Equal(t, accountOne.Balance, accountTwo.Balance)
	require.Equal(t, accountOne.Currency, accountTwo.Currency)
	require.WithinDuration(t, accountOne.CreatedAt, accountTwo.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	accountOne := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      accountOne.ID,
		Balance: utils.RandomMoney(),
	}

	accountTwo, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountTwo)

	require.Equal(t, accountOne.ID, accountTwo.ID)
	require.Equal(t, accountOne.Owner, accountTwo.Owner)
	require.Equal(t, arg.Balance, accountTwo.Balance)
	require.Equal(t, accountOne.Currency, accountTwo.Currency)
	require.WithinDuration(t, accountOne.CreatedAt, accountTwo.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	accountOne := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), accountOne.ID)
	require.NoError(t, err)

	accountTwo, err := testQueries.GetAccount(context.Background(), accountOne.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountTwo)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
