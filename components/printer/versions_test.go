package printer_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func TestVersions(t *testing.T) {
	testingPrinter.Reset()

	versions := []*models.Version{
		{
			Id:   2,
			Name: "13.0",
		},
		{
			Id:   4,
			Name: "13.1",
		},
		{
			Id:   43,
			Name: "45.5",
		},
	}

	expected := common.Reduce(
		versions,
		func(state string, version *models.Version) string {
			idString := "#" + strconv.FormatUint(version.Id, 10)
			return state + fmt.Sprintf("[%s] %s\n", printer.Red(idString), printer.Cyan(version.Name))
		},
		"")

	printer.Versions(versions)

	if testingPrinter.Result != expected {
		t.Errorf("\nExpected:\n%sbut got:\n%s", expected, testingPrinter.Result)
	}
}
