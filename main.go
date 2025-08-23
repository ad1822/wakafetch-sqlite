package main

import (
	"time"

	"github.com/ad1822/wakafetch-sqlite/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	to := time.Now().UTC()
	from := to.AddDate(0, 0, -100)
	sqlite.FetchDataForHeatMap(from, to)
}
