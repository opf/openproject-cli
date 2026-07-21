package budget

import (
	"fmt"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/budgets"
	"github.com/opf/openproject-cli/components/resources/projects"
)

var listProjectId string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists budgets of a project",
	Long:  "Get a list of all budgets for a given project.",
	RunE:  listBudgets,
}

func listBudgets(_ *cobra.Command, _ []string) error {
	if err := projects.ValidateIdentifier(listProjectId); err != nil {
		printer.ErrorText(fmt.Sprintf("--project: %s", err.Error()))
		return openerrors.ErrHandled
	}

	all, err := budgets.AllForProject(listProjectId)
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	printer.Budgets(all)
	return nil
}
