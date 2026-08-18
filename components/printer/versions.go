package printer

import "github.com/opf/openproject-cli/models"

func Versions(versions []*models.Version) {
	activeRenderer.Versions(versions)
}
