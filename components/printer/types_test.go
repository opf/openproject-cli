package printer_test

import (
	"fmt"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func TestTypes(t *testing.T) {
	testingPrinter.Reset()

	types := []*models.Type{
		{Id: 42, Name: "Lightsaber"},
		{Id: 1337, Name: "Leet"},
	}

	var expected string
	expected += fmt.Sprintf("%s %s\n", printer.Red(fmt.Sprintf("  #%d", types[0].Id)), printer.Cyan(types[0].Name))
	expected += fmt.Sprintf("%s %s\n", printer.Red(fmt.Sprintf("#%d", types[1].Id)), printer.Cyan(types[1].Name))

	printer.Types(types)

	if testingPrinter.Result != expected {
		t.Errorf("Expected \n%s\tbut got \n%s", expected, testingPrinter.Result)
	}
}
