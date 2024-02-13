package printer

import "strconv"

func Number(number int64) {
	activePrinter.Printf("%s\n", Cyan(strconv.FormatInt(number, 10)))
}
