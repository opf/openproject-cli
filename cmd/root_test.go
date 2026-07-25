package cmd

import (
	"errors"
	"strings"
	"testing"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
)

func TestInvalidOutputFormatStopsCommand(t *testing.T) {
	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)
	outputFormat = "yaml"
	t.Cleanup(func() {
		outputFormat = "text"
		_ = printer.InitRenderer("text")
	})

	err := rootCmd.PersistentPreRunE(rootCmd, nil)
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("PersistentPreRunE error = %v, want ErrHandled", err)
	}
	if count := strings.Count(testingPrinter.ErrResult, "[ERROR]"); count != 1 {
		t.Errorf("error diagnostic count = %d, want 1; stderr: %q", count, testingPrinter.ErrResult)
	}
	if testingPrinter.Result != "" {
		t.Errorf("stdout = %q, want empty", testingPrinter.Result)
	}
}
