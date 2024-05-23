package work_packages

import (
	"strings"

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

func StatusFilter(status string) requests.Filter {
	return requests.Filter{
		Operator: "=",
		Name:     "status",
		Values:   strings.Split(status, ","),
	}
}
