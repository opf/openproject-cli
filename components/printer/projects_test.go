package printer_test

import (
	"os"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

var testingPrinter = &printer.TestingPrinter{}

func TestMain(m *testing.M) {
	printer.Init(testingPrinter)

	code := m.Run()
	os.Exit(code)
}

func TestProject(t *testing.T) {
	testingPrinter.Reset()

	expected := "[\033[31m#42\033[0m] \033[36mExample\033[0m\n"

	project := models.Project{Id: 42, Name: "Example"}

	printer.Project(&project)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}

func TestProjects(t *testing.T) {
	testingPrinter.Reset()

	expected := "[\033[31m#42\033[0m] \033[36mFoo\033[0m\n" +
		"[\033[31m#45\033[0m] \033[36mBar\033[0m\n" +
		"[\033[31m#123\033[0m] \033[36mBaz\033[0m\n"

	projects := []*models.Project{
		{Id: 42, Name: "Foo"},
		{Id: 45, Name: "Bar"},
		{Id: 123, Name: "Baz"},
	}

	printer.Projects(projects)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}
