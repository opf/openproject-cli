package requests

import (
	"fmt"
	"github.com/opf/openproject-cli/models/types"
	"net/url"
	"strings"
)

type Query struct {
	pageSize int
	filters  []Filter
}

type Filter struct {
	operator string
	name     string
	values   []string
}

func (filter Filter) String() string {
	return fmt.Sprintf(
		"{\"%s\":{\"operator\":\"%s\",\"values\":[\"%s\"]}}",
		filter.name,
		filter.operator,
		strings.Join(filter.values, "\",\""),
	)
}

func (query Query) String() string {
	var filtersQuery = ""
	if len(query.filters) > 0 {
		var fStr = make([]string, len(query.filters))
		for idx, f := range query.filters {
			fStr[idx] = f.String()
		}
		filtersString := fmt.Sprintf("[%s]", strings.Join(fStr, ","))
		filtersQuery = fmt.Sprintf("&filters=%s", url.QueryEscape(filtersString))
	}

	return fmt.Sprintf("pageSize=%d%s", query.pageSize, filtersQuery)
}

func NewNotificationReasonFilter(reason types.Reason) Filter {
	return Filter{
		operator: "=",
		name:     "reason",
		values:   []string{string(reason)},
	}
}

func NewNotificationReadFilter(read bool) Filter {
	var bStr string
	if read {
		bStr = "t"
	} else {
		bStr = "f"
	}

	return Filter{
		operator: "=",
		name:     "readIAN",
		values:   []string{bStr},
	}
}

func NewQuery(filters []Filter) Query {
	return Query{
		pageSize: 100,
		filters:  filters,
	}
}
