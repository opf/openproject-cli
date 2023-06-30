package common

import (
	"strconv"
	"strings"
)

func SanitizeLineBreaks(input string) string {
	return strings.Replace(input, "\n", "", -1)
}

func ParseId(input string) (bool, uint64) {
	var inputAsId = false
	id, err := strconv.ParseUint(input, 10, 64)
	if err == nil {
		inputAsId = true
	}

	return inputAsId, id
}
