package vocabulary

import (
	"context"
	"log"

	"github.com/oa-dmitriev/telegram-bot/internal/repository"
)

const (
	defaultCapacity = 100
)

func (r *Repo) GetList(ctx context.Context, userID, offset, limit int64) ([]*repository.DBVocabulary, error) {
	res := make([]*repository.DBVocabulary, 0, defaultCapacity)

	rows, err := r.QueryContext(ctx, querySQLSelectAllVocabulary, userID, offset, limit)
	if err != nil {
		log.Println("could not exec select query for vocabulary table, error: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		dbVocab := repository.DBVocabulary{}
		if err := rows.Scan(&dbVocab.Word, &dbVocab.Definition); err != nil {
			log.Println("could not scan the row for vocabulary table, error: ", err)
			return nil, err
		}
		res = append(res, &dbVocab)
	}
	return res, rows.Err()
}

func (r *Repo) Add(ctx context.Context, vocab *repository.DBVocabulary) error {
	args := sqlArgs(vocab)
	if _, err := r.ExecContext(ctx, querySQLInsertVocabulary, args...); err != nil {
		log.Printf("could not exec insert query for vocabulary table, error: %s", err)
		return err
	}
	log.Printf("inserted word [%s] for user [%d]", vocab.Word, vocab.UserID)
	return nil
}

func (r *Repo) DeleteWord(ctx context.Context, userID int64, word string) error {
	if _, err := r.ExecContext(ctx, querySQLDeleteWordFromVocabulary, userID, word); err != nil {
		log.Printf("could not exec delete word query for vocabulary table, error: %s", err)
		return err
	}
	log.Printf("deleted word [%s] for the user with id [%d]", word, userID)
	return nil
}

func (r *Repo) Delete(ctx context.Context, userID int64) error {
	if _, err := r.ExecContext(ctx, querySQLDeleteVocabulary, userID); err != nil {
		log.Printf("could not exec delete vocabulary query, error: %s", err)
		return err
	}
	log.Printf("deleted vocabulary for the user with id [%d]", userID)
	return nil
}
