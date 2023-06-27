package printer_test

import (
	"fmt"
	"github.com/opf/openproject-cli/components/routes"
	"strconv"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

//#8 TASK Setup conference website
//[In progress]
//Assignee: OpenProject Admin
//
//Open: https://openproject.local/work_packages/8
//
//# Awesome
//
//**Bacon ipsum** _d

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
