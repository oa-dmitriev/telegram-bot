package user

import "github.com/oa-dmitriev/telegram-bot/internal/repository"

const (
	querySQLInsertUser = `
INSERT INTO users(id, username, first_name, last_name)
VALUES($1::BIGINT, $2::TEXT, $3::TEXT, $4::TEXT);
`
	querySQLDeleteUser = `
DELETE FROM users 
WHERE id = $1::BIGINT;	
`
)

func sqlArgs(user *repository.DBUser) []interface{} {
	return []interface{}{
		user.ID,
		user.Username,
		user.FirstName,
		user.LastName,
	}
}
