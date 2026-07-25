package workpackage

import (
	stderrors "errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var updateActionFlag string
var updateAssigneeFlag uint64
var updateAttachFlag string
var updateDescriptionFlag string
var updateSubjectFlag string
var updateTypeFlag string
var updateStatusFlag string
var updateSetFlags []string
var updateDryRun bool

var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Updates the work package",
	Long: `Update a work package referenced by its numeric ID (e.g. 12345) or project-based identifier (e.g. PROJ-123). Each update
provided by a flag is executed on its own.`,
	RunE: updateWorkPackage,
}

func updateWorkPackage(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [id], but got %d", len(args)))
		return openerrors.ErrHandled
	}

	id := args[0]
	if err := work_packages.ValidateIdentifier(id); err != nil {
		printer.ErrorText(err.Error())
		return openerrors.ErrHandled
	}

	if len(updateSetFlags) > 0 && updateHasFieldFlags(cmd) {
		printer.ErrorText("cannot combine --set with other update flags")
		return openerrors.ErrHandled
	}

	if len(updateSetFlags) > 0 {
		if updateDryRun {
			plan, err := work_packages.DryRunUpdateFields(id, updateSetFlags)
			if err != nil {
				printer.Error(err)
				return openerrors.ErrHandled
			}
			printer.WorkPackageUpdatePlan(plan)
			return nil
		}
		if err := work_packages.UpdateFields(id, updateSetFlags); err != nil {
			printer.Error(err)
			return openerrors.ErrHandled
		}
		payload, err := work_packages.Inspect(id)
		if err != nil {
			printer.Error(err)
			return openerrors.ErrHandled
		}
		printer.WorkPackageDetails(payload)
		return nil
	}

	options := updateOptions(cmd)
	if len(options) == 0 {
		printer.ErrorText("No update options provided. Use --help to see available flags.")
		return openerrors.ErrHandled
	}

	if updateDryRun {
		plan, err := work_packages.DryRunUpdate(id, options)
		if err != nil {
			printer.Error(err)
			return openerrors.ErrHandled
		}
		printer.WorkPackageUpdatePlan(plan)
		return nil
	}

	workPackage, err := work_packages.Update(id, options)
	if err != nil {
		if !stderrors.Is(err, openerrors.ErrHandled) {
			printer.Error(err)
		}
		return openerrors.ErrHandled
	}

	printer.Info("-- ")
	printer.WorkPackage(workPackage)
	return nil
}

func updateHasFieldFlags(cmd *cobra.Command) bool {
	for _, name := range []string{"subject", "type", "assignee", "description", "status", "action", "attach"} {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	return false
}

func updateOptions(cmd *cobra.Command) map[work_packages.UpdateOption]string {
	options := make(map[work_packages.UpdateOption]string)
	if len(updateActionFlag) > 0 {
		options[work_packages.UpdateCustomAction] = updateActionFlag
	}
	if updateAssigneeFlag > 0 {
		options[work_packages.UpdateAssignee] = strconv.FormatUint(updateAssigneeFlag, 10)
	}
	if len(updateAttachFlag) > 0 {
		options[work_packages.UpdateAttachment] = updateAttachFlag
	}
	if cmd.Flags().Changed("description") {
		options[work_packages.UpdateDescription] = updateDescriptionFlag
	}
	if len(updateSubjectFlag) > 0 {
		options[work_packages.UpdateSubject] = updateSubjectFlag
	}
	if len(updateTypeFlag) > 0 {
		options[work_packages.UpdateType] = updateTypeFlag
	}
	if len(updateStatusFlag) > 0 {
		options[work_packages.UpdateStatus] = updateStatusFlag
	}
	return options
}
