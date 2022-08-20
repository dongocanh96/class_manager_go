package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

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

type UpdateUserTxParams struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	Fullname          string    `json:"fullname"`
	Email             string    `json:"email"`
	PhoneNumber       string    `json:"phone_number"`
}

type UpdateUserTxResult struct {
	User User `json:"user"`
}

func (store *Store) UpdateUserTx(ctx context.Context, arg UpdateUserTxParams) (UpdateUserTxResult, error) {
	var result UpdateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			ID:                arg.ID,
			Username:          arg.Username,
			HashedPassword:    arg.HashedPassword,
			PasswordChangedAt: arg.PasswordChangedAt,
			Fullname:          arg.Fullname,
			Email:             arg.Email,
			PhoneNumber:       arg.PhoneNumber,
		})

		if err != nil {
			return err
		}

		return nil

	})

	return result, err
}
