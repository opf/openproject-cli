package workpackage

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
)

// newCreateTestCmd builds a fresh *cobra.Command with the same flags
// registered on the real createCmd, bound to the shared package-level
// vars. Using a fresh command per test keeps pflag's Changed() state
// isolated between test cases.
func newCreateTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "create"}
	cmd.Flags().StringVarP(&createProjectId, "project", "p", "", "")
	cmd.Flags().BoolVarP(&createOpenInBrowser, "open", "o", false, "")
	cmd.Flags().StringVarP(&createTypeFlag, "type", "t", "", "")
	cmd.Flags().Uint64Var(&createAssigneeFlag, "assignee", 0, "")
	cmd.Flags().StringVar(&createDescriptionFlag, "description", "", "")
	cmd.Flags().Uint64Var(&createParentID, "parent", 0, "")
	cmd.Flags().BoolVar(&createDryRun, "dry-run", false, "")
	return cmd
}

func resetCreateFlags() {
	createProjectId = ""
	createOpenInBrowser = false
	createTypeFlag = ""
	createAssigneeFlag = 0
	createDescriptionFlag = ""
	createParentID = 0
	createDryRun = false
}

func TestCreateDryRunRendersPlanWithoutPosting(t *testing.T) {
	testingPrinter, requestCount := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
	})

	resetCreateFlags()
	cmd := newCreateTestCmd()
	if err := cmd.Flags().Parse([]string{"--project", "1482", "--parent", "74316", "--dry-run"}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(resetCreateFlags)

	if err := createWorkPackage(cmd, []string{"Draft subject"}); err != nil {
		t.Fatalf("createWorkPackage error = %v, want nil", err)
	}

	if *requestCount != 0 {
		t.Errorf("request count = %d, want 0 (no POST for dry-run)", *requestCount)
	}

	var plan map[string]any
	if err := json.Unmarshal([]byte(testingPrinter.Result), &plan); err != nil {
		t.Fatalf("failed to parse plan output: %v", err)
	}
	workPackage, ok := plan["work_package"].(map[string]any)
	if !ok {
		t.Fatalf("expected work_package object, got %#v", plan["work_package"])
	}
	if workPackage["subject"] != "Draft subject" {
		t.Errorf("expected plan subject %q, got %#v", "Draft subject", workPackage["subject"])
	}
	if plan["parent_id"] != float64(74316) {
		t.Errorf("expected plan parent_id 74316, got %#v", plan["parent_id"])
	}
}

func TestCreateParentIncludesParentLinkInPost(t *testing.T) {
	var postBody map[string]any

	testingPrinter, requestCount := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
		}
		if err := json.NewDecoder(request.Body).Decode(&postBody); err != nil {
			t.Fatal(err)
		}
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"id":42,"displayId":"42","subject":"Created"}`))
	})

	resetCreateFlags()
	cmd := newCreateTestCmd()
	if err := cmd.Flags().Parse([]string{"--project", "1", "--parent", "42"}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(resetCreateFlags)

	if err := createWorkPackage(cmd, []string{"Created"}); err != nil {
		t.Fatalf("createWorkPackage error = %v, want nil", err)
	}

	if *requestCount != 1 {
		t.Errorf("request count = %d, want 1", *requestCount)
	}
	links, ok := postBody["_links"].(map[string]any)
	if !ok {
		t.Fatalf("expected links object in post body, got %#v", postBody["_links"])
	}
	parent, ok := links["parent"].(map[string]any)
	if !ok {
		t.Fatalf("expected parent link object, got %#v", links["parent"])
	}
	if parent["href"] != "/api/v3/work_packages/42" {
		t.Errorf("expected parent href /api/v3/work_packages/42, got %#v", parent["href"])
	}
	if strings.Contains(testingPrinter.ErrResult, "[ERROR]") {
		t.Errorf("expected no error diagnostic, got: %q", testingPrinter.ErrResult)
	}
}

func TestCreateDryRunWithOpenErrors(t *testing.T) {
	_, requestCount := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
	})

	resetCreateFlags()
	cmd := newCreateTestCmd()
	if err := cmd.Flags().Parse([]string{"--project", "1", "--dry-run", "--open"}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(resetCreateFlags)

	err := createWorkPackage(cmd, []string{"Created"})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("createWorkPackage error = %v, want ErrHandled", err)
	}
	if *requestCount != 0 {
		t.Errorf("request count = %d, want 0", *requestCount)
	}
}

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
