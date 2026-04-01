package printer_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func TestBudget(t *testing.T) {
	testingPrinter.Reset()

	budget := models.Budget{Id: 7, Subject: "Conception"}

	idString := "#" + strconv.FormatUint(budget.Id, 10)
	expected := fmt.Sprintf("%s %s\n", printer.Red(idString), printer.Cyan(budget.Subject))

	printer.Budget(&budget)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}

func TestBudgets_AlignsIds(t *testing.T) {
	testingPrinter.Reset()

	budgets := []*models.Budget{
		{Id: 7, Subject: "Conception"},
		{Id: 42, Subject: "Développement"},
		{Id: 100, Subject: "Recette & déploiement"},
	}

	// IDs are right-aligned to the width of the longest ID (3 digits)
	expected := "" +
		fmt.Sprintf("%s %s\n", printer.Red("  #7"), printer.Cyan("Conception")) +
		fmt.Sprintf("%s %s\n", printer.Red(" #42"), printer.Cyan("Développement")) +
		fmt.Sprintf("%s %s\n", printer.Red("#100"), printer.Cyan("Recette & déploiement"))

	printer.Budgets(budgets)

	if testingPrinter.Result != expected {
		t.Errorf("Expected %s, but got %s", expected, testingPrinter.Result)
	}
}
