package printer

import "fmt"

const (
	reset = "\033[0m"
	cRed  = "\033[31m"
	cCyan = "\033[1;36m"
)

func red(s string) string {
	return fmt.Sprintf("%s%s%s", cRed, s, reset)
}
func cyan(s string) string {
	return fmt.Sprintf("%s%s%s", cCyan, s, reset)
}
