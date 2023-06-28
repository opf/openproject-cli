package users

import (
	"strconv"

	"github.com/opf/openproject-cli/components/requests"
)

func IdFilter(ids []uint64) requests.Filter {
	idsStr := make([]string, len(ids))
	for idx, id := range ids {
		idsStr[idx] = strconv.FormatUint(id, 10)
	}

	return requests.Filter{
		Operator: "=",
		Name:     "id",
		Values:   idsStr,
	}
}

func NameFilter(name string) requests.Filter {
	return requests.Filter{
		Operator: "~",
		Name:     "name",
		Values:   []string{name},
	}
}
