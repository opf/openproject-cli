package resources

import "github.com/opf/openproject-cli/components/requests"

func TypeAheadFilter(input string) requests.Filter {
	return requests.Filter{
		Operator: "**",
		Name:     "typeahead",
		Values:   []string{input},
	}
}
