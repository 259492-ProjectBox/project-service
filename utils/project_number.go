package utils

import (
	"fmt"
	"regexp"
)

func FormatProjectID(nextProjectNumber int) string {
	return fmt.Sprintf("P%04d", nextProjectNumber)
}

func IsValidProjectNumberFormat(projectNo string) error {
	matched, err := regexp.MatchString(`^P\d{4}-\d/\d{2}$`, projectNo)
	if err != nil {
		return fmt.Errorf("error validating project number format: %w", err)
	}
	if !matched {
		return fmt.Errorf("project number format is invalid: expected format P####-#/## (e.g., P0001-2/66)")
	}
	return nil
}
