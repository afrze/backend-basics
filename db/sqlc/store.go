package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provide all function to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameter of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// var txKey = struct{}{}
// Custom type definition
//type contextKey string

// Define a constant for the key:
//const txKey = contextKey("transaction")

// TransferTx performs a money transfer from one account to other
// It creates a transfer record, add account entries, and update accounts' balance within a single database transactions
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		//txName := ctx.Value(txKey)

		//fmt.Println(txName, "Create Transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		//fmt.Println(txName, "Create Entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, "Create Entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, "Get account 1")
		// get account -> update its balance
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		//fmt.Println(txName, "update account 1 balance")
		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountID,
		// 	Balance: account1.Balance - arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		/*
			if arg.FromAccountID < arg.ToAccountID {

				result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
					ID:     arg.FromAccountID,
					Amount: -arg.Amount,
				})
				if err != nil {
					return err
				}

				//fmt.Println(txName, "get account 2")
				// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
				// if err != nil {
				// 	return err
				// }

				//fmt.Println(txName, "Update account 2 balance")
				// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				// 	ID:      arg.ToAccountID,
				// 	Balance: account2.Balance + arg.Amount,
				// })
				result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
					ID:     arg.ToAccountID,
					Amount: arg.Amount,
				})
				if err != nil {
					return err
				}

			} else {
				result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
					ID:     arg.ToAccountID,
					Amount: arg.Amount,
				})
				if err != nil {
					return err
				}

				result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
					ID:     arg.FromAccountID,
					Amount: -arg.Amount,
				})
				if err != nil {
					return err
				}
			}
		*/

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, _ = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)

		} else {
			result.ToAccount, result.FromAccount, _ = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}
