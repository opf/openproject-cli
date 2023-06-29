package update

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var (
	actionFlag   string
	assigneeFlag uint64
	attachFlag   string
	subjectFlag  string
	typeFlag     string
)

var workPackageCmd = &cobra.Command{
	Use:   "workpackage [id]",
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
	var options = make(map[work_packages.UpdateOption]string)
	if len(actionFlag) > 0 {
		options[work_packages.UpdateCustomAction] = actionFlag
	}
	if assigneeFlag > 0 {
		options[work_packages.UpdateAssignee] = strconv.FormatUint(assigneeFlag, 10)
	}
	if len(attachFlag) > 0 {
		options[work_packages.UpdateAttachment] = attachFlag
	}
	if len(subjectFlag) > 0 {
		options[work_packages.UpdateSubject] = subjectFlag
	}
	if len(typeFlag) > 0 {
		options[work_packages.UpdateType] = typeFlag
	}

	return options
}
