package list

import (
	"fmt"
	"os"
	"strconv"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/projects"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/models"
	"github.com/spf13/cobra"
)

var assignee string
var project uint64
var version string

var workPackagesCmd = &cobra.Command{
	Use:     "workpackages",
	Aliases: []string{"wps"},
	Short:   "Lists work packages",
	Long:    "Get a list of visible work packages. Filter flags can be applied.",
	Run:     listWorkPackages,
}

func listWorkPackages(_ *cobra.Command, _ []string) {
	if len(version) != 0 && project == 0 {
		printer.ErrorText("Version flag (--version) can only be used in conjunction with project flag (-p or --project).")
	}

	if all, err := work_packages.All(filterOptions()); err == nil {
		printer.WorkPackages(all)
	} else {
		printer.Error(err)
	}
}

func filterOptions() *map[work_packages.FilterOption]string {
	options := make(map[work_packages.FilterOption]string)

	if project > 0 {
		options[work_packages.Project] = strconv.FormatUint(project, 10)
	}

	if len(assignee) > 0 {
		options[work_packages.Assignee] = assignee
	}

	if len(version) > 0 {
		options[work_packages.Version] = validatedVersionId(version)
	}

	return &options
}

func validatedVersionId(version string) string {
	project := projects.Lookup(project)
	if project == nil {
		printer.ErrorText(fmt.Sprintf("Project with ID (%d) not found.", project.Id))
	}

	versions := projects.AvailableVersions(project.Id)

	filteredVersions := common.Filter(versions, func(v *models.Version) bool {
		return v.Name == version
	})

	if len(filteredVersions) != 1 {
		printer.Info(fmt.Sprintf(
			"No unique available version from input '%s' found for project [#%d]. Please use one of the versions listed below.",
			version,
			project.Id,
		))

		printer.Versions(versions)

		os.Exit(-1)
	}

	return strconv.FormatUint(filteredVersions[0].Id, 10)
}
