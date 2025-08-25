package render

import (
	"fmt"
	"log"
	"strings"

	"github.com/ad1822/wakafetch-sqlite/sqlite"
	"github.com/fatih/color"
)

func RenderDashboard(period string) []string {
	db, err := sqlite.ConnectToSqlite()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	labelColor := color.New(color.FgBlue, color.Bold)
	valueColor := color.New(color.FgWhite)
	cyan := color.New(color.FgCyan).SprintFunc()

	var lines []string

	totalTime, _ := sqlite.FetchTotalTime(db, period)
	dailyAvg, _ := sqlite.FetchDailyAvg(db)
	topProject, _ := sqlite.FetchTopProject(db, period)
	topEditor, _ := sqlite.FetchTopEditor(db, period)
	topOS, _ := sqlite.FetchTopOS(db, period)
	languages, _ := sqlite.FetchLanguagesCount(db, period)
	projects, _ := sqlite.FetchProjectsCount(db, period)
	editorTimes, _ := sqlite.FetchTimeByEditor(db, period)

	lines = append(lines, "") // empty line
	if period == "today" {
		lines = append(lines, cyan("╭─[  Daily Stats ] ──────────"))
	} else {
		lines = append(lines, cyan("╭─[  Stats ] ─────────────"))
	}

	lines = append(lines, sqlite.FormatStatLine("|  Total Time    ", sqlite.FormatHrsMin(totalTime), labelColor, valueColor))
	lines = append(lines, sqlite.FormatStatLine("|  Daily Avg     ", sqlite.FormatHrsMin(dailyAvg), labelColor, valueColor))
	lines = append(lines, sqlite.FormatStatLine("|  Top Project   ", topProject, labelColor, valueColor))
	lines = append(lines, sqlite.FormatStatLine("|  Top Editor    ", topEditor, labelColor, valueColor))
	lines = append(lines, sqlite.FormatStatLine("|  Top OS        ", topOS, labelColor, valueColor))
	lines = append(lines, sqlite.FormatStatLine("|  Languages     ", fmt.Sprintf("%d", languages), labelColor, valueColor))
	lines = append(lines, sqlite.FormatStatLine("|  Projects      ", fmt.Sprintf("%d", projects), labelColor, valueColor))

	// Editor breakdown
	maxLen := 15
	for editor, seconds := range editorTimes {
		padding := strings.Repeat(" ", maxLen-len(editor))
		lines = append(lines, fmt.Sprintf("|  %s%s%s",
			labelColor.Sprintf(editor),
			padding,
			valueColor.Sprintf(sqlite.FormatHrsMin(seconds)),
		))
	}

	lines = append(lines, cyan("╰"))

	return lines
}
