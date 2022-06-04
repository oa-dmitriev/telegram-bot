package user

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/oa-dmitriev/telegram-bot/internal/repository"
)

var ErrUserNotFound error = errors.New("user not found")

func (r *Repo) Add(ctx context.Context, user *repository.DBUser) error {
	args := sqlArgs(user)
	if _, err := r.ExecContext(ctx, querySQLInsertUser, args...); err != nil {
		log.Printf("could not exec query for user table, error: %s", err)
		return err
	}
	return nil
}

func (r *Repo) GetUser(ctx context.Context, userID int64) (*repository.DBUser, error) {
	rows := r.QueryRowContext(ctx, querySQLGetUser, userID)
	if errors.Is(rows.Err(), sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	user := repository.DBUser{}
	if err := rows.Scan(&user); err != nil {
		log.Println("could not exec querySQLGetUser query, error: ", err)
		return nil, err
	}
	return &user, nil
}
