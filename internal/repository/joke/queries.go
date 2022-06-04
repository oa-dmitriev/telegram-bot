package joke

const (
	querySQLInsertJoke = `
INSERT INTO jokes(id, joke_type, joke, setup, delivery, category)
VALUES($1::BIGINT, $2::TEXT, $3::TEXT, $4::TEXT, $5::TEXT, $6::TEXT);
`

	querySQLSelectJokes = `
SELECT id, joke_type, setup, delivery, category, created_at
FROM jokes
ORDER BY id
LIMIT $2::INT
OFFSET $3::INT;
`
)
