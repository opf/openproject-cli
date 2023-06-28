package printer

import (
	"fmt"

	"github.com/opf/openproject-cli/models"
)

func Versions(versions []*models.Version) {
	for _, a := range versions {
		printVersion(a)
	}
}

func printVersion(version *models.Version) {
	id := fmt.Sprintf("#%d", version.Id)
	activePrinter.Printf("[%s] %s\n", Red(id), Cyan(version.Name))
}
