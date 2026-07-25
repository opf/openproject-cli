package printer

import (
	"encoding/json"
	"fmt"

	"github.com/opf/openproject-cli/components/errors"
)

type apiErrorModel struct {
	Type            string `json:"_type,omitempty"`
	ErrorIdentifier string `json:"errorIdentifier,omitempty"`
	Messsage        string `json:"message,omitempty"`
}

// Output writes requested data to standard out. Use it (or a renderer) for
// command results only; diagnostics belong on standard error.
func Output(msg string) {
	activePrinter.Println(msg)
}

// Info writes a progress or context message to standard error.
func Info(msg string) {
	activePrinter.Eprintln(msg)
}

func Input(prompt string) {
	activePrinter.Eprintf(prompt)
}

// Warning writes a non-fatal diagnostic to standard error so it never corrupts
// machine-readable output (e.g. JSON) written to standard out.
func Warning(msg string) {
	activePrinter.Eprintln(fmt.Sprintf("%s %s", Yellow("[WARNING]"), msg))
}

func Done() {
	activePrinter.Eprintln(Green("DONE"))
}

func Error(err error) {
	switch err.(type) {
	case *errors.ResponseError:
		err := err.(*errors.ResponseError)
		responseError(err.Status(), err.Response())
	default:
		activePrinter.Eprintf("%s Program exited with error: %+v\n", Red("[ERROR]"), err)
	}
}

func Debug(verboseFlag bool, msg string) {
	if verboseFlag {
		activePrinter.Eprintf("%s %s\n", Magenta("[DEBUG]"), msg)
	}
}

func ErrorText(msg string) {
	activePrinter.Eprintf("%s %s\n", Red("[ERROR]"), msg)
}

func responseError(status int, body []byte) {
	var apiErr apiErrorModel
	if err := json.Unmarshal(body, &apiErr); err == nil {
		apiError(status, apiErr)
		return
	}

	var bodyRepresentation string
	if len(body) >= 256 {
		bodyRepresentation = string(body[:256]) + "\n..."
	} else {
		bodyRepresentation = string(body)
	}

	activePrinter.Eprintf(
		"%s Bad response from server: (%d)\n\n%s\n",
		Red("[ERROR]"),
		status,
		bodyRepresentation,
	)
}

func apiError(status int, err apiErrorModel) {
	activePrinter.Eprintf(
		"%s API request failure (%d): %s\n%s\n",
		Red("[ERROR]"),
		status,
		Yellow(err.ErrorIdentifier),
		err.Messsage,
	)
}

func indent(spaces int) (res string) {
	for len(res) < spaces {
		res += " "
	}
	return res
}
