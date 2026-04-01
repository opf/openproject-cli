package printer

import "github.com/opf/openproject-cli/models"

func Types(types []*models.Type) {
	activeRenderer.Types(types)
}
