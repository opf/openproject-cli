package timeentry_test

import (
	"strings"
	"testing"

	"github.com/opf/openproject-cli/cmd/timeentry"
)

func TestRootCmd_UsesHyphenatedName(t *testing.T) {
	use := timeentry.RootCmd.Use
	if !strings.HasPrefix(use, "time-entry") {
		t.Errorf("expected command name 'time-entry', got %q", use)
	}
}
