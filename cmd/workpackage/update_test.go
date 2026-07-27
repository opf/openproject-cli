package workpackage

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
)

// newUpdateTestCmd builds a fresh *cobra.Command with the same flags
// registered on the real updateCmd, bound to the shared package-level
// vars. Using a fresh command per test keeps pflag's Changed() state
// isolated between test cases.
func newUpdateTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "update"}
	cmd.Flags().StringVarP(&updateActionFlag, "action", "a", "", "")
	cmd.Flags().Uint64Var(&updateAssigneeFlag, "assignee", 0, "")
	cmd.Flags().StringVar(&updateAttachFlag, "attach", "", "")
	cmd.Flags().StringVar(&updateDescriptionFlag, "description", "", "")
	cmd.Flags().StringVar(&updateSubjectFlag, "subject", "", "")
	cmd.Flags().StringVarP(&updateTypeFlag, "type", "t", "", "")
	cmd.Flags().StringVar(&updateStatusFlag, "status", "", "")
	cmd.Flags().StringArrayVar(&updateSetFlags, "set", nil, "")
	cmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "")
	return cmd
}

func resetUpdateFlags() {
	updateActionFlag = ""
	updateAssigneeFlag = 0
	updateAttachFlag = ""
	updateDescriptionFlag = ""
	updateSubjectFlag = ""
	updateTypeFlag = ""
	updateStatusFlag = ""
	updateSetFlags = nil
	updateDryRun = false
}

func TestUpdateSetDryRunRendersPlanWithoutPatching(t *testing.T) {
	testingPrinter, requestCount := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/work_packages/42":
			_, _ = io.WriteString(response, `{
				"id": 42,
				"subject": "Example",
				"lockVersion": 3,
				"_links": {
					"self": {"href": "/api/v3/work_packages/42"},
					"project": {"href": "/api/v3/projects/1", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "New"},
					"type": {"href": "/api/v3/types/6", "title": "Feature"}
				}
			}`)
		case request.URL.Path == "/api/v3/work_packages/schemas/1-6":
			_, _ = io.WriteString(response, `{
				"customField130": {"name": "Story points", "type": "Integer", "writable": true}
			}`)
		default:
			t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
		}
	})

	resetUpdateFlags()
	cmd := newUpdateTestCmd()
	if err := cmd.Flags().Parse([]string{"--set", "Story points=5", "--dry-run"}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(resetUpdateFlags)

	if err := updateWorkPackage(cmd, []string{"42"}); err != nil {
		t.Fatalf("updateWorkPackage error = %v, want nil", err)
	}

	if *requestCount != 2 {
		t.Errorf("request count = %d, want 2 (fetch + schema, no patch)", *requestCount)
	}
	if !strings.Contains(testingPrinter.Result, `"resolved_fields"`) {
		t.Errorf("expected update plan output, got: %q", testingPrinter.Result)
	}
	var plan map[string]any
	if err := json.Unmarshal([]byte(testingPrinter.Result), &plan); err != nil {
		t.Fatalf("failed to parse plan output: %v", err)
	}
	resolved, ok := plan["resolved_fields"].(map[string]any)
	if !ok {
		t.Fatalf("expected resolved_fields object, got %#v", plan["resolved_fields"])
	}
	if _, ok := resolved["Story points"]; !ok {
		t.Errorf("expected resolved field for Story points, got %#v", resolved)
	}
}

func TestUpdateSetLivePatchesAndRendersDetails(t *testing.T) {
	var patchBody map[string]any

	testingPrinter, requestCount := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/work_packages/42":
			_, _ = io.WriteString(response, `{
				"id": 42,
				"subject": "Example",
				"lockVersion": 3,
				"_links": {
					"self": {"href": "/api/v3/work_packages/42"},
					"project": {"href": "/api/v3/projects/1", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "New"},
					"type": {"href": "/api/v3/types/6", "title": "Feature"}
				}
			}`)
		case request.Method == http.MethodPatch && request.URL.Path == "/api/v3/work_packages/42":
			if err := json.NewDecoder(request.Body).Decode(&patchBody); err != nil {
				t.Fatal(err)
			}
			response.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(response, `{}`)
		case request.URL.Path == "/api/v3/work_packages/schemas/1-6":
			_, _ = io.WriteString(response, `{
				"customField130": {"name": "Story points", "type": "Integer", "writable": true}
			}`)
		default:
			t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
		}
	})

	resetUpdateFlags()
	cmd := newUpdateTestCmd()
	if err := cmd.Flags().Parse([]string{"--set", "Story points=5"}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(resetUpdateFlags)

	if err := updateWorkPackage(cmd, []string{"42"}); err != nil {
		t.Fatalf("updateWorkPackage error = %v, want nil", err)
	}

	if patchBody["customField130"] != float64(5) {
		t.Errorf("expected patch body to include customField130=5, got %#v", patchBody)
	}
	// UpdateFields: fetch + schema + patch (3). Inspect: fetch + schema (2).
	if *requestCount != 5 {
		t.Errorf("request count = %d, want 5 (fetch, schema, patch, fetch, schema)", *requestCount)
	}
	if !strings.Contains(testingPrinter.Result, `"subject"`) {
		t.Errorf("expected work package details output, got: %q", testingPrinter.Result)
	}
}

