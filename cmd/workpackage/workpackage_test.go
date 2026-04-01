package workpackage_test

import (
	"strings"
	"testing"

	"github.com/opf/openproject-cli/cmd/workpackage"
)

func TestRootCmd_UsesHyphenatedName(t *testing.T) {
	use := workpackage.RootCmd.Use
	if !strings.HasPrefix(use, "work-package") {
		t.Errorf("expected command name 'work-package', got %q", use)
	}
}
