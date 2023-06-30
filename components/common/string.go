package common

import "strings"

func SanitizeLineBreaks(input string) string {
	return strings.Replace(input, "\n", "", -1)
}
