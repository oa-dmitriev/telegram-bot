package repository

type DBVocabulary struct {
	UserID     int64  `db:"user_id"`
	Word       string `db:"word"`
	Definition string `db:"definition"`
}