func TestUpdateSetCombinedWithOtherFlagsErrors(t *testing.T) {
	testingPrinter, requestCount := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
	})

	resetUpdateFlags()
	cmd := newUpdateTestCmd()
	if err := cmd.Flags().Parse([]string{"--set", "Story points=5", "--subject", "New subject"}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(resetUpdateFlags)

	err := updateWorkPackage(cmd, []string{"42"})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("updateWorkPackage error = %v, want ErrHandled", err)
	}
	if *requestCount != 0 {
		t.Errorf("request count = %d, want 0", *requestCount)
	}
	if !strings.Contains(testingPrinter.ErrResult, "cannot combine --set with other update flags") {
		t.Errorf("stderr = %q, want combine-flags diagnostic", testingPrinter.ErrResult)
	}
}

func TestUpdateStatusPatchesWorkPackage(t *testing.T) {
	var patchBody map[string]any

	testingPrinter, requestCount := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/work_packages/42":
			_, _ = io.WriteString(response, `{
				"id": 42,
				"subject": "Example",
				"lockVersion": 3,
				"_links": {
					"self": {"href": "/api/v3/work_packages/42"},
					"project": {"href": "/api/v3/projects/1", "title": "CLI"},
					"status": {"href": "/api/v3/statuses/1", "title": "New"}
				}
			}`)
		case request.Method == http.MethodPatch && request.URL.Path == "/api/v3/work_packages/42":
			if err := json.NewDecoder(request.Body).Decode(&patchBody); err != nil {
				t.Fatal(err)
			}
			response.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(response, `{}`)
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/statuses":
			_, _ = io.WriteString(response, `{
				"_embedded": {
					"elements": [
						{"id": 1, "name": "New"},
						{"id": 2, "name": "In development"}
					]
				}
			}`)
		default:
			t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
		}
	})

	resetUpdateFlags()
	cmd := newUpdateTestCmd()
	if err := cmd.Flags().Parse([]string{"--status", "In development"}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(resetUpdateFlags)

	if err := updateWorkPackage(cmd, []string{"42"}); err != nil {
		t.Fatalf("updateWorkPackage error = %v, want nil", err)
	}

	links, ok := patchBody["_links"].(map[string]any)
	if !ok {
		t.Fatalf("expected links object in patch body, got %#v", patchBody["_links"])
	}
	status, ok := links["status"].(map[string]any)
	if !ok {
		t.Fatalf("expected status link object, got %#v", links["status"])
	}
	if status["href"] != "/api/v3/statuses/2" {
		t.Errorf("expected status href /api/v3/statuses/2, got %#v", status["href"])
	}
	if *requestCount == 0 {
		t.Errorf("expected requests to have been made")
	}
	if strings.Contains(testingPrinter.ErrResult, "[ERROR]") {
		t.Errorf("expected no error diagnostic, got: %q", testingPrinter.ErrResult)
	}
}

func TestUpdateStatusDryRunRendersResolvedStatusWithoutPatching(t *testing.T) {
	testingPrinter, requestCount := initWorkPackageTestServer(t, func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/work_packages/42":
			_, _ = io.WriteString(response, `{
				"id": 42,
				"subject": "Example",
				"lockVersion": 3,
				"_links": {
					"self": {"href": "/api/v3/work_packages/42"},
					"project": {"href": "/api/v3/projects/1", "title": "CLI"},
					"status": {"href": "/api/v3/statuses/1", "title": "New"}
				}
			}`)
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/statuses":
			_, _ = io.WriteString(response, `{
				"_embedded": {
					"elements": [
						{"id": 1, "name": "New"},
						{"id": 2, "name": "In development"}
					]
				}
			}`)
		default:
			t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
		}
	})

	resetUpdateFlags()
	cmd := newUpdateTestCmd()
	if err := cmd.Flags().Parse([]string{"--status", "in development", "--dry-run"}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(resetUpdateFlags)

	if err := updateWorkPackage(cmd, []string{"42"}); err != nil {
		t.Fatalf("updateWorkPackage error = %v, want nil", err)
	}

	if !strings.Contains(testingPrinter.Result, `"status": "In development"`) &&
		!strings.Contains(testingPrinter.Result, `"status":"In development"`) {
		t.Errorf("expected resolved status name in plan output, got: %q", testingPrinter.Result)
	}
	// only the work package fetch and the status lookup, no PATCH.
	if *requestCount != 2 {
		t.Errorf("request count = %d, want 2 (fetch + status list, no patch)", *requestCount)
	}
}
