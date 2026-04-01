package printer

import "github.com/opf/openproject-cli/models"

func Projects(projects []*models.Project) {
	activeRenderer.Projects(projects)
}

func Project(project *models.Project) {
	activeRenderer.Project(project)
}
