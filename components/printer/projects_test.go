package printer_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func TestProject(t *testing.T) {
	testingPrinter.Reset()

	project := models.Project{Id: 42, Name: "Example"}

	idString := "#" + strconv.FormatUint(project.Id, 10)
	expected := fmt.Sprintf("%s %s\n", printer.Red(idString), printer.Cyan(project.Name))

	printer.Project(&project)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}

func TestProjects(t *testing.T) {
	testingPrinter.Reset()

	projects := []*models.Project{
		{Id: 42, Name: "Foo"},
		{Id: 45, Name: "Bar"},
		{Id: 123, Name: "Baz"},
	}

	expected := common.Reduce[*models.Project, string](
		projects,
		func(state string, project *models.Project) string {
			idString := "#" + strconv.FormatUint(project.Id, 10)

			return state + fmt.Sprintf("%s %s\n", printer.Red(idString), printer.Cyan(project.Name))
		},
		"")

	printer.Projects(projects)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}
