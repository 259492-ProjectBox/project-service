package utils

import (
	"fmt"
	"regexp"
)

func FormatProjectID(semester, academicYear int, nextProjectNumber int) string {
	return fmt.Sprintf("P%d/%d/%04d", semester, academicYear%100, nextProjectNumber)
}

func IsValidProjectNumberFormat(projectNo string) error {
	matched, err := regexp.MatchString(`^P\d{3}-\d/\d{2}$`, projectNo)
	if err != nil {
		return fmt.Errorf("error validating project number format: %w", err)
	}
	if !matched {
		return fmt.Errorf("project number format is invalid: expected format P###-#/## (e.g., P001-2/66)")
	}
	return nil
}
