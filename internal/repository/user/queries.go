package user

import "github.com/oa-dmitriev/telegram-bot/internal/repository"

const (
	querySQLInsertUser = `
INSERT INTO users(id, username, first_name, last_name, chat_id)
VALUES($1::BIGINT, $2::TEXT, $3::TEXT, $4::TEXT, $5::BIGINT);
`
	querySQLDeleteUser = `
DELETE FROM users 
WHERE id = $1::BIGINT;	
`
	querySQLGetUser = `
SELECT id, username, first_name, last_name, chat_id 
FROM users 
WHERE id = $1::BIGINT;
`
)

func sqlArgs(user *repository.DBUser) []any {
	return []any{
		user.ID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.ChatID,
	}
}
