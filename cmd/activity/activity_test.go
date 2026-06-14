package activity_test

import (
	"strings"
	"testing"

	"github.com/opf/openproject-cli/cmd/activity"
)

func TestRootCmd_UsesSingularNoun(t *testing.T) {
	use := activity.RootCmd.Use
	if !strings.HasPrefix(use, "activity") {
		t.Errorf("expected command name 'activity', got %q", use)
	}
}
