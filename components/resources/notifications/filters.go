package notifications

import "github.com/opf/openproject-cli/components/requests"

func ReasonFilter(reason string) requests.Filter {
	return requests.Filter{
		Operator: "=",
		Name:     "reason",
		Values:   []string{reason},
	}
}

func ReadFilter(read bool) requests.Filter {
	var bStr string
	if read {
		bStr = "t"
	} else {
		bStr = "f"
	}

	return requests.Filter{
		Operator: "=",
		Name:     "readIAN",
		Values:   []string{bStr},
	}
}
