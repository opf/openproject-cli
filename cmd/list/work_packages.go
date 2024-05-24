package list

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/projects"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/models"
	"github.com/spf13/cobra"
)

var assignee string
var projectId uint64
var version string
var showTotal bool
var statusFilter string
var typeFilter string

var workPackagesCmd = &cobra.Command{
	Use:     "workpackages",
	Aliases: []string{"wps"},
	Short:   "Lists work packages",
	Long:    "Get a list of visible work packages. Filter flags can be applied.",
	Run:     listWorkPackages,
}

func listWorkPackages(_ *cobra.Command, _ []string) {
	if len(version) != 0 && projectId == 0 {
		printer.ErrorText("Version flag (--version) can only be used in conjunction with projectId flag (-p or --project-id).")
		return
	}

	collection, err := work_packages.All(filterOptions())
	switch {
	case err == nil && showTotal:
		printer.Number(collection.Total)
	case err == nil:
		printer.WorkPackages(collection.Items)
	default:
		printer.Error(err)
	}
}

func filterOptions() *map[work_packages.FilterOption]string {
	options := make(map[work_packages.FilterOption]string)

	if projectId > 0 {
		options[work_packages.Project] = strconv.FormatUint(projectId, 10)
	}

	if len(assignee) > 0 {
		options[work_packages.Assignee] = assignee
	}

	if len(statusFilter) > 0 {
		options[work_packages.Status] = validateStatusFilterValue(statusFilter)
	}

	if len(typeFilter) > 0 {
		options[work_packages.Type] = validateStatusFilterValue(typeFilter)
	}

	if len(version) > 0 {
		options[work_packages.Version] = validatedVersionId(version)
	}

	return &options
}

func validatedVersionId(version string) string {
	project, err := projects.Lookup(projectId)
	if err != nil {
		printer.Error(err)
	}

	versions, err := projects.AvailableVersions(project.Id)
	if err != nil {
		printer.Error(err)
	}

	filteredVersions := common.Filter(versions, func(v *models.Version) bool {
		return v.Name == version
	})

	if len(filteredVersions) != 1 {
		printer.Info(fmt.Sprintf(
			"No unique available version from input %s found for projectId %s. Please use one of the versions listed below.",
			printer.Cyan(version),
			printer.Red(fmt.Sprintf("#%d", project.Id)),
		))

		printer.Versions(versions)

		os.Exit(-1)
	}

	return strconv.FormatUint(filteredVersions[0].Id, 10)
}

func validateStatusFilterValue(status string) string {
	matched, err := regexp.Match(`^(open)$|^(closed)$|^(!?[0-9,]+)$`, []byte(status))
	if err != nil {
		printer.Error(err)
	}

	if !matched {
		printer.ErrorText(fmt.Sprintf("Invalid status filter value %s.", printer.Yellow(status)))
		os.Exit(-1)
	}

	return status
}
