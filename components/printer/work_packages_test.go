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
		Id:       42,
		Subject:  "Test",
		Type:     "TASK",
		Assignee: "Aaron",
		Status:   "New",
		Description: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam non est sed nunc euismod luctus. Donec vehicula scelerisque efficitur. Nunc arcu ligula, dictum maximus consequat id, tincidunt vitae augue. Vestibulum ut tellus id nisi faucibus efficitur id eu tortor. Vestibulum sed vehicula turpis, sit amet eleifend massa. Proin eu quam justo. Nulla id libero sit amet turpis venenatis mollis. Quisque iaculis lectus non ligula faucibus, ut pellentesque velit sodales. Vivamus nibh est, molestie at laoreet nec, lacinia porttitor nisl. Nulla eget urna in enim porttitor tempus. Nullam velit nunc, ultrices eget molestie vitae, tincidunt vitae felis.

Integer augue purus, mollis a vestibulum quis, sagittis ac lacus. Ut vitae tempor tellus. Cras neque turpis, malesuada nec tincidunt vel, mattis vel dolor. Ut hendrerit magna ac suscipit convallis. Ut quis nisi vel metus facilisis sagittis eu eget orci. Pellentesque laoreet metus vitae nulla fringilla, sed lacinia sem laoreet. Maecenas velit erat, luctus ac metus eget, hendrerit tincidunt dolor. Cras mattis orci sem, sed convallis arcu venenatis nec. Donec imperdiet mattis ante, quis euismod lorem viverra ac. Pellentesque in efficitur magna, at ullamcorper ipsum. Vivamus vulputate, tellus et blandit mollis, elit nisl posuere dui, nec molestie metus arcu in lectus. Vivamus eget congue libero, ut congue dolor. Interdum et malesuada fames ac ante ipsum primis in faucibus. Suspendisse blandit, velit quis euismod tincidunt, nunc lectus rutrum nisi, et commodo enim ligula nec sem. Pellentesque nec tincidunt sapien.`,
		LockVersion: 0,
	}

	var expectedDescription = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam non est sed nunc
euismod luctus. Donec vehicula scelerisque efficitur. Nunc arcu ligula, dictum
maximus consequat id, tincidunt vitae augue. Vestibulum ut tellus id nisi
faucibus efficitur id eu tortor. Vestibulum sed vehicula turpis, sit amet
eleifend massa. Proin eu quam justo. Nulla id libero sit amet turpis venenatis
mollis. Quisque iaculis lectus non ligula faucibus, ut pellentesque velit
sodales. Vivamus nibh est, molestie at laoreet nec, lacinia porttitor nisl.
Nulla eget urna in enim porttitor tempus. Nullam velit nunc, ultrices eget
molestie vitae, tincidunt vitae felis.

Integer augue purus, mollis a vestibulum quis, sagittis ac lacus. Ut vitae
tempor tellus. Cras neque turpis, malesuada nec tincidunt vel, mattis vel dolor.
Ut hendrerit magna ac suscipit convallis. Ut quis nisi vel metus facilisis
sagittis eu eget orci. Pellentesque laoreet metus vitae nulla fringilla, sed
lacinia sem laoreet. Maecenas velit erat, luctus ac metus eget, hendrerit
tincidunt dolor. Cras mattis orci sem, sed convallis arcu venenatis nec. Donec
imperdiet mattis ante, quis euismod lorem viverra ac. Pellentesque in efficitur
magna, at ullamcorper ipsum. Vivamus vulputate, tellus et blandit mollis, elit
nisl posuere dui, nec molestie metus arcu in lectus. Vivamus eget congue libero,
ut congue dolor. Interdum et malesuada fames ac ante ipsum primis in faucibus.
Suspendisse blandit, velit quis euismod tincidunt, nunc lectus rutrum nisi, et
commodo enim ligula nec sem. Pellentesque nec tincidunt sapien.`

	var expected string

	idString := "#" + strconv.FormatUint(workPackage.Id, 10)
	expected += fmt.Sprintf("%s %s %s\n", printer.Red(idString), printer.Green(workPackage.Type), printer.Cyan(workPackage.Subject))
	expected += fmt.Sprintf("[%s]\n", printer.Yellow(workPackage.Status))
	expected += fmt.Sprintf("Assignee: %s\n\n", workPackage.Assignee)
	expected += fmt.Sprintf("Open: %s\n\n", routes.WorkPackageUrl(&workPackage))
	expected += fmt.Sprintln(expectedDescription)

	printer.WorkPackage(&workPackage)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}

func TestWorkPackage_Assignee_With_Empty_String(t *testing.T) {
	testingPrinter.Reset()

	workPackage := models.WorkPackage{
		Id:          42,
		Subject:     "Test",
		Type:        "TASK",
		Assignee:    "",
		Status:      "New",
		Description: "This is an example.",
		LockVersion: 0,
	}

	var expected string

	idString := "#" + strconv.FormatUint(workPackage.Id, 10)
	expected += fmt.Sprintf("%s %s %s\n", printer.Red(idString), printer.Green(workPackage.Type), printer.Cyan(workPackage.Subject))
	expected += fmt.Sprintf("[%s]\n", printer.Yellow(workPackage.Status))
	expected += fmt.Sprintf("Assignee: %s\n\n", "-")
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
			Assignee:    "Obi-Wan",
			Status:      "In progress",
			Description: "This is one example.",
			LockVersion: 0,
		},
		{
			Id:          43,
			Subject:     "Test 2",
			Type:        "TASK",
			Assignee:    "Anakin",
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
