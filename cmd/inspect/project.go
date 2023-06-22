package inspect

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/projects"
)

var inspectProjectCmd = &cobra.Command{
	Use:   "project [id]",
	Short: "Show details about a project",
	Long:  "Show detailed information of a project refereced by it's ID.",
	Run:   inspectProject,
}

func inspectProject(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [id], but got %d", len(args)))
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		printer.ErrorText(fmt.Sprintf("'%s' is an invalid project id. Must be a number.", args[0]))
	}

	printer.Project(projects.Lookup(id))
}
