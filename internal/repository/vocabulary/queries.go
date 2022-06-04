package vocabulary

const (
	querySQLInsertVocabulary = `
INSERT INTO vocabulary(user_id, word, definition)
VALUES($1::BIGINT, $2::TEXT, $3::TEXT);
`
	querySQLDeleteWordFromVocabulary = `
DELETE FROM vocabulary
WHERE id = $1::BIGINT AND word = $2::TEXT;	
`

	querySQLDeleteVocabulary = `
DELETE FROM vocabulary 
WHERE id = $1::BIGINT;
`

	querySQLSelectAllVocabulary = `
SELECT user_id, word, definition
FROM vocabulary
WHERE user_id = $1::BIGINT 
ORDER BY word
LIMIT $2::INT
OFFSET $3::INT;	
`
)
