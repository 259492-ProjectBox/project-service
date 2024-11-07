package utils

import (
	"fmt"
	"regexp"
)

func FormatProjectNumber(sectionID, semester, academicYear int) string {
	return fmt.Sprintf("P%03d-%d/%02d", sectionID, semester, academicYear%100)
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
