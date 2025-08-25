package render

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/ad1822/wakafetch-sqlite/sqlite"
	"github.com/fatih/color"
)

func DisplayLangChart() {
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

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if !hasData {
		fmt.Println("No language data found.")
		return
	}

	// Find max count
	maxCount := 0
	for _, c := range langCount {
		if c > maxCount {
			maxCount = c
		}
	}

	if maxCount == 0 {
		fmt.Println("No activity counts available.")
		return
	}

	// Convert map to slice and sort by count
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
	// Draw the chart
	fmt.Println()

	fmt.Print(cyan("╭─"))
	fmt.Println(cyan("[  Daily Stats ] ──────────"))
	for _, item := range sorted {
		barLength := int(float64(item.Count) / float64(maxCount) * 30)
		bar := strings.Repeat("█", barLength)
		line := fmt.Sprintf("| %-10s %s\n", item.Lang, bar) // format string
		fmt.Print(barColor(line))
	}

	fmt.Println()
}
