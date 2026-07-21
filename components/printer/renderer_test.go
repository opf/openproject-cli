package printer_test

import (
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
)

func TestInitRenderer_UnknownFormat_PrintsError(t *testing.T) {
	testingPrinter.Reset()
	printer.InitRenderer("xml")
	if !strings.Contains(testingPrinter.ErrResult, "xml") {
		t.Errorf("expected error mentioning unknown format 'xml', got: %q", testingPrinter.ErrResult)
	}
}

func TestInitRenderer_KnownFormats_NoError(t *testing.T) {
	defer printer.InitRenderer("text")
	for _, format := range []string{"text", "json"} {
		testingPrinter.Reset()
		printer.InitRenderer(format)
		if strings.Contains(testingPrinter.ErrResult, "[ERROR]") {
			t.Errorf("unexpected error output for known format %q: %s", format, testingPrinter.ErrResult)
		}
	}
}
