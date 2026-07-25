package workpackage

import (
	"net/http"
	"strings"
	"testing"
)

func TestCreateBrowserFailureKeepsSuccessfulExit(t *testing.T) {
	testingPrinter, requestCount := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Errorf("request method = %s, want POST", request.Method)
		}
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"id":42,"displayId":"42","subject":"Created"}`))
	})

	t.Setenv("PATH", t.TempDir())
	createProjectId = "1"
	createOpenInBrowser = true
	createTypeFlag = ""
	createAssigneeFlag = 0
	createDescriptionFlag = ""
	t.Cleanup(func() {
		createProjectId = ""
		createOpenInBrowser = false
		createTypeFlag = ""
		createAssigneeFlag = 0
		createDescriptionFlag = ""
	})

	if err := createWorkPackage(createCmd, []string{"Created"}); err != nil {
		t.Fatalf("createWorkPackage returned an error after the API create succeeded: %v", err)
	}

	if *requestCount != 1 {
		t.Errorf("request count = %d, want 1", *requestCount)
	}
	if !strings.Contains(testingPrinter.ErrResult, "[WARNING]") ||
		!strings.Contains(testingPrinter.ErrResult, "Error opening browser") {
		t.Errorf("stderr should warn about the browser failure, got %q", testingPrinter.ErrResult)
	}
}
