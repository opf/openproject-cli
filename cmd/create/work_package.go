package create

import (
	"fmt"

	"github.com/spf13/cobra"

	componentErrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/launch"
	"github.com/opf/openproject-cli/components/presenter"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/components/routes"
)

var projectId uint64
var parentWorkPackageID uint64
var shouldOpenWorkPackageInBrowser bool
var printCreatedWorkPackageAsJSON bool
var dryRunCreateWorkPackage bool
var typeFlag string
var descriptionFlag string
var descriptionFlagChanged bool

var createWorkPackageCmd = &cobra.Command{
	Use:   "workpackage [subject]",
	Short: "Create work package in project",
	Long:  "Create a new work package with the given subject in a project",
	Run:   createWorkPackage,
}

func createWorkPackage(cmd *cobra.Command, args []string) {
	if cmd != nil {
		descriptionFlagChanged = cmd.Flags().Changed("description")
	}

	if len(args) != 1 {
		printCreateError("invalid_argument", fmt.Sprintf("Expected 1 argument [subject], but got %d", len(args)))
		return
	}

	if err := validateCreateWorkPackageFlags(); err != nil {
		printCreateError("conflicting_arguments", err.Error())
		return
	}

	subject := args[0]
	options := createOptions(subject)

	if dryRunCreateWorkPackage {
		plan, err := work_packages.DryRunCreate(projectId, options)
		if err != nil {
			printCreateError(createErrorCode(err), err.Error())
			return
		}

		data, err := presenter.MarshalJSON(plan)
		if err != nil {
			printer.Error(err)
			return
		}

		printer.Info(string(data))
		return
	}

	workPackage, err := work_packages.Create(projectId, options)
	if err != nil {
		printCreateError(createErrorCode(err), err.Error())
		return
	}

	if printCreatedWorkPackageAsJSON {
		payload, err := work_packages.Inspect(workPackage.Id)
		if err != nil {
			printCreateError("post_apply_inspect_failed", err.Error())
			return
		}

		data, err := presenter.MarshalJSON(payload)
		if err != nil {
			printer.Error(err)
			return
		}

		printer.Info(string(data))
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

	if parentWorkPackageID > 0 {
		options[work_packages.CreateParent] = fmt.Sprintf("%d", parentWorkPackageID)
	}

	if descriptionFlagChanged {
		options[work_packages.CreateDescription] = descriptionFlag
	}

	return options
}

func validateCreateWorkPackageFlags() error {
	if shouldOpenWorkPackageInBrowser && printCreatedWorkPackageAsJSON {
		return fmt.Errorf("cannot use --open together with --json")
	}

	if dryRunCreateWorkPackage && !printCreatedWorkPackageAsJSON {
		return fmt.Errorf("cannot use --dry-run without --json")
	}

	return nil
}

func printCreateError(code, message string) {
	if !printCreatedWorkPackageAsJSON {
		printer.ErrorText(message)
		return
	}

	data, err := presenter.MarshalError(code, message)
	if err != nil {
		printer.Error(err)
		return
	}

	printer.Info(string(data))
}

func createErrorCode(err error) string {
	if componentErrors.IsCustom(err) {
		return "invalid_argument"
	}

	return "api_error"
}
