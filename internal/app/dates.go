package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func weekStartOf(t time.Time) time.Time {
	year, month, day := t.Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	weekday := int(start.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return start.AddDate(0, 0, -(weekday - 1))
}

func weekLabel(start time.Time, span int) string {
	if span < 0 {
		return "All Weeks"
	}
	if span <= 1 {
		_, weekNum := start.ISOWeek()
		end := start.AddDate(0, 0, 6)
		return "Week " + strconv.Itoa(weekNum) + " - " + start.Format("Jan 2") + " - " + end.Format("Jan 2, 2006")
	}
	end := start.AddDate(0, 0, span*7-1)
	_, startWeek := start.ISOWeek()
	_, endWeek := end.ISOWeek()
	return fmt.Sprintf("Weeks %d-%d - %s - %s", startWeek, endWeek, start.Format("Jan 2"), end.Format("Jan 2, 2006"))
}

func parseExamDate(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	layouts := []string{
		"Jan 2, 2006 @ 15:04",
		"Jan 2, 2006",
		"2006-01-02 15:04",
		"2006-01-02",
		"2.1.2006 15:04",
		"2.1.2006",
		"02.01.2006 15:04",
		"02.01.2006",
		"2/1/2006 15:04",
		"2/1/2006",
		"02/01/2006 15:04",
		"02/01/2006",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func parseTodoDate(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	if t, err := time.ParseInLocation("2006-01-02", value, time.Local); err == nil {
		return t, true
	}
	return time.Time{}, false
}
