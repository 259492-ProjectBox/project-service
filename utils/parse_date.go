package utils

import (
	"fmt"
	"time"
)

func ParseDateTime(dateStr string) (time.Time, error) {
	const dateFormat = "02-01-2006" // Reference format for "DD-MM-YYYY"
	parsedDate, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	// Ensure the time is set to midnight explicitly (though it's the default)
	return parsedDate.Truncate(24 * time.Hour), nil
}

// parse from time.Time to string with format "DD-MM-YYYY"
func FormatDate(date time.Time) string {
	return date.Format("02-01-2006")
}
