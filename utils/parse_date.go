package utils

import (
	"fmt"
	"time"
)

func ParseDateTime(dateStr string) (time.Time, error) {
	const dateFormat = "2006-01-02" // Reference format for "YYYY-MM-DD"
	parsedDate, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	// Ensure the time is set to midnight explicitly (though it's the default)
	return parsedDate.Truncate(24 * time.Hour), nil
}
