package activity

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

func TestListWorkPackageActivitiesReportsFetchErrorsOnce(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		http.Error(response, "activity failure", http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	testingPrinter := initActivityTest(t, server.URL)
	err := listWorkPackageActivities("42")

	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("listWorkPackageActivities error = %v, want ErrHandled", err)
	}
	if count := strings.Count(testingPrinter.ErrResult, "[ERROR]"); count != 1 {
		t.Errorf("error diagnostic count = %d, want 1; stderr: %q", count, testingPrinter.ErrResult)
	}
	if testingPrinter.Result != "" {
		t.Errorf("stdout = %q, want empty", testingPrinter.Result)
	}
}

func TestListWorkPackageActivitiesReportsUserFetchErrorsOnce(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		if request.URL.Path == "/api/v3/work_packages/42/activities" {
			_, _ = response.Write([]byte(`{
				"_embedded": {
					"elements": [{
						"id": 1,
						"comment": {"raw": ""},
						"details": [],
						"_links": {"user": {"href": "/api/v3/users/7"}}
					}]
				}
			}`))
			return
		}
		http.Error(response, "user failure", http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	testingPrinter := initActivityTest(t, server.URL)
	err := listWorkPackageActivities("42")

	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("listWorkPackageActivities error = %v, want ErrHandled", err)
	}
	if count := strings.Count(testingPrinter.ErrResult, "[ERROR]"); count != 1 {
		t.Errorf("error diagnostic count = %d, want 1; stderr: %q", count, testingPrinter.ErrResult)
	}
	if testingPrinter.Result != "" {
		t.Errorf("stdout = %q, want empty", testingPrinter.Result)
	}
}

func initActivityTest(t *testing.T, serverURL string) *printer.TestingPrinter {
	t.Helper()

	host, err := url.Parse(serverURL)
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
	return testingPrinter
}
