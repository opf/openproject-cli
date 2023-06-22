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

func black(s string) string {
	return fmt.Sprintf("%s%s%s", cBlack, s, reset)
}

func red(s string) string {
	return fmt.Sprintf("%s%s%s", cRed, s, reset)
}

func green(s string) string {
	return fmt.Sprintf("%s%s%s", cGreen, s, reset)
}

func yellow(s string) string {
	return fmt.Sprintf("%s%s%s", cYellow, s, reset)
}

func blue(s string) string {
	return fmt.Sprintf("%s%s%s", cBlue, s, reset)
}

func magenta(s string) string {
	return fmt.Sprintf("%s%s%s", cMagenta, s, reset)
}

func cyan(s string) string {
	return fmt.Sprintf("%s%s%s", cCyan, s, reset)
}

func white(s string) string {
	return fmt.Sprintf("%s%s%s", cWhite, s, reset)
}
