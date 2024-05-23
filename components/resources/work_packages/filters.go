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
