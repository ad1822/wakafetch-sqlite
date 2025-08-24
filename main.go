package main

import (
	"time"

	// "github.com/ad1822/wakafetch-sqlite/render"

	"github.com/ad1822/wakafetch-sqlite/render"
	"github.com/ad1822/wakafetch-sqlite/sqlite"
	heatmap "github.com/ad1822/wakafetch-sqlite/ui"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	to := time.Now().UTC()
	from := to.AddDate(-1, -3, 0)
	activities, err := sqlite.FetchDataForHeatMap(from, to)
	if err != nil {
		panic(err)
	}

	render.DisplayDashboard("today")
	heatmap.RenderHeatmap(activities)

	// sqlite.DisplayDashboard("all")
}
