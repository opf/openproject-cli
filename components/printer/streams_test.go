package printer_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
)

// Diagnostics must go to standard error so that machine-readable output
// (e.g. --format json) on standard out stays parseable.
func TestDiagnosticsGoToStderr(t *testing.T) {
	diagnostics := map[string]func(){
		"Info":      func() { printer.Info("progress") },
		"Done":      func() { printer.Done() },
		"Error":     func() { printer.Error(errors.New("boom")) },
		"ErrorText": func() { printer.ErrorText("boom") },
		"Warning":   func() { printer.Warning("boom") },
		"Debug":     func() { printer.Debug(true, "boom") },
	}

	for name, call := range diagnostics {
		testingPrinter.Reset()
		call()

		if len(testingPrinter.ErrResult) == 0 {
			t.Errorf("%s: expected output on stderr, got none", name)
		}
		if len(testingPrinter.Result) > 0 {
			t.Errorf("%s: unexpected output on stdout: %q", name, testingPrinter.Result)
		}
	}
}

func TestOutputGoesToStdout(t *testing.T) {
	testingPrinter.Reset()

	printer.Output("data")

	if !strings.Contains(testingPrinter.Result, "data") {
		t.Errorf("expected data on stdout, got: %q", testingPrinter.Result)
	}
	if len(testingPrinter.ErrResult) > 0 {
		t.Errorf("unexpected output on stderr: %q", testingPrinter.ErrResult)
	}
}
