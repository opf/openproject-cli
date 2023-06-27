package printer

import (
	"fmt"
	"strings"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/models"
)

func Types(types []*models.Type) {
	var maxIdLength = 0
	for _, t := range types {
		maxIdLength = common.Max(maxIdLength, idLength(t.Id))
	}

	for _, t := range types {
		printType(t, maxIdLength)
	}
}

func printType(workPackageType *models.Type, maxIdLength int) {
	var parts []string

	diff := maxIdLength - idLength(workPackageType.Id)
	idStr := fmt.Sprintf("%s#%d", indent(diff), workPackageType.Id)
	parts = append(parts, Red(idStr))

	parts = append(parts, Yellow(workPackageType.Name))
	activePrinter.Println(strings.Join(parts, " "))
}
