package list

import (
	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/models"
)

var assignee string

var workPackagesCmd = &cobra.Command{
	Use:   "workpackages",
	Short: "Lists work packages",
	Long:  "Get a list of visible work packages. Filter flags can be applied.",
	Run:   listWorkPackages,
}

func listWorkPackages(_ *cobra.Command, _ []string) {
	var principal *models.Principal

	if assignee != "" {
		principal = &models.Principal{Name: assignee}
	}

	printer.WorkPackages(work_packages.All(principal))
}
