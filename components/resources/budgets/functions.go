package budgets

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func Lookup(id uint64) (*models.Budget, error) {
	response, err := requests.Get(paths.Budget(id), nil)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.BudgetDto](response)
	return element.Convert(), nil
}

func AllForProject(projectId uint64) ([]*models.Budget, error) {
	query := requests.NewPaginatedQuery(-1, nil)
	response, err := requests.Get(paths.ProjectBudgets(projectId), &query)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.BudgetCollectionDto](response)
	return element.Convert(), nil
}
