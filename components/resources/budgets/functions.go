package budgets

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func AllForProject(projectId uint64) ([]*models.Budget, error) {
	query := requests.NewPaginatedQuery(-1, nil)
	response, err := requests.Get(paths.ProjectBudgets(projectId), &query)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.BudgetCollectionDto](response)
	return element.Convert(), nil
}
