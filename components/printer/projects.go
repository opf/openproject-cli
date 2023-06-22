package printer

import (
	"fmt"

	"github.com/opf/openproject-cli/models"
)

func Projects(projects []*models.Project) {
	for _, p := range projects {
		printProject(p)
	}
}

func Project(project *models.Project) {
	printProject(project)
}

func printProject(p *models.Project) {
	id := fmt.Sprintf("#%d", p.Id)
	fmt.Printf("[%s] %s\n", red(id), cyan(p.Name))
}
