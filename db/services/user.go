package services

import (
	"context"
	"simple-bank-system/db/pkg"
)

type CreateUserParams struct {
	Username       string
	HashedPassword string
	FullName       string
	Email          string
}

func (r *DB) CreateUser(ctx context.Context, user CreateUserParams) (res *pkg.User, err error) {
	query := `INSERT INTO users(username, hashed_password, full_name, email
		) VALUES(
			$1, $2, $3, $4, $5
		) RETURNING username, hashed_password, full_name, email, password_change_at, created_at;`
	err = r.db.QueryRow(ctx, query, user.Username, user.HashedPassword, user.FullName, user.Email).Scan(&res.Username, &res.HashedPassword, &res.FullName, &res.Email, &res.PasswordChangeAt, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *DB) GetUser(ctx context.Context, username string) (*pkg.User, error) {
	row := r.db.QueryRow(ctx, "SELECT * FROM accounts WHERE id=$1;", &username)

	var user pkg.User
	if err := row.Scan(&user.Username, &user.HashedPassword, &user.FullName, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}
