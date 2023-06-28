package inspect

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/launch"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/projects"
	"github.com/opf/openproject-cli/components/routes"
)

var shouldOpenProjectInBrowser bool

var inspectProjectCmd = &cobra.Command{
	Use:   "project [id]",
	Short: "Show details about a project",
	Long:  "Show detailed information of a project refereced by it's ID.",
	Run:   inspectProject,
}

func inspectProject(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [id], but got %d", len(args)))
		return
	}

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		printer.ErrorText(fmt.Sprintf("'%s' is an invalid project id. Must be a number.", args[0]))
		return
	}

	project, err := projects.Lookup(id)
	if err != nil {
		printer.Error(err)
		return
	}

	if shouldOpenProjectInBrowser {
		err = launch.Browser(routes.ProjectUrl(project))
		if err != nil {
			printer.ErrorText(fmt.Sprintf("Error opening browser: %+v", err))
		}
	} else {
		printer.Project(project)
	}
}
