package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transaction
// Because each *Queries only do 1 operations, we need a transaction struct, which is called composition
// Extend struct using composition instead of inheritance, so all functionality of Queries can be used here
// and can support transaction by adding more functions to Store struct
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// create a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execute a function within a dB transaction
// execTC takes a ctx and callback fn as input, start new DB tsx, create new Queries obj with that tsx,
// call the callback fn with created Queries, then commit/rollback the tsx, based on the error returned by execTC
func (store *SQLStore) execTC(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // TxOptions to set Isolation level, for beginning can use nil
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

// input param of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// result of trasnfer transction (updated result)
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// var txKey = struct{}{}

// perform money transfer from one account to another
// creates a transfer record,
// add account entries,
// update accounts' balance within single DB TX
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTC(ctx, func(q *Queries) error {
		var err error

		// txName := ctx.Value(txKey)

		// fmt.Println(txName, "Create trasnfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "Create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "Create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// update accounts' balance
		// fmt.Println(txName, "get account 1")

		// CAN REMOVE THIS AND USE addAccountBalance to save 1 query
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "update account 1 balance")
		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountID,
		// 	Balance: account1.Balance - arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// --------------------------------------------------------------------------------------------
		// REFACTOR to use AddMoney func to reduce code redundancy
		// if arg.FromAccountID < arg.ToAccountID {
		// 	result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 		ID:     arg.FromAccountID,
		// 		Amount: -arg.Amount,
		// 	})
		// 	if err != nil {
		// 		return err
		// 	}

		// 	result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 		ID:     arg.ToAccountID,
		// 		Amount: arg.Amount,
		// 	})
		// 	if err != nil {
		// 		return err
		// 	}
		// } else {
		// 	result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 		ID:     arg.ToAccountID,
		// 		Amount: arg.Amount,
		// 	})
		// 	if err != nil {
		// 		return err
		// 	}

		// 	result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 		ID:     arg.FromAccountID,
		// 		Amount: -arg.Amount,
		// 	})
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		// --------------------------------------------------------------------------------------------
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})

	return result, err
}

func AddMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return // naked return here, means default return account1, account2, err
	}

	account2, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return // without any parameters would be fine.
}
