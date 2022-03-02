package repository

import "context"

type UserRepo interface {
	Add(ctx context.Context, user *DBUser) error
}

type VocabularyRepo interface {
	GetList(ctx context.Context, userID, limit, offset int64) ([]*DBVocabulary, error)
	Add(ctx context.Context, vocab *DBVocabulary) error
	DeleteWord(ctx context.Context, userID int64, word string) error
	Delete(ctx context.Context, userID int64) error
}