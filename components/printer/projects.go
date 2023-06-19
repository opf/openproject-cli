package printer

import (
	"fmt"
	
	"github.com/opf/openproject-cli/components/resources/projects"
)

func Projects(v interface{})  {
	list, ok := v.([]*projects.Project)
	if ok {
		for _, p := range list {
			printProject(p)
		}
	}
	
	single, ok := v.(*projects.Project)
	if ok {
		printProject(single)
	}
}

func printProject(p *projects.Project)  {
	id := fmt.Sprintf("#%d", p.Id)
	fmt.Printf("[%s] %s\n", red(id), cyan(p.Name))
}