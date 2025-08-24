package sqlite

import (
	"database/sql"
	"fmt"
)

func ConnectToSqlite() (*sql.DB, error) {
	dbPath := "/home/ad/.local/share/wakapi/wakapi_db.db"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close() // close if ping fails
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Successful connection, do not close here
	return db, nil
}
