package db

import (
	"context"
	"database/sql"
	"fmt"
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

type UpdateUserInfoTxParams struct {
	ID          int64          `json:"id"`
	Username    sql.NullString `json:"username"`
	Fullname    sql.NullString `json:"fullname"`
	Email       sql.NullString `json:"email"`
	PhoneNumber sql.NullString `json:"phone_number"`
}

type UpdateUserInfoTxResult struct {
	User User `json:"user"`
}

func (store *Store) UpdateUserInfoTx(ctx context.Context, arg UpdateUserInfoTxParams) (UpdateUserInfoTxResult, error) {
	var result UpdateUserInfoTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.GetUserForUpdate(ctx, arg.ID)

		result.User, err = q.UpdateUserInfo(ctx, UpdateUserInfoParams{
			ID:          arg.ID,
			Username:    arg.Username,
			Fullname:    arg.Fullname,
			Email:       arg.Email,
			PhoneNumber: arg.PhoneNumber,
		})

		if err != nil {
			return err
		}

		return nil

	})

	return result, err
}
