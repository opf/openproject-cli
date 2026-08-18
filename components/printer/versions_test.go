package printer_test

import (
	"fmt"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func TestVersions(t *testing.T) {
	testingPrinter.Reset()

	versions := []*models.Version{
		{Id: 2, Name: "13.0"},
		{Id: 4, Name: "13.1"},
		{Id: 43, Name: "45.5"},
	}

	expected := fmt.Sprintf("%s %s\n", printer.Red(" #2"), printer.Cyan("13.0"))
	expected += fmt.Sprintf("%s %s\n", printer.Red(" #4"), printer.Cyan("13.1"))
	expected += fmt.Sprintf("%s %s\n", printer.Red("#43"), printer.Cyan("45.5"))

	printer.Versions(versions)

	if testingPrinter.Result != expected {
		t.Errorf("\nExpected:\n%sbut got:\n%s", expected, testingPrinter.Result)
	}
}
