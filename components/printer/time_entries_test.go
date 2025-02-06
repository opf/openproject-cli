package printer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func TestTimeEntry_Default(t *testing.T) {
	testingPrinter.Reset()

	timeEntry := models.TimeEntry{
		Id:          42,
		Comment:     "Glad it works!",
		Project:     "OpenProject",
		WorkPackage: "Create Go CLI",
		SpentOn:     "2009-11-10",
		Hours:       time.Hour + 45*time.Minute,
		Activity:    "Development",
	}

	expected := fmt.Sprintf("%s %s %s %s %s %s %s\n",
		printer.Red("#42"),
		printer.Green("DEVELOPMENT"),
		printer.Cyan("2009-11-10"),
		"1.75h",
		printer.Yellow("OpenProject"),
		printer.Cyan("Create Go CLI"),
		"Glad it works!")

	printer.TimeEntry(&timeEntry)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}

func TestTimeEntry_Ongoing(t *testing.T) {
	testingPrinter.Reset()

	timeEntry := models.TimeEntry{
		Id:          42,
		Comment:     "Glad it works!",
		Project:     "OpenProject",
		WorkPackage: "Create Go CLI",
		SpentOn:     "2009-11-10",
		Hours:       time.Hour + 45*time.Minute,
		Activity:    "Development",
		Ongoing:     true,
	}

	expected := fmt.Sprintf("%s %s %s %s %s %s %s (%s)\n",
		printer.Red("#42"),
		printer.Green("DEVELOPMENT"),
		printer.Cyan("2009-11-10"),
		"1.75h",
		printer.Yellow("OpenProject"),
		printer.Cyan("Create Go CLI"),
		"Glad it works!",
		printer.Yellow("ongoing"))

	printer.TimeEntry(&timeEntry)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}

func TestTimeEntryList(t *testing.T) {
	testingPrinter.Reset()

	timeEntry := []*models.TimeEntry{
		{
			Id:          8,
			Comment:     "Almost lost count!",
			Project:     "Mobile App",
			WorkPackage: "Sprint Meeting",
			SpentOn:     "2025-02-04",
			Hours:       4 * time.Hour,
			Activity:    "Accounting",
		},
		{
			Id:          61,
			Project:     "Frontend App",
			WorkPackage: "Rewrite in React",
			SpentOn:     "2024-11-08",
			Hours:       2*time.Hour + 15*time.Minute,
			Activity:    "Development",
		},
		{
			Id:          99,
			Comment:     "",
			Project:     "Customer Feedback",
			WorkPackage: "Stakeholder Testing",
			SpentOn:     "2025-01-02",
			Hours:       35 * time.Minute,
			Activity:    "Coordination",
			Ongoing:     true,
		},
	}

	expected := fmt.Sprintf("%s %s   %s %s %s        %s %s\n",
		printer.Red(" #8"),
		printer.Green("ACCOUNTING"),
		printer.Cyan("2025-02-04"),
		"4.00h",
		printer.Yellow("Mobile App"),
		printer.Cyan("Sprint Meeting"),
		"Almost lost count!")
	expected += fmt.Sprintf("%s %s  %s %s %s      %s\n",
		printer.Red("#61"),
		printer.Green("DEVELOPMENT"),
		printer.Cyan("2024-11-08"),
		"2.25h",
		printer.Yellow("Frontend App"),
		printer.Cyan("Rewrite in React"))
	expected += fmt.Sprintf("%s %s %s %s %s %s (%s)\n",
		printer.Red("#99"),
		printer.Green("COORDINATION"),
		printer.Cyan("2025-01-02"),
		"0.58h",
		printer.Yellow("Customer Feedback"),
		printer.Cyan("Stakeholder Testing"),
		printer.Yellow("ongoing"))

	printer.TimeEntryList(timeEntry)

	if testingPrinter.Result != expected {
		t.Errorf("Expected \n%s, but got \n%s", expected, testingPrinter.Result)
	}
}
