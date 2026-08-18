package version

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
)

func TestListVersionsInvalidProjectReturnsError(t *testing.T) {
	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)
	listProjectId = ""
	t.Cleanup(func() { listProjectId = "" })

	err := listVersions(nil, nil)
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("listVersions error = %v, want ErrHandled", err)
	}
	if count := strings.Count(testingPrinter.ErrResult, "[ERROR]"); count != 1 {
		t.Errorf("error diagnostic count = %d, want 1; stderr: %q", count, testingPrinter.ErrResult)
	}
}

func TestListVersionsPrintsAvailableVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{
			"_type": "Collection",
			"_embedded": {
				"elements": [
					{"id": 3, "name": "v17"},
					{"id": 7, "name": "v42"}
				]
			}
		}`))
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)

	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)
	listProjectId = "example"
	t.Cleanup(func() { listProjectId = "" })

	if err := listVersions(nil, nil); err != nil {
		t.Fatalf("listVersions error = %v, want nil", err)
	}

	if !strings.Contains(testingPrinter.Result, "v17") || !strings.Contains(testingPrinter.Result, "v42") {
		t.Errorf("expected output to contain versions 'v17' and 'v42', got: %q", testingPrinter.Result)
	}
}
