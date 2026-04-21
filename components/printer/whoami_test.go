package printer_test

import (
	"fmt"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func TestWhoami(t *testing.T) {
	testingPrinter.Reset()

	user := &models.User{Id: 42, Name: "Jane Doe"}
	printer.Whoami("default", "https://example.com", user)

	expected := fmt.Sprintf(
		"Profile: %s\nServer:  %s\nUser:    %s %s\n",
		printer.Yellow("default"),
		printer.Cyan("https://example.com"),
		printer.Red("#42"),
		printer.Cyan("Jane Doe"),
	)

	if testingPrinter.Result != expected {
		t.Errorf("Whoami output mismatch\ngot:  %q\nwant: %q", testingPrinter.Result, expected)
	}
}
