package time_entries

import (
	"github.com/opf/openproject-cli/components/requests"
)

type FilterOption int

const (
	User FilterOption = iota
)

func UserFilter(name string) requests.Filter {
	return requests.Filter{
		Operator: "=",
		Name:     "user",
		Values:   []string{name},
	}
}
