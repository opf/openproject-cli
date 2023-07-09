package common

import (
	"strconv"
	"strings"
)

func SanitizeLineBreaks(input string) string {
	noCRLF := strings.Replace(input, "\r\n", "", -1)
	return strings.Replace(noCRLF, "\n", "", -1)
}

func ParseId(input string) (bool, uint64) {
	var inputAsId = false
	id, err := strconv.ParseUint(input, 10, 64)
	if err == nil {
		inputAsId = true
	}

	return inputAsId, id
}
