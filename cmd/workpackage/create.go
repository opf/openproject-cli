package workpackage

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/launch"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/projects"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/components/routes"
)

var createProjectId string
var createOpenInBrowser bool
var createTypeFlag string
var createAssigneeFlag uint64
var createDescriptionFlag string

var createCmd = &cobra.Command{
	Use:   "create [subject]",
	Short: "Create work package in project",
	Long:  "Create a new work package with the given subject in a project",
	RunE:  createWorkPackage,
}

func createWorkPackage(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [subject], but got %d", len(args)))
		return openerrors.ErrHandled
	}

	subject := args[0]
	if err := projects.ValidateIdentifier(createProjectId); err != nil {
		printer.ErrorText(fmt.Sprintf("--project: %s", err.Error()))
		return openerrors.ErrHandled
	}

	workPackage, err := work_packages.Create(createProjectId, createOptions(cmd, subject))
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	if createOpenInBrowser {
		err = launch.Browser(routes.WorkPackageUrl(workPackage))
		if err != nil {
			printer.ErrorText(fmt.Sprintf("Error opening browser: %+v", err))
		}
	} else {
		printer.WorkPackage(workPackage)
	}
	return nil
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
