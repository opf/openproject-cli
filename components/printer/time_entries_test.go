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
		SpentOn:     time.Date(2009, time.November, 10, 0, 0, 0, 0, time.UTC),
		Hours:       time.Hour + 45*time.Minute,
		Activity:    "Development",
	}

	expected := fmt.Sprintf("%s %s %s %s %s %s %s\n",
		printer.Red("#42"),
		printer.Green("DEVELOPMENT"),
		printer.Cyan("Tue Nov 10"),
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
		SpentOn:     time.Date(2009, time.November, 10, 0, 0, 0, 0, time.UTC),
		Hours:       time.Hour + 45*time.Minute,
		Activity:    "Development",
		Ongoing:     true,
	}

	expected := fmt.Sprintf("%s %s %s %s %s %s %s (%s)\n",
		printer.Red("#42"),
		printer.Green("DEVELOPMENT"),
		printer.Cyan("Tue Nov 10"),
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
			Hours:       4 * time.Hour,
			SpentOn:     time.Date(2025, time.February, 4, 0, 0, 0, 0, time.UTC),
			Activity:    "Accounting",
		},
		{
			Id:          61,
			Project:     "Frontend App",
			WorkPackage: "Rewrite in React",
			Hours:       2*time.Hour + 15*time.Minute,
			SpentOn:     time.Date(2024, time.November, 8, 0, 0, 0, 0, time.UTC),
			Activity:    "Development",
		},
		{
			Id:          99,
			Comment:     "",
			Project:     "Customer Feedback",
			WorkPackage: "Stakeholder Testing",
			Hours:       35 * time.Minute,
			SpentOn:     time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
			Activity:    "Coordination",
			Ongoing:     true,
		},
	}

	expected := fmt.Sprintf("%s %s   %s %s %s        %s %s\n",
		printer.Red(" #8"),
		printer.Green("ACCOUNTING"),
		printer.Cyan("Tue Feb  4"),
		"4.00h",
		printer.Yellow("Mobile App"),
		printer.Cyan("Sprint Meeting"),
		"Almost lost count!")
	expected += fmt.Sprintf("%s %s  %s %s %s      %s\n",
		printer.Red("#61"),
		printer.Green("DEVELOPMENT"),
		printer.Cyan("Fri Nov  8"),
		"2.25h",
		printer.Yellow("Frontend App"),
		printer.Cyan("Rewrite in React"))
	expected += fmt.Sprintf("%s %s %s %s %s %s (%s)\n",
		printer.Red("#99"),
		printer.Green("COORDINATION"),
		printer.Cyan("Thu Jan  2"),
		"0.58h",
		printer.Yellow("Customer Feedback"),
		printer.Cyan("Stakeholder Testing"),
		printer.Yellow("ongoing"))

	printer.TimeEntryList(timeEntry)

	if testingPrinter.Result != expected {
		t.Errorf("Expected \n%s, but got \n%s", expected, testingPrinter.Result)
	}
}
