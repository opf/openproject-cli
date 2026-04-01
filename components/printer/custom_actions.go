package printer

import "github.com/opf/openproject-cli/models"

func CustomActions(actions []*models.CustomAction) {
	activeRenderer.CustomActions(actions)
}
