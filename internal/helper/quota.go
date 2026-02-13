package helper

import (
	"time"
)

func CalculateNextResetTime(now time.Time, duration int64, unit string) int64 {
	if unit == "day" {
		days := duration / 86400
		if days <= 0 {
			days = 1
		}
		// Calculate next reset time as the next occurrence of midnight UTC that's X days away
		year, month, day := now.UTC().Date()
		todayMidnight := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		nextReset := todayMidnight.AddDate(0, 0, int(days))
		
		// Safety check: if for some reason nextReset is still not after now, increment
		for !nextReset.After(now.UTC()) {
			nextReset = nextReset.AddDate(0, 0, int(days))
		}

		return nextReset.Unix()
	}
	// For minutes and hours, use relative timing as before
	return now.Unix() + duration
}

func IsAlignedToMidnight(timestamp int64) bool {
	if timestamp == 0 {
		return true
	}
	t := time.Unix(timestamp, 0).UTC()
	return t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0
}
