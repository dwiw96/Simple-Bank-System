package services

import (
	"context"
	"errors"
	"simple-bank-system/db/pkg"
	"simple-bank-system/util"

	"github.com/jackc/pgconn"
)

type CreateUserParams struct {
	Username       string
	HashedPassword string
	FullName       string
	Email          string
}

func (r *DB) CreateUser(ctx context.Context, user CreateUserParams) (*pkg.User, error) {
	var res pkg.User
	hashedPass, err := util.HashingPassword(user.HashedPassword)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO users(username, hashed_password, full_name, email
		) VALUES(
			$1, $2, $3, $4
		) RETURNING username, full_name, email, password_change_at, created_at;`
	err = r.db.QueryRow(ctx, query, user.Username, hashedPass, user.FullName, user.Email).Scan(&res.Username, &res.FullName, &res.Email, &res.PasswordChangeAt, &res.CreatedAt)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			// 23505 (unique_violation) ->  Duplicate account with the same currency
			if pgxError.Code == "23505" {
				return nil, util.ErrUser
			}
		}
		return nil, err
	}

	return &res, nil
}

func (r *DB) GetUser(ctx context.Context, username string) (*pkg.User, error) {
	row := r.db.QueryRow(ctx, "SELECT * FROM accounts WHERE id=$1;", &username)

	var user pkg.User
	if err := row.Scan(&user.Username, &user.HashedPassword, &user.FullName, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}
