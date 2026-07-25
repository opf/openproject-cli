package printer_test

import (
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
)

func TestWarning(t *testing.T) {
	testingPrinter.Reset()

	printer.Warning("watch out")

	if !strings.Contains(testingPrinter.ErrResult, "[WARNING]") {
		t.Errorf("warning missing [WARNING] tag: %q", testingPrinter.ErrResult)
	}
	if !strings.Contains(testingPrinter.ErrResult, "watch out") {
		t.Errorf("warning missing message: %q", testingPrinter.ErrResult)
	}
}
