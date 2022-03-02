package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const entryPointWaitingSeconds = 15

type Options struct {
	DatabaseURL string
}

func NewDBWrapper(ctx context.Context, opts Options) (*sql.DB, error) {
	db, err := sql.Open("postgres", opts.DatabaseURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Ready(db *sql.DB) error {
	var errPing error
	for i := 0; i < entryPointWaitingSeconds; i++ {
		log.Printf("trying to ping database, attempt %d", i)
		if errPing = db.Ping(); errPing == nil {
			log.Printf("connected to database")
			return nil
		}
		time.Sleep(time.Second)
	}
	return errPing
}
