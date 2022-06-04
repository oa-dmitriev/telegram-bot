package repository

import "time"

type DBJoke struct {
	ID         int       `db:"id"`
	Type       string    `db:"joke_type"`
	Joke       string    `db:"joke"`
	Setup      string    `db:"setup"`
	Delivery   string    `db:"delivery"`
	Category   string    `db:"category"`
	Created_at time.Time `db:"created_at"`
}
