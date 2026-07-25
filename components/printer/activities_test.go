package printer_test

import (
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

// TestActivities_UserLookup_FirstUser reproduces the sort.Search bug:
// the old predicate (Id == target) is not monotone, and len(users)-1 as the
// upper bound means the last user is never searched. Together they cause
// sort.Search to return the wrong index when the target is the first user.
func TestActivities_UserLookup_FirstUser(t *testing.T) {
	testingPrinter.Reset()
	printer.InitRenderer("text")

	users := []*models.User{
		{Id: 1, Name: "Alice"},
		{Id: 3, Name: "Bob"},
		{Id: 5, Name: "Charlie"},
	}
	activity := &models.Activity{
		Id:        42,
		UserId:    1, // Alice — first in the slice
		UpdatedAt: "2024-01-01",
	}

	printer.Activities([]*models.Activity{activity}, users)

	if !strings.Contains(testingPrinter.Result, printer.Green("Alice")) {
		t.Errorf("expected user 'Alice' in output, got: %q", testingPrinter.Result)
	}
}
