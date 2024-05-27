package work_packages

import (
	"strings"

	"github.com/opf/openproject-cli/components/requests"
)

type FilterOption int

const (
	Assignee FilterOption = iota
	Version
	Project
	Status
	Type
	IncludeSubProjects
)

var InputValidationExpression = map[FilterOption]string{
	Status: "^(open)$|^(closed)$|^(!?[0-9,]+)$",
	Type:   "^(!?[0-9,]+)$",
}

func (f FilterOption) String() string {
	switch f {
	case Assignee:
		return "assignee"
	case Version:
		return "version"
	case Project:
		return "project"
	case Status:
		return "status"
	case Type:
		return "type"
	case IncludeSubProjects:
		return "include-sub-projects"
	default:
		return "filter"
	}
}

func AssigneeFilter(name string) requests.Filter {
	return requests.Filter{
		Operator: "=",
		Name:     "assignee",
		Values:   []string{name},
	}
}

func VersionFilter(version string) requests.Filter {
	return requests.Filter{
		Operator: "=",
		Name:     "version",
		Values:   []string{version},
	}
}

func TypeFilter(workPackageType string) requests.Filter {
	var operator string
	var values []string

	switch {
	case strings.Index(workPackageType, "!") == 0:
		operator = "!"
		values = strings.Split(workPackageType[1:], ",")
	default:
		operator = "="
		values = strings.Split(workPackageType, ",")
	}

	return requests.Filter{
		Operator: operator,
		Name:     "type",
		Values:   values,
	}
}

func StatusFilter(status string) requests.Filter {
	var operator string
	var values []string

	switch {
	case status == "open":
		operator = "o"
		values = []string{}
	case status == "closed":
		operator = "c"
		values = []string{}
	case strings.Index(status, "!") == 0:
		operator = "!"
		values = strings.Split(status[1:], ",")
	default:
		operator = "="
		values = strings.Split(status, ",")
	}

	return requests.Filter{
		Operator: operator,
		Name:     "status",
		Values:   values,
	}
}
