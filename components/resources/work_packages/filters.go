package work_packages

import (
	"github.com/opf/openproject-cli/components/requests"
)

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
