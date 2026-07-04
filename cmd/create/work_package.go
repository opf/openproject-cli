package create

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/launch"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/components/routes"
)

var projectId uint64
var shouldOpenWorkPackageInBrowser bool
var typeFlag string
var descriptionFlag string
var descriptionFileFlag string
var parentFlag uint64

var createWorkPackageCmd = &cobra.Command{
	Use:   "workpackage [subject]",
	Short: "Create work package in project",
	Long:  "Create a new work package with the given subject in a project",
	Run:   createWorkPackage,
}

func createWorkPackage(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [subject], but got %d", len(args)))
		return
	}

	subject := args[0]
	workPackage, err := work_packages.Create(projectId, createOptions(subject))
	if err != nil {
		printer.Error(err)
		return
	}

	if shouldOpenWorkPackageInBrowser {
		err = launch.Browser(routes.WorkPackageUrl(workPackage))
		if err != nil {
			printer.ErrorText(fmt.Sprintf("Error opening browser: %+v", err))
		}
	} else {
		printer.WorkPackage(workPackage)
	}
}

func createOptions(subject string) map[work_packages.CreateOption]string {
	var options = make(map[work_packages.CreateOption]string)

	options[work_packages.CreateSubject] = subject

	if len(typeFlag) > 0 {
		options[work_packages.CreateType] = typeFlag
	}
	if len(descriptionFlag) > 0 {
		options[work_packages.CreateDescription] = descriptionFlag
	}
	if len(descriptionFileFlag) > 0 {
		options[work_packages.CreateDescriptionFile] = descriptionFileFlag
	}
	if parentFlag > 0 {
		options[work_packages.CreateParent] = strconv.FormatUint(parentFlag, 10)
	}

	return options
}
