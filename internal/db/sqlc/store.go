package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := s.Queries.WithTx(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	AmountCents   int64 `json:"amount_cents"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// Create the transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		// Create entry for the sender (negative amount)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:   pgtype.Int8{Int64: arg.FromAccountID, Valid: true},
			AmountCents: -arg.AmountCents,
		})
		if err != nil {
			return err
		}

		// Create entry for the recipient (positive amount)
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:   pgtype.Int8{Int64: arg.ToAccountID, Valid: true},
			AmountCents: arg.AmountCents,
		})
		if err != nil {
			return err
		}

		// Update account balances
		if arg.FromAccountID < arg.ToAccountID {
			if result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.AmountCents, arg.ToAccountID, arg.AmountCents); err != nil {
				return err
			}
		} else {
			if result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.AmountCents, arg.FromAccountID, -arg.AmountCents); err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amountCents1 int64,
	accountID2 int64,
	amountCents2 int64,
) (account1 Account, account2 Account, err error) {
	if account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amountCents1,
	}); err != nil {
		return
	}

	if account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amountCents2,
	}); err != nil {
		return
	}

	return
}
