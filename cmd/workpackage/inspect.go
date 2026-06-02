package workpackage

import (
	"fmt"

	"github.com/spf13/cobra"

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
	Run:   inspectWorkPackage,
}

func inspectWorkPackage(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [id], but got %d", len(args)))
		return
	}

	id := args[0]
	if err := work_packages.ValidateIdentifier(id); err != nil {
		printer.ErrorText(err.Error())
		return
	}

	if inspectHasListingFlag() {
		switch {
		case inspectListAvailableTypes:
			inspectAvailableTypes(id)
		}
		return
	}

	workPackage, err := work_packages.Lookup(id)
	if err != nil {
		printer.Error(err)
		return
	}

	if inspectOpenInBrowser {
		err = launch.Browser(routes.WorkPackageUrl(workPackage))
		if err != nil {
			printer.ErrorText(fmt.Sprintf("Error opening browser: %+v", err))
		}
	} else {
		printer.WorkPackage(workPackage)
	}
}

func inspectAvailableTypes(id string) {
	types, err := work_packages.AvailableTypes(id)
	if err != nil {
		printer.Error(err)
		return
	}
	printer.Types(types)
}

func inspectHasListingFlag() bool {
	return inspectListAvailableTypes
}
