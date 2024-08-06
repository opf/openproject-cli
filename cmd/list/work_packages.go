package list

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources"
	"github.com/opf/openproject-cli/components/resources/projects"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/components/resources/work_packages/filters"
	"github.com/opf/openproject-cli/models"
	"github.com/spf13/cobra"
)

var assignee string
var projectId uint64
var showTotal bool
var statusFilter string
var typeFilter string
var includeSubProjects bool

var activeFilters = map[string]resources.Filter{
	"notSubProject": filters.NewNotSubProjectFilter(),
	"notVersion":    filters.NewNotVersionFilter(),
	"subProject":    filters.NewSubProjectFilter(),
	"timestamp":     filters.NewTimestampFilter(),
	"version":       filters.NewVersionFilter(),
}

var workPackagesCmd = &cobra.Command{
	Use:     "workpackages",
	Aliases: []string{"wps"},
	Short:   "Lists work packages",
	Long:    "Get a list of visible work packages. Filter flags can be applied.",
	Run:     listWorkPackages,
}

func listWorkPackages(_ *cobra.Command, _ []string) {
	// This needs to be removed, once all filters are built the "new" way
	if errorText := validateCommandFlagComposition(); len(errorText) > 0 {
		printer.ErrorText(errorText)
		return
	}

	query, err := buildQuery()
	if err != nil {
		printer.ErrorText(err.Error())
		return
	}

	collection, err := work_packages.All(filterOptions(), query, showTotal)
	switch {
	case err == nil && showTotal:
		printer.Number(collection.Total)
	case err == nil:
		printer.WorkPackages(collection.Items)
	default:
		printer.Error(err)
	}
}

func validateCommandFlagComposition() (errorText string) {
	switch {
	case len(activeFilters["version"].Value()) != 0 && projectId == 0:
		return "Version flag (--version) can only be used in conjunction with projectId flag (-p or --project-id)."
	case len(activeFilters["notVersion"].Value()) != 0 && projectId == 0:
		return "Not version filter flag (--not-version) can only be used in conjunction with projectId flag (-p or --project-id)."
	case len(activeFilters["subProject"].Value()) > 0 || len(activeFilters["notSubProject"].Value()) > 0:
		if !includeSubProjects || projectId == 0 {
			return `Sub project filter flags (--sub-project or --not-sub-project) can only be used
in conjunction with setting the flag --include-sub-projects and setting a
project with the projectId flag (-p or --project-id).`
		}
	}

	return ""
}

func buildQuery() (requests.Query, error) {
	var q requests.Query

	for _, filter := range activeFilters {
		if filter.Value() == filter.DefaultValue() {
			continue
		}

		err := filter.ValidateInput()
		if err != nil {
			return requests.NewEmptyQuery(), err
		}

		q = q.Merge(filter.Query())
	}

	return q, nil
}

func filterOptions() *map[work_packages.FilterOption]string {
	options := make(map[work_packages.FilterOption]string)

	options[work_packages.IncludeSubProjects] = strconv.FormatBool(includeSubProjects)

	if projectId > 0 {
		options[work_packages.Project] = strconv.FormatUint(projectId, 10)
	}

	if len(assignee) > 0 {
		options[work_packages.Assignee] = assignee
	}

	if len(statusFilter) > 0 {
		options[work_packages.Status] = validateFilterValue(work_packages.Status, statusFilter)
	}

	if len(typeFilter) > 0 {
		options[work_packages.Type] = validateFilterValue(work_packages.Type, typeFilter)
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

func validateFilterValue(filter work_packages.FilterOption, value string) string {
	matched, err := regexp.Match(work_packages.InputValidationExpression[filter], []byte(value))
	if err != nil {
		printer.Error(err)
	}

	if !matched {
		printer.ErrorText(fmt.Sprintf("Invalid %s value %s.", filter, printer.Yellow(value)))
		os.Exit(-1)
	}

	return value
}
