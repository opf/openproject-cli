package printer_test

import (
	"encoding/json"
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

func TestWhoamiList_Text_SeparatesEntries(t *testing.T) {
	testingPrinter.Reset()

	printer.WhoamiList([]printer.WhoamiEntry{
		{Profile: "default", Host: "https://a.example.com", User: &models.User{Id: 1, Name: "A"}},
		{Profile: "work", Host: "https://b.example.com", User: &models.User{Id: 2, Name: "B"}},
	})

	expected := fmt.Sprintf(
		"Profile: %s\nServer:  %s\nUser:    %s %s\n\nProfile: %s\nServer:  %s\nUser:    %s %s\n",
		printer.Yellow("default"),
		printer.Cyan("https://a.example.com"),
		printer.Red("#1"),
		printer.Cyan("A"),
		printer.Yellow("work"),
		printer.Cyan("https://b.example.com"),
		printer.Red("#2"),
		printer.Cyan("B"),
	)

	if testingPrinter.Result != expected {
		t.Errorf("WhoamiList output mismatch\ngot:  %q\nwant: %q", testingPrinter.Result, expected)
	}
}

func TestWhoamiList_Json_EmitsSingleArray(t *testing.T) {
	printer.InitRenderer("json")
	defer printer.InitRenderer("text")
	testingPrinter.Reset()

	printer.WhoamiList([]printer.WhoamiEntry{
		{Profile: "default", Host: "https://a.example.com", User: &models.User{Id: 1, Name: "A"}},
		{Profile: "work", Host: "https://b.example.com", User: &models.User{Id: 2, Name: "B"}},
	})

	var parsed []map[string]any
	if err := json.Unmarshal([]byte(testingPrinter.Result), &parsed); err != nil {
		t.Fatalf("output is not a single JSON array: %v\noutput: %q", err, testingPrinter.Result)
	}
	if len(parsed) != 2 {
		t.Errorf("expected 2 entries, got %d", len(parsed))
	}
	if parsed[1]["profile"] != "work" {
		t.Errorf("second entry profile = %v, want work", parsed[1]["profile"])
	}
}
