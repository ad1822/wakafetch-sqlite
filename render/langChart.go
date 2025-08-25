package render

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/ad1822/wakafetch-sqlite/sqlite"
	"github.com/fatih/color"
)

func RenderLangChart() []string {
	db, err := sqlite.ConnectToSqlite()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := sqlite.GetLangData(db)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	langCount := make(map[string]int)
	hasData := false
	for rows.Next() {
		var language string
		var count int
		if err := rows.Scan(&language, &count); err != nil {
			log.Fatal(err)
		}
		langCount[language] = count
		hasData = true
	}

	if !hasData {
		return []string{"No language data found."}
	}

	maxCount := 0
	for _, c := range langCount {
		if c > maxCount {
			maxCount = c
		}
	}

	type kv struct {
		Lang  string
		Count int
	}
	var sorted []kv
	for k, v := range langCount {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	barColor := color.New(color.FgBlue, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	var lines []string
	lines = append(lines, "")
	lines = append(lines, cyan("╭─[  Languages ] ──────────"))

	for _, item := range sorted {
		barLength := int(float64(item.Count) / float64(maxCount) * 30)
		bar := strings.Repeat("󱪿", barLength)
		line := fmt.Sprintf("%s %-10s %s",
			cyan("|"),     // border cyan
			item.Lang,     // plain text (or you can color language separately)
			barColor(bar), // only bar colored
		)
		lines = append(lines, line)
	}

	lines = append(lines, "")
	return lines
}
