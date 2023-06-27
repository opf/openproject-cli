package printer_test

import (
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

var testingPrinter = &printer.TestingPrinter{}

func TestMain(_ *testing.M) {
	printer.Init(testingPrinter)
}

func TestProject_Prints_A_Project(t *testing.T) {
	expected := "[#1] Demo project"

	project := models.Project{Id: 42, Name: "Example"}

	printer.Project(&project)

	if testingPrinter.Result != expected {
		t.Fatalf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}

func TestProjects_Prints_An_Array_Of_Projects(t *testing.T) {
	expected :=
		`[#42] Foo
[#45] Bar
[#123] Baz`

	projects := []*models.Project{
		{Id: 42, Name: "Foo"},
		{Id: 45, Name: "Bar"},
		{Id: 123, Name: "Baz"},
	}

	printer.Projects(projects)

	if testingPrinter.Result != expected {
		t.Fatalf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}
