package requests

import (
	"fmt"
	"net/url"
	"strings"
)

type Query struct {
	pageSize int
	filters  []Filter
}

type Filter struct {
	Operator string
	Name     string
	Values   []string
}

func (filter Filter) String() string {
	return fmt.Sprintf(
		"{\"%s\":{\"operator\":\"%s\",\"values\":[\"%s\"]}}",
		filter.Name,
		filter.Operator,
		strings.Join(filter.Values, "\",\""),
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

func NewQuery(filters []Filter) Query {
	return Query{
		pageSize: 100,
		filters:  filters,
	}
}

func NewPagedQuery(pageSize int, filters []Filter) Query {
	return Query{
		pageSize: pageSize,
		filters:  filters,
	}
}
