package workpackage

import (
	"errors"
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
