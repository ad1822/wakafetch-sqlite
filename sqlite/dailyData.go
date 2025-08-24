package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
)

const userID = "ad1822"

// ConnectToSqlite opens SQLite connection

// formatHrsMin converts seconds to "Xh Ym"
func formatHrsMin(seconds int64) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// printStat prints aligned label-value with colors
func printStat(label, value string, labelColor, valueColor *color.Color) {
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("%s%s %s\n", cyan(string(label[0])), labelColor.Sprint(label[1:]), valueColor.Sprint(value))
}

// FetchTotalTime fetches total coding time in seconds
func FetchTotalTime(db *sql.DB, period string) (int64, error) {
	query := `
        SELECT IFNULL(SUM(duration)/1000000000,0)
        FROM durations
        WHERE user_id = ?
    `
	if period == "today" {
		query += " AND DATE(time) = DATE('now', 'localtime')"
	}

	var total sql.NullInt64
	err := db.QueryRow(query, userID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total.Int64, nil
}

// FetchDailyAvg calculates daily average coding time
func FetchDailyAvg(db *sql.DB) (int64, error) {
	var total, days sql.NullInt64
	db.QueryRow(`
        SELECT IFNULL(SUM(duration)/1000000000,0) FROM durations WHERE user_id = ?
    `, userID).Scan(&total)
	db.QueryRow(`
        SELECT COUNT(DISTINCT DATE(time,'localtime')) FROM durations WHERE user_id = ?
    `, userID).Scan(&days)

	if days.Int64 == 0 {
		return 0, nil
	}
	return total.Int64 / days.Int64, nil
}

// FetchTopEditor fetches top editor
func FetchTopEditor(db *sql.DB, period string) (string, error) {
	query := `
        SELECT editor FROM durations WHERE user_id = ?
    `
	if period == "today" {
		query += " AND DATE(time) = DATE('now','localtime')"
	}
	query += " GROUP BY editor ORDER BY SUM(duration) DESC LIMIT 1"

	var editor sql.NullString
	err := db.QueryRow(query, userID).Scan(&editor)
	if err != nil {
		return "", err
	}
	return editor.String, nil
}

// FetchTopProject fetches top project
func FetchTopProject(db *sql.DB, period string) (string, error) {
	query := `
        SELECT project FROM durations WHERE user_id = ?
    `
	if period == "today" {
		query += " AND DATE(time) = DATE('now','localtime')"
	}
	query += " GROUP BY project ORDER BY SUM(duration) DESC LIMIT 1"

	var project sql.NullString
	err := db.QueryRow(query, userID).Scan(&project)
	if err != nil {
		return "", err
	}
	return project.String, nil
}

// FetchTopOS fetches top OS
func FetchTopOS(db *sql.DB, period string) (string, error) {
	query := `
        SELECT machine FROM durations WHERE user_id = ?
    `
	if period == "today" {
		query += " AND DATE(time) = DATE('now','localtime')"
	}
	query += " GROUP BY machine ORDER BY SUM(duration) DESC LIMIT 1"

	var os sql.NullString
	err := db.QueryRow(query, userID).Scan(&os)
	if err != nil {
		return "", err
	}
	return os.String, nil
}

// FetchLanguagesCount
func FetchLanguagesCount(db *sql.DB, period string) (int64, error) {
	query := `
        SELECT COUNT(DISTINCT language) FROM durations WHERE user_id = ?
    `
	if period == "today" {
		query += " AND DATE(time) = DATE('now','localtime')"
	}

	var count sql.NullInt64
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count.Int64, nil
}

// FetchProjectsCount
func FetchProjectsCount(db *sql.DB, period string) (int64, error) {
	query := `
        SELECT COUNT(DISTINCT project) FROM durations WHERE user_id = ?
    `
	if period == "today" {
		query += " AND DATE(time) = DATE('now','localtime')"
	}

	var count sql.NullInt64
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count.Int64, nil
}

// FetchTimeByEditor returns a map of editor -> seconds
func FetchTimeByEditor(db *sql.DB, period string) (map[string]int64, error) {
	query := `
        SELECT editor, SUM(duration)/1000000000 AS seconds
        FROM durations
        WHERE user_id = ?
    `
	if period == "today" {
		query += " AND DATE(time) = DATE('now','localtime')"
	}
	query += " GROUP BY editor ORDER BY seconds DESC"

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var editor string
		var seconds sql.NullInt64
		if err := rows.Scan(&editor, &seconds); err != nil {
			return nil, err
		}
		result[editor] = seconds.Int64
	}
	return result, nil
}

// DisplayDashboard prints all stats
func DisplayDashboard(period string) {
	db, err := ConnectToSqlite()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	labelColor := color.New(color.FgBlue, color.Bold)
	// borderColor := color.New(color.FgCyan, color.Bold)
	valueColor := color.New(color.FgWhite)

	totalTime, _ := FetchTotalTime(db, period)
	dailyAvg, _ := FetchDailyAvg(db)
	topProject, _ := FetchTopProject(db, period)
	topEditor, _ := FetchTopEditor(db, period)
	topOS, _ := FetchTopOS(db, period)
	languages, _ := FetchLanguagesCount(db, period)
	projects, _ := FetchProjectsCount(db, period)
	editorTimes, _ := FetchTimeByEditor(db, period)

	// Print summary block
	// fmt.Println()

	// for range 15 {
	// 	fmt.Print("──")
	// }

	fmt.Println()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Print(cyan("╭─"))
	if period == "today" {
		fmt.Println(cyan("[  Daily Stats ] ──────────"))
	} else {
		fmt.Println(cyan("[  Stats ] ────────────"))
	}

	printStat("|"+"  Total Time    ", formatHrsMin(totalTime), labelColor, valueColor)
	printStat("|"+"  Daily Avg     ", formatHrsMin(dailyAvg), labelColor, valueColor)
	printStat("|"+"  Top Project   ", topProject, labelColor, valueColor)
	printStat("|"+"  Top Editor    ", topEditor, labelColor, valueColor)
	printStat("|"+"  Top OS        ", topOS, labelColor, valueColor)
	printStat("|"+"  Languages     ", fmt.Sprintf("%d", languages), labelColor, valueColor)
	printStat("|"+"  Projects      ", fmt.Sprintf("%d", projects), labelColor, valueColor)

	// Editor breakdown

	maxLen := 15
	for editor, seconds := range editorTimes {
		padding := strings.Repeat(" ", maxLen-len(editor))
		fmt.Printf(cyan("|")+"  %s%s%s\n",
			labelColor.Sprintf(editor),
			padding,
			valueColor.Sprintf(formatHrsMin(seconds)),
		)
	}

	fmt.Println(cyan("╰"))
	fmt.Println()
}
