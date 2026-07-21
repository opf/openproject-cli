package project

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
	"github.com/opf/openproject-cli/components/routes"
)

func TestInspectBrowserFailureReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("request method = %s, want GET", request.Method)
		}
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"id":1,"identifier":"example","name":"Example"}`))
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)
	routes.Init(host)

	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)
	t.Setenv("PATH", t.TempDir())
	openInBrowser = true
	t.Cleanup(func() {
		openInBrowser = false
	})

	err = inspectProject(nil, []string{"1"})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("inspectProject error = %v, want ErrHandled", err)
	}
	if count := strings.Count(testingPrinter.ErrResult, "[ERROR]"); count != 1 {
		t.Errorf("error diagnostic count = %d, want 1; stderr: %q", count, testingPrinter.ErrResult)
	}
}
