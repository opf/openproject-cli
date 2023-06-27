package printer

import (
	"fmt"

	"github.com/opf/openproject-cli/models"
)

func Types(types []*models.Type) {
	for _, t := range types {
		printType(t)
	}
}

func printType(workPackageType *models.Type) {
	id := fmt.Sprintf("#%d", workPackageType.Id)
	activePrinter.Printf("%s %s\n", Red(id), Yellow(workPackageType.Name))
}
