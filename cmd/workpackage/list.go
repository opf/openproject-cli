package workpackage

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources"
	"github.com/opf/openproject-cli/components/resources/projects"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/components/resources/work_packages/filters"
)

var listAssignee string
var listParentId string
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
	RunE:  listWorkPackages,
}

func listWorkPackages(_ *cobra.Command, _ []string) error {
	if errorText := validateCommandFlagComposition(); len(errorText) > 0 {
		printer.ErrorText(errorText)
		return openerrors.ErrHandled
	}

	if len(listProjectId) > 0 {
		if err := projects.ValidateIdentifier(listProjectId); err != nil {
			printer.ErrorText(fmt.Sprintf("--project: %s", err.Error()))
			return openerrors.ErrHandled
		}
	}

	if len(listParentId) > 0 {
		if err := work_packages.ValidateIdentifier(listParentId); err != nil {
			printer.ErrorText(fmt.Sprintf("--parent-id: %s", err.Error()))
			return openerrors.ErrHandled
		}
		parentWp, err := work_packages.Lookup(listParentId)
		if err != nil {
			if isNotFound(err) {
				printer.ErrorText(fmt.Sprintf("--parent-id: work package %s not found.", listParentId))
			} else {
				printer.Error(err)
			}
			return openerrors.ErrHandled
		}
		listParentId = fmt.Sprintf("%d", parentWp.Id)
	}

	query, err := buildQuery()
	if err != nil {
		printer.ErrorText(err.Error())
		return openerrors.ErrHandled
	}

	options, err := filterOptions()
	if err != nil {
		printer.ErrorText(err.Error())
		return openerrors.ErrHandled
	}

	collection, err := work_packages.All(options, query, listShowTotal)
	switch {
	case err == nil && listShowTotal:
		printer.Number(collection.Total)
	case err == nil:
		printer.WorkPackages(collection.Items)
	case isNotFound(err) && len(listProjectId) > 0:
		printer.ErrorText(fmt.Sprintf("--project: no project found with identifier or ID '%s'", listProjectId))
		return openerrors.ErrHandled
	default:
		printer.Error(err)
		return openerrors.ErrHandled
	}
	return nil
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

func filterOptions() (*map[work_packages.FilterOption]string, error) {
	options := make(map[work_packages.FilterOption]string)

	options[work_packages.IncludeSubProjects] = strconv.FormatBool(listIncludeSubProjects)

	if len(listParentId) > 0 {
		options[work_packages.Parent] = listParentId
	}

	if len(listProjectId) > 0 {
		options[work_packages.Project] = listProjectId
	}

	if len(listAssignee) > 0 {
		options[work_packages.Assignee] = listAssignee
	}

	if len(listStatusFilter) > 0 {
		value, err := validateFilterValue(work_packages.Status, listStatusFilter)
		if err != nil {
			return nil, err
		}
		options[work_packages.Status] = value
	}

	if len(listTypeFilter) > 0 {
		value, err := validateFilterValue(work_packages.Type, listTypeFilter)
		if err != nil {
			return nil, err
		}
		options[work_packages.Type] = value
	}

	return &options, nil
}

func validateFilterValue(filter work_packages.FilterOption, value string) (string, error) {
	matched, err := regexp.Match(work_packages.InputValidationExpression[filter], []byte(value))
	if err != nil {
		return "", err
	}

	if !matched {
		return "", fmt.Errorf("invalid %s value %s", filter, printer.Yellow(value))
	}

	return value, nil
}
