package project

import (
	"fmt"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
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
	RunE:  inspectProject,
}

func inspectProject(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [id|identifier], but got %d", len(args)))
		return openerrors.ErrHandled
	}

	if err := projects.ValidateIdentifier(args[0]); err != nil {
		printer.ErrorText(err.Error())
		return openerrors.ErrHandled
	}

	project, err := projects.Lookup(args[0])
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	if openInBrowser {
		err = launch.Browser(routes.ProjectUrl(project))
		if err != nil {
			printer.ErrorText(fmt.Sprintf("Error opening browser: %+v", err))
			return openerrors.ErrHandled
		}
	} else {
		printer.Project(project)
	}
	return nil
}
