package joke

import "database/sql"

type Repo struct {
	*sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		DB: db,
	}
}
