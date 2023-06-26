package work_packages

import (
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/models"
)

func AssigneeFilter(principal *models.Principal) requests.Filter {
	return requests.Filter{
		Operator: "=",
		Name:     "assignee",
		Values:   []string{principal.Name},
	}
}
