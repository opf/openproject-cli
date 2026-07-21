package workpackage

import (
	"fmt"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/launch"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/components/routes"
)

var inspectOpenInBrowser bool
var inspectListAvailableTypes bool

var inspectCmd = &cobra.Command{
	Use:   "inspect [id]",
	Short: "Show details about a work package",
	Long:  "Show detailed information of a work package referenced by its numeric ID (e.g. 12345) or project-based identifier (e.g. PROJ-123).",
	RunE:  inspectWorkPackage,
}

func inspectWorkPackage(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [id], but got %d", len(args)))
		return openerrors.ErrHandled
	}

	id := args[0]
	if err := work_packages.ValidateIdentifier(id); err != nil {
		printer.ErrorText(err.Error())
		return openerrors.ErrHandled
	}

	if inspectHasListingFlag() {
		switch {
		case inspectListAvailableTypes:
			return inspectAvailableTypes(id)
		}
		return nil
	}

	workPackage, err := work_packages.Lookup(id)
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	if inspectOpenInBrowser {
		err = launch.Browser(routes.WorkPackageUrl(workPackage))
		if err != nil {
			printer.ErrorText(fmt.Sprintf("Error opening browser: %+v", err))
			return openerrors.ErrHandled
		}
	} else {
		printer.WorkPackage(workPackage)
	}
	return nil
}

func inspectAvailableTypes(id string) error {
	types, err := work_packages.AvailableTypes(id)
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}
	printer.Types(types)
	return nil
}

func inspectHasListingFlag() bool {
	return inspectListAvailableTypes
}
