package update

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var (
	actionFlag string
	attachFlag string
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
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		printer.ErrorText(fmt.Sprintf("'%s' is an invalid work package id. Must be a number.", args[0]))
	}

	var opts = make(map[work_packages.UpdateOption]string)
	if len(actionFlag) > 0 {
		opts[work_packages.Action] = actionFlag
	}
	if len(attachFlag) > 0 {
		opts[work_packages.Attach] = attachFlag
	}

	work_packages.Update(id, opts)
}
