package list

import (
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/spf13/cobra"
)

var workPackagesCmd = &cobra.Command{
	Use:   "workpackages",
	Short: "Lists work packages",
	Long:  "Get a list of visible work packages.",
	Run:   listWorkPackages,
}

func listWorkPackages(cmd *cobra.Command, args []string) {
	printer.WorkPackages(work_packages.All())
}
