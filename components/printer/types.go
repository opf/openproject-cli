package printer

import (
	"fmt"

	"github.com/opf/openproject-cli/models"
)

func Types(types []*models.Type) {
	activeRenderer.Types(types)
}

// AvailableTypes lists type names on standard error as a guidance hint, so it
// never pollutes machine-readable output (e.g. JSON) on a failed command.
func AvailableTypes(types []*models.Type) {
	Info("Available types:")
	for _, t := range types {
		Info(fmt.Sprintf("  - %s", t.Name))
	}
}
