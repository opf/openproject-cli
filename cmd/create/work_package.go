package create

import (
	"fmt"
	"github.com/opf/openproject-cli/components/launch"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/components/routes"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
)

var projectId int64
var shouldOpenWorkPackageInBrowser bool

var createWorkPackageCmd = &cobra.Command{
	Use:   "workpackage [subject]",
	Short: "Create work package in project",
	Long:  "Create a new work package with the given subject in a project",
	Run:   createWorkPackage,
}

func createWorkPackage(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [subject], but got %d", len(args)))
	}

	subject := args[0]

	workPackage := work_packages.Create(projectId, subject)

	if shouldOpenWorkPackageInBrowser {
		err := launch.Browser(routes.WorkPackageUrl(workPackage))
		if err != nil {
			printer.ErrorText(fmt.Sprintf("Error opening browser: %s", err))
		}
	} else {
		printer.WorkPackage(workPackage)
	}
}
