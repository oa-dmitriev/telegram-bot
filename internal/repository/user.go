package repository

type DBUser struct {
	ID        int64  `db:"id"`
	Username  string `db:"username"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	ChatID    int64  `db:"chat_id"`
}
