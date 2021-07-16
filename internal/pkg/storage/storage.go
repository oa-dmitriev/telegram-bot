package storage

import (
	"database/sql"
	"os"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

func InitRedis() (*redis.Client, error) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}
	conn := redis.NewClient(opt)
	status := conn.Ping()
	if status.Err() != nil {
		return nil, status.Err()
	}
	return conn, nil
}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS vocabulary(
			user_id bigint NOT NULL,
			word text NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}
	return db, nil
}
