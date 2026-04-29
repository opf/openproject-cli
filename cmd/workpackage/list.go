package workpackage

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"

	opErrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources"
	"github.com/opf/openproject-cli/components/resources/projects"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/components/resources/work_packages/filters"
	"github.com/opf/openproject-cli/models"
)

var listAssignee string
var listParentId uint64
var listProjectId string
var listShowTotal bool
var listStatusFilter string
var listTypeFilter string
var listIncludeSubProjects bool

var activeFilters = map[string]resources.Filter{
	"notSubProject": filters.NewNotSubProjectFilter(),
	"notVersion":    filters.NewNotVersionFilter(),
	"subProject":    filters.NewSubProjectFilter(),
	"timestamp":     filters.NewTimestampFilter(),
	"version":       filters.NewVersionFilter(),
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists work packages",
	Long:  "Get a list of visible work packages. Filter flags can be applied.",
	Run:   listWorkPackages,
}

func listWorkPackages(_ *cobra.Command, _ []string) {
	if errorText := validateCommandFlagComposition(); len(errorText) > 0 {
		printer.ErrorText(errorText)
		return
	}

	if len(listProjectId) > 0 {
		if err := projects.ValidateIdentifier(listProjectId); err != nil {
			printer.ErrorText(fmt.Sprintf("--project: %s", err.Error()))
			return
		}
	}

	if listParentId > 0 {
		if _, err := work_packages.Lookup(listParentId); err != nil {
			printer.ErrorText(fmt.Sprintf("--parent-id: work package #%d not found.", listParentId))
			return
		}
	}

	query, err := buildQuery()
	if err != nil {
		printer.ErrorText(err.Error())
		return
	}

	collection, err := work_packages.All(filterOptions(), query, listShowTotal)
	switch {
	case err == nil && listShowTotal:
		printer.Number(collection.Total)
	case err == nil:
		printer.WorkPackages(collection.Items)
	case isNotFound(err) && len(listProjectId) > 0:
		printer.ErrorText(fmt.Sprintf("--project: no project found with identifier or ID '%s'", listProjectId))
	default:
		printer.Error(err)
	}
}

func validateCommandFlagComposition() (errorText string) {
	switch {
	case len(activeFilters["version"].Value()) != 0 && len(listProjectId) == 0:
		return "Version flag (--version) can only be used in conjunction with project flag (-p or --project)."
	case len(activeFilters["notVersion"].Value()) != 0 && len(listProjectId) == 0:
		return "Not version filter flag (--not-version) can only be used in conjunction with project flag (-p or --project)."
	case len(activeFilters["subProject"].Value()) > 0 || len(activeFilters["notSubProject"].Value()) > 0:
		if !listIncludeSubProjects || len(listProjectId) == 0 {
			return `Sub project filter flags (--sub-project or --not-sub-project) can only be used
in conjunction with setting the flag --include-sub-projects and setting a
project with the project flag (-p or --project).`
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

	options[work_packages.IncludeSubProjects] = strconv.FormatBool(listIncludeSubProjects)

	if listParentId > 0 {
		options[work_packages.Parent] = strconv.FormatUint(listParentId, 10)
	}

	if len(listProjectId) > 0 {
		options[work_packages.Project] = listProjectId
	}

	if len(listAssignee) > 0 {
		options[work_packages.Assignee] = listAssignee
	}

	if len(listStatusFilter) > 0 {
		options[work_packages.Status] = validateFilterValue(work_packages.Status, listStatusFilter)
	}

	if len(listTypeFilter) > 0 {
		options[work_packages.Type] = validateFilterValue(work_packages.Type, listTypeFilter)
	}

	return &options
}

func validatedVersionId(version string) string {
	project, err := projects.Lookup(listProjectId)
	if err != nil {
		printer.Error(err)
	}

	versions, err := projects.AvailableVersions(project.Identifier)
	if err != nil {
		printer.Error(err)
	}

	filteredVersions := common.Filter(versions, func(v *models.Version) bool {
		return v.Name == version
	})

	if len(filteredVersions) != 1 {
		printer.Info(fmt.Sprintf(
			"No unique available version from input %s found for project %s. Please use one of the versions listed below.",
			printer.Cyan(version),
			printer.Red(project.Identifier),
		))

		printer.Versions(versions)

		os.Exit(-1)
	}

	return strconv.FormatUint(filteredVersions[0].Id, 10)
}

func isNotFound(err error) bool {
	if respErr, ok := err.(*opErrors.ResponseError); ok {
		return respErr.Status() == 404
	}
	return false
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
