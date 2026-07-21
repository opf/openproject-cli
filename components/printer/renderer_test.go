package printer_test

import (
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
)

func TestInitRenderer_UnknownFormat_PrintsError(t *testing.T) {
	testingPrinter.Reset()
	err := printer.InitRenderer("xml")
	if err == nil {
		t.Fatal("expected unknown format to return an error")
	}
	if !strings.Contains(testingPrinter.ErrResult, "xml") {
		t.Errorf("expected error mentioning unknown format 'xml', got: %q", testingPrinter.ErrResult)
	}
	if count := strings.Count(testingPrinter.ErrResult, "[ERROR]"); count != 1 {
		t.Errorf("error diagnostic count = %d, want 1", count)
	}
	if testingPrinter.Result != "" {
		t.Errorf("stdout = %q, want empty", testingPrinter.Result)
	}
}

func TestInitRenderer_KnownFormats_NoError(t *testing.T) {
	defer printer.InitRenderer("text")
	for _, format := range []string{"text", "json"} {
		testingPrinter.Reset()
		if err := printer.InitRenderer(format); err != nil {
			t.Errorf("InitRenderer(%q) returned an error: %v", format, err)
		}
		if strings.Contains(testingPrinter.ErrResult, "[ERROR]") {
			t.Errorf("unexpected error output for known format %q: %s", format, testingPrinter.ErrResult)
		}
	}
}
