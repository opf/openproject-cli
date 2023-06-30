package printer

import (
	"encoding/json"
	"github.com/opf/openproject-cli/components/errors"
)

type apiErrorModel struct {
	Type            string `json:"_type,omitempty"`
	ErrorIdentifier string `json:"errorIdentifier,omitempty"`
	Messsage        string `json:"message,omitempty"`
}

func Info(msg string) {
	activePrinter.Println(msg)
}

func Done() {
	activePrinter.Println(Green("DONE"))
}

func Error(err error) {
	switch err.(type) {
	case *errors.ResponseError:
		err := err.(*errors.ResponseError)
		responseError(err.Status(), err.Response())
	default:
		activePrinter.Printf("%s Program exited with error: %+v\n", Red("[ERROR]"), err)
	}
}

func ErrorText(msg string) {
	activePrinter.Printf("%s %s\n", Red("[ERROR]"), msg)
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

	activePrinter.Printf(
		"%s Bad response from server: (%d)\n\n%s\n",
		Red("[ERROR]"),
		status,
		bodyRepresentation,
	)
}

func apiError(status int, err apiErrorModel) {
	activePrinter.Printf(
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
