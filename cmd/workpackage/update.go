package workpackage

import (
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

	workPackage, err := work_packages.Update(id, updateOptions(cmd))
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	printer.Info("-- ")
	printer.WorkPackage(workPackage)
	return nil
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
	return options
}
