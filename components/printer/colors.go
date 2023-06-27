package printer

import "fmt"

const (
	cBlack   = "\033[30m"
	cRed     = "\033[31m"
	cGreen   = "\033[32m"
	cYellow  = "\033[33m"
	cBlue    = "\033[34m"
	cMagenta = "\033[35m"
	cCyan    = "\033[36m"
	cWhite   = "\033[37m"
	reset    = "\033[0m"
)

func Black(s string) string {
	return fmt.Sprintf("%s%s%s", cBlack, s, reset)
}

func Red(s string) string {
	return fmt.Sprintf("%s%s%s", cRed, s, reset)
}

func Green(s string) string {
	return fmt.Sprintf("%s%s%s", cGreen, s, reset)
}

func Yellow(s string) string {
	return fmt.Sprintf("%s%s%s", cYellow, s, reset)
}

func Blue(s string) string {
	return fmt.Sprintf("%s%s%s", cBlue, s, reset)
}

func Magenta(s string) string {
	return fmt.Sprintf("%s%s%s", cMagenta, s, reset)
}

func Cyan(s string) string {
	return fmt.Sprintf("%s%s%s", cCyan, s, reset)
}

func White(s string) string {
	return fmt.Sprintf("%s%s%s", cWhite, s, reset)
}
