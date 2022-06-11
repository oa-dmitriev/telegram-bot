package user

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/oa-dmitriev/telegram-bot/internal/repository"
)

var ErrUserNotFound error = errors.New("user not found")

const defaultCapacity = 100

func (r *Repo) Add(ctx context.Context, user *repository.DBUser) error {
	args := sqlArgs(user)
	if _, err := r.ExecContext(ctx, querySQLInsertUser, args...); err != nil {
		log.Printf("could not exec query for user table, error: %s", err)
		return err
	}
	return nil
}

func (r *Repo) GetUser(ctx context.Context, userID int64) (*repository.DBUser, error) {
	row := r.QueryRowContext(ctx, querySQLGetUser, userID)
	user := repository.DBUser{}
	if err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.ChatID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		log.Println("could not exec querySQLGetUser query, error: ", err)
		return nil, err
	}
	return &user, nil
}

func (r *Repo) GetAllUsers(ctx context.Context) ([]*repository.DBUser, error) {
	rows, err := r.QueryContext(ctx, querySQLGetAllUsers)
	if err != nil {
		log.Println("could not exec querySQLGetAllUsers: ", err)
		return nil, err
	}
	defer rows.Close()

	res := make([]*repository.DBUser, 0, defaultCapacity)
	for rows.Next() {
		user := repository.DBUser{}
		if err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.ChatID); err != nil {
			log.Println("could not scan the row for vocabulary table, error: ", err)
			return nil, err
		}
		res = append(res, &user)
	}
	return res, nil
}

func sqlArgs(user *repository.DBUser) []any {
	return []any{
		user.ID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.ChatID,
	}
}
