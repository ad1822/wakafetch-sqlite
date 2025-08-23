package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ad1822/wakafetch-sqlite/types"
)

func FetchDataForHeatMap(from, to time.Time) error {
	dbPath := "/home/ad/.local/share/wakapi/wakapi_db.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	fmt.Println("Connected to SQLite database successfully!")

	query := `
		SELECT strftime('%Y-%m-%d', time) AS day, COUNT(*)
		FROM heartbeats
		WHERE user_id = ?
		  AND time >= ? 
		  AND time < ? 
		GROUP BY day
		ORDER BY day ASC;
	`

	rows, err := db.Query(query, "ad1822", from.Format(time.RFC3339), to.Format(time.RFC3339))
	if err != nil {
		return err
	}
	defer rows.Close()

	var sum int64
	sum = 0

	var activities []types.DailyActivity
	for rows.Next() {
		var day string
		var count int64
		if err := rows.Scan(&day, &count); err != nil {
			return err
		}
		sum = count + sum
		parsedDate, err := time.Parse("2006-01-02", day)
		if err != nil {
			return fmt.Errorf("invalid date format in DB: %w", err)
		}
		activities = append(activities, types.DailyActivity{
			Count: count,
			Date:  parsedDate,
		})
		fmt.Println(day, count)
	}
	fmt.Println(sum)
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
