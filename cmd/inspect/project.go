package inspect

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/projects"
)

var inspectProjectCmd = &cobra.Command{
	Use:   "project [id]",
	Short: "Show details about a project",
	Long:  "Show detailed information of a project refereced by it's ID.",
	Run:   execute,
}

func execute(cmd *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])
	all := projects.Find(id)
	printer.Projects(all)
}
