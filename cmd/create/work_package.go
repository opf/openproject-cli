package create

import (
	"fmt"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
)

var createWorkPackageCmd = &cobra.Command{
	Use:   "workpackage [project_id] [subject]",
	Short: "Create work package in project",
	Long:  "Create a new work package with the given subject in a project selected by a project's ID",
	Run:   createWorkPackage,
}

func createWorkPackage(_ *cobra.Command, args []string) {
	if len(args) != 2 {
		printer.ErrorText(fmt.Sprintf("Expected 2 arguments [project_id] and [subject], but got %d", len(args)))
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		printer.ErrorText(fmt.Sprintf("'%s' is an invalid project id. Must be a number.", args[0]))
	}

	subject := args[1]

	workPackage := work_packages.Create(id, subject)

	printer.WorkPackage(workPackage)
}
