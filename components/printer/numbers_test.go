package printer_test

import (
	"fmt"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
)

func TestNumber(t *testing.T) {
	testingPrinter.Reset()

	expected := fmt.Sprintf("%s\n", printer.Cyan("5"))
	printer.Number(5)

	if testingPrinter.Result != expected {
		t.Errorf("\nExpected:\n%sbut got:\n%s", expected, testingPrinter.Result)
	}
}
