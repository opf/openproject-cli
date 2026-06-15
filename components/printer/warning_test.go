package printer_test

import (
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
)

func TestWarning(t *testing.T) {
	testingPrinter.Reset()

	printer.Warning("watch out")

	if !strings.Contains(testingPrinter.Result, "[WARNING]") {
		t.Errorf("warning missing [WARNING] tag: %q", testingPrinter.Result)
	}
	if !strings.Contains(testingPrinter.Result, "watch out") {
		t.Errorf("warning missing message: %q", testingPrinter.Result)
	}
}
