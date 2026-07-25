package budget

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/budgets"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [id]",
	Short: "Show details about a budget",
	Long:  "Show detailed information of a budget referenced by its ID.",
	RunE:  inspectBudget,
}

func inspectBudget(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [id], but got %d", len(args)))
		return openerrors.ErrHandled
	}

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		printer.ErrorText(fmt.Sprintf("'%s' is an invalid budget id. Must be a number.", args[0]))
		return openerrors.ErrHandled
	}

	budget, err := budgets.Lookup(id)
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	printer.Budget(budget)
	return nil
}
