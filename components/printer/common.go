package printer

import (
	"fmt"
	"os"
)

func Error(err error) {
	fmt.Printf("%s Program exited with error: %+v", red("[ERROR]"), err)
	os.Exit(-1)
}

func ErrorText(msg string) {
	fmt.Printf("%s %s", red("[ERROR]"), msg)
	os.Exit(-1)
}

func ResponseError(status int, body []byte) {
	var bodyRepresentation string
	if len(body) >= 256 {
		bodyRepresentation = string(body[:256]) + "\n..."
	} else {
		bodyRepresentation = string(body)
	}

	fmt.Printf(
		"%s Bad response from server: (%d)\n\n%s",
		red("[ERROR]"),
		status,
		bodyRepresentation,
	)
	os.Exit(-1)
}
