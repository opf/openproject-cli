package projects

import (
	"fmt"
	"regexp"

	"github.com/opf/openproject-cli/components/errors"
)

// Dashes are allowed for old project identifiers; semantic identifiers only allow letters, numbers and underscores.
var invalidIdentifierChars = regexp.MustCompile(`[^a-zA-Z0-9\-_]`)

func ValidateIdentifier(identifier string) error {
	if identifier == "" {
		return errors.Custom("project identifier cannot be empty")
	}
	if invalidIdentifierChars.MatchString(identifier) {
		return errors.Custom(fmt.Sprintf(
			"invalid project %q: only letters, numbers, - and _ are allowed",
			identifier,
		))
	}
	return nil
}
