package printer_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/routes"
	"github.com/opf/openproject-cli/models"
)

func TestWorkPackage(t *testing.T) {
	testingPrinter.Reset()

	workPackage := models.WorkPackage{
		Id:          42,
		Subject:     "Test",
		Type:        "TASK",
		Assignee:    models.Principal{Name: "Aaron"},
		Status:      "New",
		Description: "This is an example.",
		LockVersion: 0,
	}

	var expected string

	idString := "#" + strconv.FormatUint(workPackage.Id, 10)
	expected += fmt.Sprintf("%s %s %s\n", printer.Red(idString), printer.Green(workPackage.Type), printer.Cyan(workPackage.Subject))
	expected += fmt.Sprintf("[%s]\n", printer.Yellow(workPackage.Status))
	expected += fmt.Sprintf("Assignee: %s\n\n", workPackage.Assignee.Name)
	expected += fmt.Sprintf("Open: %s\n\n", routes.WorkPackageUrl(&workPackage))
	expected += fmt.Sprintln(workPackage.Description)

	printer.WorkPackage(&workPackage)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}

func TestWorkPackages(t *testing.T) {
	testingPrinter.Reset()

	workPackages := []*models.WorkPackage{
		{
			Id:          42,
			Subject:     "Test 1",
			Type:        "PHASE",
			Assignee:    models.Principal{Name: "Obi-Wan"},
			Status:      "In progress",
			Description: "This is one example.",
			LockVersion: 0,
		},
		{
			Id:          43,
			Subject:     "Test 2",
			Type:        "TASK",
			Assignee:    models.Principal{Name: "Anakin"},
			Status:      "New",
			Description: "This is another example.",
			LockVersion: 0,
		},
	}

	var expected string

	idString := "#" + strconv.FormatUint(workPackages[0].Id, 10)
	expected += fmt.Sprintf("%s %s [%s] %s\n", printer.Red(idString), printer.Green(workPackages[0].Type), printer.Yellow(workPackages[0].Status), printer.Cyan(workPackages[0].Subject))
	idString = "#" + strconv.FormatUint(workPackages[1].Id, 10)
	expected += fmt.Sprintf("%s %s [%s]         %s\n", printer.Red(idString), printer.Green(workPackages[1].Type+" "), printer.Yellow(workPackages[1].Status), printer.Cyan(workPackages[1].Subject))

	printer.WorkPackages(workPackages)

	if testingPrinter.Result != expected {
		t.Errorf("\nExpected:\n%sbut got:\n%s", expected, testingPrinter.Result)
	}
}
