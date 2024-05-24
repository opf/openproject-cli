package printer

import (
	"fmt"
	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/models"
	"strings"
)

func StatusList(status []*models.Status) {
	var maxIdLength = 0
	for _, s := range status {
		maxIdLength = common.Max(maxIdLength, idLength(s.Id))
	}

	for _, s := range status {
		printStatus(s, maxIdLength)
	}
}

func Status(status *models.Status) {
	printStatus(status, idLength(status.Id))
}

func printStatus(status *models.Status, maxIdLength int) {
	var parts []string

	diff := maxIdLength - idLength(status.Id)
	idStr := fmt.Sprintf("%s#%d", indent(diff), status.Id)
	parts = append(parts, Red(idStr))

	parts = append(parts, Cyan(status.Name))

	if status.IsDefault {
		parts = append(parts, fmt.Sprintf("(%s)", Yellow("default")))
	}

	activePrinter.Println(strings.Join(parts, " "))
}
