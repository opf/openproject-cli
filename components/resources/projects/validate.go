package projects

import (
	"fmt"
	"regexp"

	"github.com/opf/openproject-cli/components/errors"
)

var invalidIdentifierChars = regexp.MustCompile(`[^a-zA-Z0-9\-_+]`)

func ValidateIdentifier(identifier string) error {
	if identifier == "" {
		return errors.Custom("project identifier cannot be empty")
	}
	if invalidIdentifierChars.MatchString(identifier) {
		return errors.Custom(fmt.Sprintf(
			"invalid project %q: only letters, numbers, -, _, and + are allowed",
			identifier,
		))
	}
	return nil
}
