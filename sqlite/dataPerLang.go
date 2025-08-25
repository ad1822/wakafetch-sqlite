package sqlite

import (
	"database/sql"
)

func GetLangData(db *sql.DB) (*sql.Rows, error) {
	query := `
	SELECT language, COUNT(language)
	FROM durations  
	WHERE user_id=?  AND language != ''
	GROUP BY language 
	ORDER BY count(language) desc
	LIMIT 7;
	`
	return db.Query(query, "ad1822")

	// return rows, nil
}
