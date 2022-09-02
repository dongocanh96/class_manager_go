package api

import (
	"database/sql"
	"testing"

	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/dongocanh96/class_manager_go/util"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       sql.NullString{String: util.RandomString(6), Valid: true},
		HashedPassword: hashedPassword,
		Fullname:       sql.NullString{String: util.RandomString(6), Valid: true},
		Email:          sql.NullString{String: util.RandomEmail(), Valid: true},
		PhoneNumber:    sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
		IsTeacher:      util.RandomBoolean(),
	}

	return
}
