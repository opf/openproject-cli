package version

import (
	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/projects"
)

var listProjectId string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists versions",
	Long:  "Get a list of a project's versions, scoped by the provided flag (--project).",
	RunE:  listVersions,
}

func listVersions(_ *cobra.Command, _ []string) error {
	if err := projects.ValidateIdentifier(listProjectId); err != nil {
		printer.ErrorText(err.Error())
		return openerrors.ErrHandled
	}

	versions, err := projects.AvailableVersions(listProjectId)
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	printer.Versions(versions)
	return nil
}
