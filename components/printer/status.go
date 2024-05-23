package printer

import (
	"fmt"
	"github.com/opf/openproject-cli/models"
)

func StatusList(statuss []*models.Status) {
	for _, p := range statuss {
		printStatus(p)
	}
}

func Status(status *models.Status) {
	printStatus(status)
}

func printStatus(status *models.Status) {
	var defaultSuffix string

	id := fmt.Sprintf("#%d", status.Id)
	if status.IsDefault {
		defaultSuffix = fmt.Sprintf(" (%s)", Yellow("default"))
	}
	activePrinter.Printf("%s %s%s\n", Red(id), Cyan(status.Name), defaultSuffix)
}
