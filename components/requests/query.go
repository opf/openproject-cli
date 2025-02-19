package requests

import (
	"fmt"
	"github.com/opf/openproject-cli/components/common"
	"net/url"
	"slices"
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

func (filter Filter) Equals(other Filter) bool {
	valuesEqual := true

	for _, value := range filter.Values {
		if !slices.Contains(other.Values, value) {
			valuesEqual = false
			break
		}
	}

	return filter.Operator == other.Operator && filter.Name == other.Name && valuesEqual
}

func (query Query) Merge(another Query) Query {
	filters := append(query.filters, another.filters...)

	attributes := query.attributes
	if attributes == nil {
		attributes = make(map[string]string)
	}

	for key, value := range another.attributes {
		attributes[key] = value
	}

	return Query{
		attributes: attributes,
		filters:    filters,
	}
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

func (query Query) Equals(other Query) bool {
	filtersEqual := common.All(
		query.filters,
		func(filter Filter) bool {
			filterExists := false
			for _, f := range other.filters {
				if filter.Equals(f) {
					filterExists = true
					break
				}
			}

			return filterExists
		})

	attributesEqual := true
	for idx, value := range query.attributes {
		if other.attributes[idx] != value {
			attributesEqual = false
			break
		}
	}

	return filtersEqual && attributesEqual
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

func NewEmptyQuery() Query {
	return Query{attributes: make(map[string]string), filters: []Filter{}}
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
