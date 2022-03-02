package user

import (
	"context"
	"log"

	"github.com/oa-dmitriev/telegram-bot/internal/repository"
)

func (r *Repo) Add(ctx context.Context, user *repository.DBUser) error {
	args := sqlArgs(user)
	if _, err := r.ExecContext(ctx, querySQLInsertUser, args...); err != nil {
		log.Printf("could not exec query for user table, error: %s", err)
		return err
	}

	log.Printf("inserted user with [%d] id", user.ID)
	return nil
}
