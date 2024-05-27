package requests

import (
	"fmt"
	"net/url"
	"strings"
)

type Query struct {
	filters    []Filter
	attributes map[string]string
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
	queryStr := filtersQueryAttribute(query.filters)
	for key, value := range query.attributes {
		if len(queryStr) > 0 {
			queryStr += "&"
		}
		queryStr += fmt.Sprintf("%s=%s", key, url.QueryEscape(value))
	}

	return queryStr
}

func filtersQueryAttribute(filters []Filter) string {
	if len(filters) == 0 {
		return ""
	}

	var fStr = make([]string, len(filters))
	for idx, f := range filters {
		fStr[idx] = f.String()
	}

	filtersString := fmt.Sprintf("[%s]", strings.Join(fStr, ","))
	return fmt.Sprintf("filters=%s", url.QueryEscape(filtersString))
}

func NewQuery(attributes map[string]string, filters []Filter) Query {
	if attributes == nil {
		return Query{attributes: make(map[string]string), filters: filters}
	}

	return Query{attributes: attributes, filters: filters}
}

func NewFilterQuery(filters []Filter) Query {
	attributes := map[string]string{
		"pageSize": "100",
	}

	return Query{attributes: attributes, filters: filters}
}

func NewUnpaginatedQuery(attributes map[string]string, filters []Filter) Query {
	var attr = attributes
	if attr == nil {
		attr = make(map[string]string)
	}

	attr["pageSize"] = "-1"
	return Query{attributes: attr, filters: filters}
}

func NewPaginatedQuery(pageSize int, filters []Filter) Query {
	attributes := map[string]string{
		"pageSize": fmt.Sprintf("%d", pageSize),
	}

	return Query{attributes: attributes, filters: filters}
}
