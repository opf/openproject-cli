package user

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
)

func TestSearchUsersRendersEmptyJSONArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"_embedded":{"elements":[]},"total":0,"count":0}`))
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)

	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)
	if err := printer.InitRenderer("json"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = printer.InitRenderer("text")
	})

	if err := searchUser(nil, []string{"missing person"}); err != nil {
		t.Fatalf("searchUser returned an error: %v", err)
	}

	if testingPrinter.Result != "[]\n" {
		t.Errorf("stdout = %q, want %q", testingPrinter.Result, "[]\n")
	}
	if !strings.Contains(testingPrinter.ErrResult, "missing person") {
		t.Errorf("stderr should contain the complete query, got %q", testingPrinter.ErrResult)
	}
}

func TestSearchUserRequiresOneArgument(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		requestCount++
		_, _ = response.Write([]byte(`{"_embedded":{"elements":[]}}`))
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)

	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)

	if err := searchUser(nil, []string{"missing", "person"}); err == nil {
		t.Fatal("searchUser returned nil for two arguments")
	}
	if requestCount != 0 {
		t.Errorf("request count = %d, want 0", requestCount)
	}
}
