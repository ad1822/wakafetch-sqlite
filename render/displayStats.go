package render

import (
	"fmt"
	"log"
	"strings"

	"github.com/ad1822/wakafetch-sqlite/sqlite"
	"github.com/fatih/color"
)

func DisplayDashboard(period string) {
	db, err := sqlite.ConnectToSqlite()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	labelColor := color.New(color.FgBlue, color.Bold)
	valueColor := color.New(color.FgWhite)

	totalTime, _ := sqlite.FetchTotalTime(db, period)
	dailyAvg, _ := sqlite.FetchDailyAvg(db)
	topProject, _ := sqlite.FetchTopProject(db, period)
	topEditor, _ := sqlite.FetchTopEditor(db, period)
	topOS, _ := sqlite.FetchTopOS(db, period)
	languages, _ := sqlite.FetchLanguagesCount(db, period)
	projects, _ := sqlite.FetchProjectsCount(db, period)
	editorTimes, _ := sqlite.FetchTimeByEditor(db, period)

	fmt.Println()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Print(cyan("╭─"))
	if period == "today" {
		fmt.Println(cyan("[  Daily Stats ] ──────────"))
	} else {
		fmt.Println(cyan("[  Stats ] ────────────"))
	}

	sqlite.PrintStat("|"+"  Total Time    ", sqlite.FormatHrsMin(totalTime), labelColor, valueColor)
	sqlite.PrintStat("|"+"  Daily Avg     ", sqlite.FormatHrsMin(dailyAvg), labelColor, valueColor)
	sqlite.PrintStat("|"+"  Top Project   ", topProject, labelColor, valueColor)
	sqlite.PrintStat("|"+"  Top Editor    ", topEditor, labelColor, valueColor)
	sqlite.PrintStat("|"+"  Top OS        ", topOS, labelColor, valueColor)
	sqlite.PrintStat("|"+"  Languages     ", fmt.Sprintf("%d", languages), labelColor, valueColor)
	sqlite.PrintStat("|"+"  Projects      ", fmt.Sprintf("%d", projects), labelColor, valueColor)

	// Editor breakdown

	maxLen := 15
	for editor, seconds := range editorTimes {
		padding := strings.Repeat(" ", maxLen-len(editor))
		fmt.Printf(cyan("|")+"  %s%s%s\n",
			labelColor.Sprintf(editor),
			padding,
			valueColor.Sprintf(sqlite.FormatHrsMin(seconds)),
		)
	}

	fmt.Println(cyan("╰"))
}
