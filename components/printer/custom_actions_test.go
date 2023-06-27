package printer_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func TestCustomActions(t *testing.T) {
	testingPrinter.Reset()

	actions := []*models.CustomAction{
		{
			Id:   2,
			Name: "Claim",
		},
		{
			Id:   4,
			Name: "Developed",
		},
		{
			Id:   43,
			Name: "Tested",
		},
	}

	expected := common.Reduce[*models.CustomAction, string](
		actions,
		func(state string, action *models.CustomAction) string {
			idString := "#" + strconv.FormatUint(action.Id, 10)
			return state + fmt.Sprintf("[%s] %s\n", printer.Red(idString), printer.Cyan(action.Name))
		},
		"")

	printer.CustomActions(actions)

	if testingPrinter.Result != expected {
		t.Errorf("\nExpected:\n%sbut got:\n%s", expected, testingPrinter.Result)
	}
}
