package work_packages_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

func TestCreateRejectsInvalidTypeBeforePosting(t *testing.T) {
	postCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1":
			_, _ = response.Write([]byte(`{"_links":{"types":{"href":"/api/v3/projects/1/types"}}}`))
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1/types":
			_, _ = response.Write([]byte(`{"_embedded":{"elements":[]}}`))
		case request.Method == http.MethodPost:
			postCount++
			_, _ = response.Write([]byte(`{"id":42}`))
		default:
			http.Error(response, "unexpected request", http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)
	printer.Init(&printer.TestingPrinter{})

	_, err = work_packages.Create("1", map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "Subject",
		work_packages.CreateType:    "Missing",
	})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("Create error = %v, want ErrHandled", err)
	}
	if postCount != 0 {
		t.Errorf("POST count = %d, want 0", postCount)
	}
}

// On an unresolved --type, the "available types" hint is guidance and must go
// to stderr, never to stdout, so --format json output stays clean.
func TestCreateInvalidTypeHintGoesToStderr(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1":
			_, _ = response.Write([]byte(`{"_links":{"types":{"href":"/api/v3/projects/1/types"}}}`))
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1/types":
			_, _ = response.Write([]byte(`{"_embedded":{"elements":[{"id":7,"name":"Task"}]}}`))
		default:
			http.Error(response, "unexpected request", http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)
	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)

	_, err = work_packages.Create("1", map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "Subject",
		work_packages.CreateType:    "Missing",
	})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("Create error = %v, want ErrHandled", err)
	}
	if testingPrinter.Result != "" {
		t.Errorf("stdout = %q, want empty on failed create", testingPrinter.Result)
	}
	if !strings.Contains(testingPrinter.ErrResult, "Task") {
		t.Errorf("stderr = %q, want it to list available type %q", testingPrinter.ErrResult, "Task")
	}
}

func TestDryRunCreateIncludesParentAndDescription(t *testing.T) {
	plan, err := work_packages.DryRunCreate("1482", map[work_packages.CreateOption]string{
		work_packages.CreateSubject:     "Build reusable skill",
		work_packages.CreateParent:      "74316",
		work_packages.CreateDescription: "Body",
	})
	if err != nil {
		t.Fatalf("DryRunCreate returned error: %v", err)
	}

	if !plan.Valid || plan.Operation != "create" {
		t.Fatalf("unexpected plan: %+v", plan)
	}
	if plan.ProjectID != "1482" {
		t.Errorf("ProjectID = %q, want 1482", plan.ProjectID)
	}
	if plan.ParentID == nil || *plan.ParentID != 74316 {
		t.Fatalf("expected parent id 74316, got %+v", plan.ParentID)
	}
	if plan.WorkPackage.Subject != "Build reusable skill" {
		t.Errorf("Subject = %q, want %q", plan.WorkPackage.Subject, "Build reusable skill")
	}
	if plan.WorkPackage.Description != "Body" {
		t.Errorf("Description = %q, want %q", plan.WorkPackage.Description, "Body")
	}
}

func TestDryRunCreateWithTypeResolvesTypeName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/api/v3/projects/1482":
			_, _ = response.Write([]byte(`{
				"id": 1482,
				"_links": {"types": {"href": "/api/v3/projects/1482/types/available"}}
			}`))
		case "/api/v3/projects/1482/types/available":
			_, _ = response.Write([]byte(`{
				"_embedded": {
					"elements": [
						{"id": 7, "name": "Implementation", "_links": {"self": {"href": "/api/v3/types/7"}}},
						{"id": 6, "name": "Feature", "_links": {"self": {"href": "/api/v3/types/6"}}}
					]
				}
			}`))
		default:
			http.Error(response, "unexpected request", http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)

	plan, err := work_packages.DryRunCreate("1482", map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "Build reusable skill",
		work_packages.CreateType:    "Feature",
	})
	if err != nil {
		t.Fatalf("DryRunCreate returned error: %v", err)
	}
	if plan.WorkPackage.Type != "Feature" {
		t.Errorf("Type = %q, want %q", plan.WorkPackage.Type, "Feature")
	}
}

func TestDryRunCreateReturnsErrorForUnresolvedType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/api/v3/projects/1482":
			_, _ = response.Write([]byte(`{
				"id": 1482,
				"_links": {"types": {"href": "/api/v3/projects/1482/types/available"}}
			}`))
		case "/api/v3/projects/1482/types/available":
			_, _ = response.Write([]byte(`{
				"_embedded": {
					"elements": [
						{"id": 7, "name": "Implementation", "_links": {"self": {"href": "/api/v3/types/7"}}}
					]
				}
			}`))
		default:
			http.Error(response, "unexpected request", http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)

	_, err = work_packages.DryRunCreate("1482", map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "Build reusable skill",
		work_packages.CreateType:    "Missing",
	})
	if err == nil {
		t.Fatal("expected error for unresolved type, got nil")
	}
}

func TestDryRunCreateWithoutParentLeavesParentIDNil(t *testing.T) {
	plan, err := work_packages.DryRunCreate("1482", map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "No parent",
	})
	if err != nil {
		t.Fatalf("DryRunCreate returned error: %v", err)
	}
	if plan.ParentID != nil {
		t.Errorf("ParentID = %+v, want nil", plan.ParentID)
	}
}

func TestDryRunCreateReturnsErrorForInvalidParent(t *testing.T) {
	_, err := work_packages.DryRunCreate("1482", map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "Subject",
		work_packages.CreateParent:  "not-a-number",
	})
	if err == nil {
		t.Fatal("expected error for invalid parent id, got nil")
	}
}

func TestCreatePatchIncludesParentLink(t *testing.T) {
	var postBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodPost && request.URL.Path == "/api/v3/projects/1482/work_packages":
			if err := json.NewDecoder(request.Body).Decode(&postBody); err != nil {
				t.Fatal(err)
			}
			response.WriteHeader(http.StatusCreated)
			_, _ = response.Write([]byte(`{"id":99999,"subject":"","_links":{"self":{"href":"/api/v3/work_packages/99999"}}}`))
		default:
			http.Error(response, "unexpected request", http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)

	_, err = work_packages.Create("1482", map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "Subject",
		work_packages.CreateParent:  "74316",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	links, ok := postBody["_links"].(map[string]any)
	if !ok {
		t.Fatalf("expected _links object in POST body, got %#v", postBody["_links"])
	}
	parent, ok := links["parent"].(map[string]any)
	if !ok {
		t.Fatalf("expected parent link in POST body, got %#v", links["parent"])
	}
	if parent["href"] != "/api/v3/work_packages/74316" {
		t.Errorf("parent href = %v, want %q", parent["href"], "/api/v3/work_packages/74316")
	}
}

func TestCreateRejectsInvalidParentBeforePosting(t *testing.T) {
	postCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			postCount++
		}
		response.WriteHeader(http.StatusCreated)
		_, _ = response.Write([]byte(`{"id":1}`))
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)

	_, err = work_packages.Create("1482", map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "Subject",
		work_packages.CreateParent:  "not-a-number",
	})
	if err == nil {
		t.Fatal("expected error for invalid parent id, got nil")
	}
	if postCount != 0 {
		t.Errorf("POST count = %d, want 0", postCount)
	}
}
