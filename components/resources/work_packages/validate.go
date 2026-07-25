package work_packages

import (
	"fmt"
	"regexp"
)

// matches plain numeric IDs like 12345
var numericIdPattern = regexp.MustCompile(`^\d+$`)

// matches project-based identifiers like PROJ-123: uppercase letter start, up to 10 chars
// (letters/digits/underscores), hyphen, numeric sequence — per OpenProject identifier rules
var semanticIdPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]{0,9}-\d+$`)

func ValidateIdentifier(id string) error {
	if numericIdPattern.MatchString(id) || semanticIdPattern.MatchString(id) {
		return nil
	}
	return fmt.Errorf("'%s' is an invalid work package identifier: must be a numeric ID (e.g. '12345') or a project-based identifier (e.g. 'PROJ-123')", id)
}
