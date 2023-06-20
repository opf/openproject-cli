package printer

import (
	"fmt"
	"github.com/opf/openproject-cli/models"
)

func Projects(v interface{}) {
	list, ok := v.([]*models.Project)
	if ok {
		for _, p := range list {
			printProject(p)
		}
	}

	single, ok := v.(*models.Project)
	if ok {
		printProject(single)
	}
}

func printProject(p *models.Project) {
	id := fmt.Sprintf("#%d", p.Id)
	fmt.Printf("[%s] %s\n", red(id), cyan(p.Name))
}
