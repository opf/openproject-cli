package work_packages

import (
	"strconv"

	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func Lookup(id uint64) (*models.WorkPackage, error) {
	workPackage, err := fetch(id)
	if err != nil {
		return nil, err
	}

	return workPackage.Convert(), nil
}

func All(filterOptions *map[FilterOption]string, query requests.Query, showOnlyTotal bool) (*models.WorkPackageCollection, error) {
	var filters []requests.Filter
	var projectId *uint64
	var queryAttributes = make(map[string]string)

	for updateOpt, value := range *filterOptions {
		switch updateOpt {
		case IncludeSubProjects:
			queryAttributes["includeSubprojects"] = value
		case Assignee:
			filters = append(filters, AssigneeFilter(value))
		case Status:
			filters = append(filters, StatusFilter(value))
		case Type:
			filters = append(filters, TypeFilter(value))
		case Project:
			n, _ := strconv.ParseUint(value, 10, 64)
			projectId = &n
		}
	}

	if showOnlyTotal {
		queryAttributes["pageSize"] = "0"
	} else {
		queryAttributes["pageSize"] = "-1"
	}

	legacyQuery := requests.NewQuery(queryAttributes, filters)
	newQuery := legacyQuery.Merge(query)

	requestUrl := paths.WorkPackages()

	if projectId != nil {
		requestUrl = paths.ProjectWorkPackages(*projectId)
	}

	response, err := requests.Get(requestUrl, &newQuery)
	if err != nil {
		return nil, err
	}

	workPackageCollection := parser.Parse[dtos.WorkPackageCollectionDto](response)
	return workPackageCollection.Convert(), nil
}

func AvailableTypes(id uint64) ([]*models.Type, error) {
	workPackageDto, err := fetch(id)
	if err != nil {
		return nil, err
	}

	types, err := availableTypes(workPackageDto.Links.Project)
	if err != nil {
		return nil, err
	}

	return types.Convert(), nil
}

func fetch(id uint64) (*dtos.WorkPackageDto, error) {
	response, err := requests.Get(paths.WorkPackage(id), nil)
	if err != nil {
		return nil, err
	}

	workPackage := parser.Parse[dtos.WorkPackageDto](response)
	return &workPackage, nil
}
