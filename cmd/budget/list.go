package budget

import (
	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/budgets"
)

var listProjectId uint64

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists budgets of a project",
	Long:  "Get a list of all budgets for a given project.",
	Run:   listBudgets,
}

func listBudgets(_ *cobra.Command, _ []string) {
	all, err := budgets.AllForProject(listProjectId)
	if err != nil {
		printer.Error(err)
		return
	}

	printer.Budgets(all)
}
