package workpackage

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var updateActionFlag string
var updateAssigneeFlag uint64
var updateAttachFlag string
var updateSubjectFlag string
var updateTypeFlag string

var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Updates the work package",
	Long: `Update a work package. Each update
provided by a flag is executed on its own.`,
	Run: updateWorkPackage,
}

func updateWorkPackage(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [id], but got %d", len(args)))
		return
	}

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		printer.ErrorText(fmt.Sprintf("'%s' is an invalid work package id. Must be a number.", args[0]))
		return
	}

	if workPackage, err := work_packages.Update(id, updateOptions()); err == nil {
		printer.Info("-- ")
		printer.WorkPackage(workPackage)
	} else {
		printer.Error(err)
	}
}

func updateOptions() map[work_packages.UpdateOption]string {
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
	if len(updateSubjectFlag) > 0 {
		options[work_packages.UpdateSubject] = updateSubjectFlag
	}
	if len(updateTypeFlag) > 0 {
		options[work_packages.UpdateType] = updateTypeFlag
	}
	return options
}
