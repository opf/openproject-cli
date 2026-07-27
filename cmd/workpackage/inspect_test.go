package workpackage

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	openerrors "github.com/opf/openproject-cli/components/errors"
)

func TestInspectBrowserFailureReturnsError(t *testing.T) {
	testingPrinter, requestCount := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("request method = %s, want GET", request.Method)
		}
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"id":42,"displayId":"42","subject":"Example"}`))
	})

	t.Setenv("PATH", t.TempDir())
	inspectOpenInBrowser = true
	t.Cleanup(func() {
		inspectOpenInBrowser = false
	})

	err := inspectWorkPackage(nil, []string{"42"})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("inspectWorkPackage error = %v, want ErrHandled", err)
	}
	if *requestCount != 1 {
		t.Errorf("request count = %d, want 1", *requestCount)
	}
	if count := strings.Count(testingPrinter.ErrResult, "[ERROR]"); count != 1 {
		t.Errorf("error diagnostic count = %d, want 1; stderr: %q", count, testingPrinter.ErrResult)
	}
}

func TestInspectChildrenRendersPayload(t *testing.T) {
	testingPrinter, _ := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.URL.Path == "/api/v3/work_packages/42":
			_, _ = io.WriteString(response, `{
				"id": 42,
				"subject": "Parent",
				"_links": {
					"self": {"href": "/api/v3/work_packages/42"},
					"project": {"href": "/api/v3/projects/1", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Feature"}
				}
			}`)
		case request.URL.Path == "/api/v3/work_packages/schemas/1-6":
			_, _ = io.WriteString(response, `{}`)
		case request.URL.Path == "/api/v3/work_packages":
			if !strings.Contains(request.URL.RawQuery, "parent") {
				t.Fatalf("expected parent filter in query: %s", request.URL.RawQuery)
			}
			_, _ = io.WriteString(response, `{
				"_embedded": {
					"elements": [
						{
							"id": 43,
							"subject": "Child",
							"_links": {
								"type": {"title": "Task"},
								"status": {"title": "new"}
							}
						}
					]
				}
			}`)
		default:
			t.Fatalf("unexpected path: %s", request.URL.Path)
		}
	})

	inspectWithChildren = true
	t.Cleanup(func() {
		inspectWithChildren = false
	})

	err := inspectWorkPackage(nil, []string{"42"})
	if err != nil {
		t.Fatalf("inspectWorkPackage error = %v, want nil", err)
	}

	if !strings.Contains(testingPrinter.Result, `"id": 43`) {
		t.Errorf("expected output to include child row, got: %q", testingPrinter.Result)
	}
	if !strings.Contains(testingPrinter.Result, `"subject": "Child"`) {
		t.Errorf("expected output to include child subject, got: %q", testingPrinter.Result)
	}
}

func TestInspectChildrenWithOpenErrors(t *testing.T) {
	testingPrinter, requestCount := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		t.Fatalf("unexpected request: %s", request.URL.Path)
	})

	inspectWithChildren = true
	inspectOpenInBrowser = true
	t.Cleanup(func() {
		inspectWithChildren = false
		inspectOpenInBrowser = false
	})

	err := inspectWorkPackage(nil, []string{"42"})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("inspectWorkPackage error = %v, want ErrHandled", err)
	}
	if *requestCount != 0 {
		t.Errorf("request count = %d, want 0", *requestCount)
	}
	if !strings.Contains(testingPrinter.ErrResult, "cannot use --children together with --open or --types") {
		t.Errorf("expected error message about --children, got: %q", testingPrinter.ErrResult)
	}
}
