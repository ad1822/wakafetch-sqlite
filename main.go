package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	// "github.com/ad1822/wakafetch-sqlite/render"

	"github.com/ad1822/wakafetch-sqlite/render"
	"github.com/ad1822/wakafetch-sqlite/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	showStats := true
	showHeatmap := true

	if len(os.Args) > 1 {
		arg := strings.ToLower(os.Args[1])
		switch arg {
		case "stats":
			showHeatmap = false
		case "heatmap":
			showStats = false
		case "all":
			// default, show everything
		default:
			fmt.Println("Usage: app [all|stats|heatmap|langchart]")
			return
		}
	}

	// Fetch heatmap data once
	to := time.Now().UTC()
	from := to.AddDate(-1, -3, 0)
	activities, err := sqlite.FetchDataForHeatMap(from, to)
	if err != nil {
		panic(err)
	}

	// Render side-by-side stats and language chart
	if showStats {
		period := "all"
		left := render.RenderDashboard(period)
		right := render.RenderLangChart()

		maxLines := len(left)
		if len(right) > maxLines {
			maxLines = len(right)
		}

		for i := 0; i < maxLines; i++ {
			l := ""
			r := ""
			if i < len(left) {
				l = left[i]
			}
			if i < len(right) {
				r = right[i]
			}

			// Align chart title separately if needed
			if strings.Contains(r, "[  Languages ]") {
				fmt.Printf("%-47s %s\n", l, r)
			} else {
				fmt.Printf("%-70s %s\n", l, r)
			}
		}
	}

	// Render heatmap
	if showHeatmap {
		render.RenderHeatmap(activities)
	}
}
