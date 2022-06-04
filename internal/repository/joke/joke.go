package joke

import (
	"context"
	"log"

	"github.com/oa-dmitriev/telegram-bot/internal/repository"
)

const defaultCapacity = 5000

func (r *Repo) Add(ctx context.Context, joke *repository.DBJoke) error {
	args := sqlArgs(joke)
	if _, err := r.ExecContext(ctx, querySQLInsertJoke, args...); err != nil {
		log.Printf("could not exec querySQLInsertJoke, error: %s\n", err)
		return err
	}
	return nil
}

func (r *Repo) GetList(ctx context.Context, limit, offset int64) ([]*repository.DBJoke, error) {
	res := make([]*repository.DBJoke, 0, defaultCapacity)

	rows, err := r.QueryContext(ctx, querySQLSelectJokes, limit, offset)
	if err != nil {
		log.Println("could not exec querySQLSelectJokes, error: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		dbJoke := repository.DBJoke{}

		if err := rows.Scan(&dbJoke.ID, &dbJoke.Type, &dbJoke.Setup,
			&dbJoke.Delivery, &dbJoke.Category, &dbJoke.Created_at,
		); err != nil {
			log.Println("could not scan the row in selectJokesQuery, error: ", err)
			return nil, err
		}
		res = append(res, &dbJoke)
	}
	return res, rows.Err()
}

func sqlArgs(joke *repository.DBJoke) []any {
	return []any{
		joke.ID,
		joke.Type,
		joke.Joke,
		joke.Setup,
		joke.Delivery,
		joke.Category,
	}
}
