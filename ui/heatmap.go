package heatmap

import (
	"fmt"
	"time"

	"github.com/ad1822/wakafetch-sqlite/types"
)

const heatmapChar = "■"
const highlight = "\x1b[38;2;0;%v;0m"
const borderColor = "\x1b[38;5;245m" // gray

func RenderHeatmap(activities []types.DailyActivity) {
	if len(activities) == 0 {
		fmt.Println("No activity data.")
		return
	}

	activityMap := make(map[string]int64)
	var minDate, maxDate time.Time
	var maxCount int64

	for i, a := range activities {
		key := a.Date.Format("2006-01-02")
		activityMap[key] = a.Count
		if i == 0 || a.Date.Before(minDate) {
			// minDate = a.Date
			// For Extending Lower Limit
			minDate = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
		}
		if i == 0 || a.Date.After(maxDate) {
			// maxDate = a.Date
			// For Extending Upper Limit
			maxDate = time.Date(2025, time.October, 30, 0, 0, 0, 0, time.UTC)
		}
		if a.Count > maxCount {
			maxCount = a.Count
		}
	}

	// Align min/max dates to Monday/Sunday
	for minDate.Weekday() != time.Monday {
		minDate = minDate.AddDate(0, 0, -1)
	}
	for maxDate.Weekday() != time.Sunday {
		maxDate = maxDate.AddDate(0, 0, 1)
	}

	// Build grid: weeks as columns
	var grid [][]int64
	date := minDate
	for !date.After(maxDate) {
		weekdayIndex := int(date.Weekday()) - 1
		if weekdayIndex < 0 {
			weekdayIndex = 6 // Sunday
		}

		if len(grid) == 0 || weekdayIndex == 0 {
			grid = append(grid, make([]int64, 7))
		}

		weekIndex := len(grid) - 1
		grid[weekIndex][weekdayIndex] = activityMap[date.Format("2006-01-02")]
		date = date.AddDate(0, 0, 1)
	}

	cols := len(grid)

	// title := "WakaAPI "
	// Top border
	fmt.Print(borderColor + "╭")
	for i := 0; i < (cols + 1); i++ {
		fmt.Print("──")
	}
	fmt.Println("╮\x1b[0m")

	// Print vertical layout
	for day := 0; day < 7; day++ {
		fmt.Print(borderColor + "│ " + "\x1b[0m")

		for week := 0; week < len(grid); week++ {
			count := grid[week][day]
			green := 0
			if maxCount > 0 && count > 0 {
				green = int(float64(count) / float64(maxCount) * 255)
			}

			fmt.Printf(highlight+heatmapChar+"\x1b[0m", green)

			fmt.Print(" ")
		}

		fmt.Println(borderColor + " │\x1b[0m")
	}

	fmt.Print(borderColor + "╰")
	for i := 0; i < cols+1; i++ {
		fmt.Print("──")
	}
	fmt.Println("╯\x1b[0m")
}
