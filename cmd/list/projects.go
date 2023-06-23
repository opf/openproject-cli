package list

import (
	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/projects"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Lists projects",
	Long: `Get a list of visible projects.
The list can get filtered further.`,
	Run: listProjects,
}

func listProjects(cmd *cobra.Command, args []string) {
	printer.Projects(projects.All())
}
