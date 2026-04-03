package workpackage

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/launch"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/components/routes"
)

var createProjectId uint64
var createOpenInBrowser bool
var createTypeFlag string
var createAssigneeFlag uint64
var createDescriptionFlag string

var createCmd = &cobra.Command{
	Use:   "create [subject]",
	Short: "Create work package in project",
	Long:  "Create a new work package with the given subject in a project",
	Run:   createWorkPackage,
}

func createWorkPackage(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [subject], but got %d", len(args)))
		return
	}

	subject := args[0]
	workPackage, err := work_packages.Create(createProjectId, createOptions(cmd, subject))
	if err != nil {
		printer.Error(err)
		return
	}

	if createOpenInBrowser {
		err = launch.Browser(routes.WorkPackageUrl(workPackage))
		if err != nil {
			printer.ErrorText(fmt.Sprintf("Error opening browser: %+v", err))
		}
	} else {
		printer.WorkPackage(workPackage)
	}
}

func createOptions(cmd *cobra.Command, subject string) map[work_packages.CreateOption]string {
	options := make(map[work_packages.CreateOption]string)
	options[work_packages.CreateSubject] = subject
	if len(createTypeFlag) > 0 {
		options[work_packages.CreateType] = createTypeFlag
	}
	if createAssigneeFlag > 0 {
		options[work_packages.CreateAssignee] = strconv.FormatUint(createAssigneeFlag, 10)
	}
	if cmd.Flags().Changed("description") {
		options[work_packages.CreateDescription] = createDescriptionFlag
	}
	return options
}
