package project

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/launch"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/projects"
	"github.com/opf/openproject-cli/components/routes"
)

var openInBrowser bool

var inspectCmd = &cobra.Command{
	Use:   "inspect [id|identifier]",
	Short: "Show details about a project",
	Long:  "Show detailed information of a project referenced by its numeric ID or identifier.",
	Run:   inspectProject,
}

func inspectProject(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [id|identifier], but got %d", len(args)))
		return
	}

	if err := projects.ValidateIdentifier(args[0]); err != nil {
		printer.ErrorText(err.Error())
		return
	}

	project, err := projects.Lookup(args[0])
	if err != nil {
		printer.Error(err)
		return
	}

	if openInBrowser {
		err = launch.Browser(routes.ProjectUrl(project))
		if err != nil {
			printer.ErrorText(fmt.Sprintf("Error opening browser: %+v", err))
		}
	} else {
		printer.Project(project)
	}
}
