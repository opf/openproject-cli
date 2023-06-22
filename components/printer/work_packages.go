package printer

import (
	"fmt"

	"github.com/opf/openproject-cli/models"
)

func WorkPackage(workPackage *models.WorkPackage) {
	id := fmt.Sprintf("#%d", workPackage.Id)
	fmt.Printf("[%s] %s\n", red(id), cyan(workPackage.Subject))
}
