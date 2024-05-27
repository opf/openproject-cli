package resources

import "github.com/opf/openproject-cli/components/requests"

func TypeAheadFilter(input string) requests.Filter {
	return requests.Filter{
		Operator: "**",
		Name:     "typeahead",
		Values:   []string{input},
	}
}

type Filter interface {
	Value() string
	ValuePointer() *string
	Name() string
	ShortHand() string
	Usage() string
	ValidateInput() error
	Query() requests.Query
	DefaultValue() string
}
