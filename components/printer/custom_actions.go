package printer

import (
	"fmt"

	"github.com/opf/openproject-cli/models"
)

func CustomActions(actions []*models.CustomAction) {
	for _, a := range actions {
		printCustomAction(a)
	}
}

func printCustomAction(action *models.CustomAction) {
	id := fmt.Sprintf("#%d", action.Id)
	fmt.Printf("[%s] %s\n", red(id), cyan(action.Name))
}
