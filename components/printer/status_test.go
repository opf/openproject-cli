package printer_test

import (
	"fmt"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func TestStatus_Default(t *testing.T) {
	testingPrinter.Reset()

	status := models.Status{
		Id:        42,
		Name:      "My Default",
		IsDefault: true,
	}

	expected := fmt.Sprintf("%s %s (%s)\n", printer.Red("#42"), printer.Cyan("My Default"), printer.Yellow("default"))

	printer.Status(&status)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}

func TestStatus_NoDefault(t *testing.T) {
	testingPrinter.Reset()

	status := models.Status{
		Id:        42,
		Name:      "Another Status",
		IsDefault: false,
	}

	expected := fmt.Sprintf("%s %s\n", printer.Red("#42"), printer.Cyan("Another Status"))

	printer.Status(&status)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}

func TestStatusList(t *testing.T) {
	testingPrinter.Reset()

	status := []*models.Status{
		{Id: 2, Name: "First"},
		{Id: 45, Name: "Second"},
		{Id: 123, Name: "Third", IsDefault: true},
	}

	expected := fmt.Sprintf("%s %s\n", printer.Red("  #2"), printer.Cyan("First"))
	expected += fmt.Sprintf("%s %s\n", printer.Red(" #45"), printer.Cyan("Second"))
	expected += fmt.Sprintf("%s %s (%s)\n", printer.Red("#123"), printer.Cyan("Third"), printer.Yellow("default"))

	printer.StatusList(status)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}
