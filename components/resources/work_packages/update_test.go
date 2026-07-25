package work_packages_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

func TestUpdateValidatesTypeBeforeExecutingAction(t *testing.T) {
	mutationCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/work_packages/42":
			_, _ = response.Write([]byte(`{
				"id": 42,
				"lockVersion": 1,
				"_links": {
					"self": {"href": "/api/v3/work_packages/42"},
					"project": {"href": "/api/v3/projects/1"}
				},
				"_embedded": {
					"customActions": [{
						"name": "Claim",
						"_links": {
							"self": {"href": "/api/v3/custom_actions/1"},
							"executeImmediately": {
								"href": "/api/v3/custom_actions/1/execute",
								"method": "POST"
							}
						}
					}]
				}
			}`))
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1":
			_, _ = response.Write([]byte(`{"_links":{"types":{"href":"/api/v3/projects/1/types"}}}`))
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1/types":
			_, _ = response.Write([]byte(`{"_embedded":{"elements":[]}}`))
		case request.Method == http.MethodPost || request.Method == http.MethodPatch:
			mutationCount++
			_, _ = response.Write([]byte(`{}`))
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

	_, err = work_packages.Update("42", map[work_packages.UpdateOption]string{
		work_packages.UpdateCustomAction: "Claim",
		work_packages.UpdateType:         "Missing",
	})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("Update error = %v, want ErrHandled", err)
	}
	if mutationCount != 0 {
		t.Errorf("mutation count = %d, want 0", mutationCount)
	}
}

func TestUpdateValidatesAttachmentBeforeExecutingAction(t *testing.T) {
	mutationCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		if request.Method == http.MethodPost || request.Method == http.MethodPatch {
			mutationCount++
			_, _ = response.Write([]byte(`{}`))
			return
		}

		_, _ = response.Write([]byte(`{
			"id": 42,
			"lockVersion": 1,
			"_links": {
				"self": {"href": "/api/v3/work_packages/42"},
				"project": {"href": "/api/v3/projects/1"},
				"addAttachment": {
					"href": "/api/v3/work_packages/42/attachments",
					"method": "POST"
				}
			},
			"_embedded": {
				"customActions": [{
					"name": "Claim",
					"_links": {
						"self": {"href": "/api/v3/custom_actions/1"},
						"executeImmediately": {
							"href": "/api/v3/custom_actions/1/execute",
							"method": "POST"
						}
					}
				}]
			}
		}`))
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)
	printer.Init(&printer.TestingPrinter{})

	_, err = work_packages.Update("42", map[work_packages.UpdateOption]string{
		work_packages.UpdateCustomAction: "Claim",
		work_packages.UpdateAttachment:   "/file/does/not/exist",
	})
	if err == nil {
		t.Fatal("Update returned nil, want an attachment validation error")
	}
	if mutationCount != 0 {
		t.Errorf("mutation count = %d, want 0", mutationCount)
	}
}

func TestUpdateValidatesActionBeforePatching(t *testing.T) {
	mutationCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		if request.Method == http.MethodPatch || request.Method == http.MethodPost {
			mutationCount++
		}
		_, _ = response.Write([]byte(`{
			"id": 42,
			"lockVersion": 1,
			"_links": {
				"self": {"href": "/api/v3/work_packages/42"},
				"project": {"href": "/api/v3/projects/1"}
			},
			"_embedded": {"customActions": []}
		}`))
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)
	printer.Init(&printer.TestingPrinter{})

	_, err = work_packages.Update("42", map[work_packages.UpdateOption]string{
		work_packages.UpdateCustomAction: "Missing",
		work_packages.UpdateSubject:      "Changed",
	})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("Update error = %v, want ErrHandled", err)
	}
	if mutationCount != 0 {
		t.Errorf("mutation count = %d, want 0", mutationCount)
	}
}

func TestUpdateUsesRefetchedLockVersionAfterAction(t *testing.T) {
	fetchCount := 0
	mutationOrder := []string{}
	patchedLockVersion := 0
	patchedSubject := ""

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/work_packages/42":
			fetchCount++
			_, _ = fmt.Fprintf(response, `{
				"id": 42,
				"subject": "Before",
				"lockVersion": %d,
				"_links": {
					"self": {"href": "/api/v3/work_packages/42"},
					"project": {"href": "/api/v3/projects/1"}
				},
				"_embedded": {
					"customActions": [{
						"name": "Claim",
						"_links": {
							"self": {"href": "/api/v3/custom_actions/1"},
							"executeImmediately": {
								"href": "/api/v3/custom_actions/1/execute",
								"method": "POST"
							}
						}
					}]
				}
			}`, fetchCount)
		case request.Method == http.MethodPost && request.URL.Path == "/api/v3/custom_actions/1/execute":
			mutationOrder = append(mutationOrder, "action")
			_, _ = response.Write([]byte(`{}`))
		case request.Method == http.MethodPatch && request.URL.Path == "/api/v3/work_packages/42":
			mutationOrder = append(mutationOrder, "patch")
			var body struct {
				LockVersion int    `json:"lockVersion"`
				Subject     string `json:"subject"`
			}
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Errorf("decode PATCH body: %v", err)
			}
			patchedLockVersion = body.LockVersion
			patchedSubject = body.Subject
			_, _ = response.Write([]byte(`{}`))
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

	_, err = work_packages.Update("42", map[work_packages.UpdateOption]string{
		work_packages.UpdateCustomAction: "Claim",
		work_packages.UpdateSubject:      "Changed",
	})
	if err != nil {
		t.Fatalf("Update returned an error: %v", err)
	}
	if got := fmt.Sprint(mutationOrder); got != "[action patch]" {
		t.Errorf("mutation order = %s, want [action patch]", got)
	}
	if patchedLockVersion != 2 {
		t.Errorf("PATCH lockVersion = %d, want refetched version 2", patchedLockVersion)
	}
	if patchedSubject != "Changed" {
		t.Errorf("PATCH subject = %q, want Changed", patchedSubject)
	}
}

func TestUpdateRejectsInvalidTypeWithoutAction(t *testing.T) {
	mutationCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/work_packages/42":
			_, _ = response.Write([]byte(`{
				"id": 42,
				"lockVersion": 1,
				"_links": {
					"self": {"href": "/api/v3/work_packages/42"},
					"project": {"href": "/api/v3/projects/1"}
				},
				"_embedded": {"customActions": []}
			}`))
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1":
			_, _ = response.Write([]byte(`{"_links":{"types":{"href":"/api/v3/projects/1/types"}}}`))
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1/types":
			_, _ = response.Write([]byte(`{"_embedded":{"elements":[]}}`))
		case request.Method == http.MethodPost || request.Method == http.MethodPatch:
			mutationCount++
			_, _ = response.Write([]byte(`{}`))
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

	_, err = work_packages.Update("42", map[work_packages.UpdateOption]string{
		work_packages.UpdateSubject: "Changed",
		work_packages.UpdateType:    "Missing",
	})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("Update error = %v, want ErrHandled", err)
	}
	if mutationCount != 0 {
		t.Errorf("mutation count = %d, want 0", mutationCount)
	}
}

func TestUpdateRejectsAttachmentWithoutAddAttachmentLink(t *testing.T) {
	mutationCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		if request.Method == http.MethodPost || request.Method == http.MethodPatch {
			mutationCount++
			_, _ = response.Write([]byte(`{}`))
			return
		}

		_, _ = response.Write([]byte(`{
			"id": 42,
			"lockVersion": 1,
			"_links": {
				"self": {"href": "/api/v3/work_packages/42"},
				"project": {"href": "/api/v3/projects/1"}
			},
			"_embedded": {
				"customActions": [{
					"name": "Claim",
					"_links": {
						"self": {"href": "/api/v3/custom_actions/1"},
						"executeImmediately": {
							"href": "/api/v3/custom_actions/1/execute",
							"method": "POST"
						}
					}
				}]
			}
		}`))
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)
	printer.Init(&printer.TestingPrinter{})

	attachment := filepath.Join(t.TempDir(), "attachment.txt")
	if err := os.WriteFile(attachment, []byte("content"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err = work_packages.Update("42", map[work_packages.UpdateOption]string{
		work_packages.UpdateCustomAction: "Claim",
		work_packages.UpdateAttachment:   attachment,
	})
	if err == nil {
		t.Fatal("Update returned nil, want a missing addAttachment link error")
	}
	if !strings.Contains(err.Error(), "does not accept attachments") {
		t.Errorf("error should mention missing attachment capability, got: %v", err)
	}
	if mutationCount != 0 {
		t.Errorf("mutation count = %d, want 0", mutationCount)
	}
}

func TestUpdateRejectsFogStorageAttachments(t *testing.T) {
	mutationCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		if request.Method == http.MethodPost || request.Method == http.MethodPatch {
			mutationCount++
			_, _ = response.Write([]byte(`{}`))
			return
		}

		_, _ = response.Write([]byte(`{
			"id": 42,
			"lockVersion": 1,
			"_links": {
				"self": {"href": "/api/v3/work_packages/42"},
				"project": {"href": "/api/v3/projects/1"},
				"prepareAttachment": {
					"href": "/api/v3/work_packages/42/attachments/prepare",
					"method": "POST"
				}
			},
			"_embedded": {"customActions": []}
		}`))
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)
	printer.Init(&printer.TestingPrinter{})

	attachment := filepath.Join(t.TempDir(), "attachment.txt")
	if err := os.WriteFile(attachment, []byte("content"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err = work_packages.Update("42", map[work_packages.UpdateOption]string{
		work_packages.UpdateAttachment: attachment,
	})
	if err == nil {
		t.Fatal("Update returned nil, want a fog storage error")
	}
	if !strings.Contains(err.Error(), "fog storages") {
		t.Errorf("error should mention fog storages, got: %v", err)
	}
	if mutationCount != 0 {
		t.Errorf("mutation count = %d, want 0", mutationCount)
	}
}

func TestUpdatePatchIncludesStatus(t *testing.T) {
	var patchBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/work_packages/74416":
			_, _ = response.Write([]byte(`{
				"id": 74416,
				"subject": "Old subject",
				"lockVersion": 7,
				"_links": {
					"self": {"href": "/api/v3/work_packages/74416"},
					"project": {"href": "/api/v3/projects/1482"},
					"status": {"href": "/api/v3/statuses/1", "title": "New"}
				}
			}`))
		case request.Method == http.MethodPatch && request.URL.Path == "/api/v3/work_packages/74416":
			if err := json.NewDecoder(request.Body).Decode(&patchBody); err != nil {
				t.Fatal(err)
			}
			response.WriteHeader(http.StatusOK)
			_, _ = response.Write([]byte(`{}`))
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/statuses":
			_, _ = response.Write([]byte(`{
				"_embedded": {
					"elements": [
						{"id": 1, "name": "New"},
						{"id": 2, "name": "In development"}
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
	printer.Init(&printer.TestingPrinter{})

	_, err = work_packages.Update("74416", map[work_packages.UpdateOption]string{
		work_packages.UpdateStatus: "In development",
	})
	if err != nil {
		t.Fatal(err)
	}

	links, ok := patchBody["_links"].(map[string]any)
	if !ok {
		t.Fatalf("expected links object, got %#v", patchBody["_links"])
	}
	status, ok := links["status"].(map[string]any)
	if !ok {
		t.Fatalf("expected status link object, got %#v", links["status"])
	}
	if status["href"] != "/api/v3/statuses/2" {
		t.Fatalf("expected status href /api/v3/statuses/2, got %#v", status["href"])
	}
}

func TestUpdatePatchStatusIsCaseInsensitive(t *testing.T) {
	var patchBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/work_packages/74416":
			_, _ = response.Write([]byte(`{
				"id": 74416,
				"lockVersion": 7,
				"_links": {
					"self": {"href": "/api/v3/work_packages/74416"},
					"project": {"href": "/api/v3/projects/1482"}
				}
			}`))
		case request.Method == http.MethodPatch && request.URL.Path == "/api/v3/work_packages/74416":
			if err := json.NewDecoder(request.Body).Decode(&patchBody); err != nil {
				t.Fatal(err)
			}
			response.WriteHeader(http.StatusOK)
			_, _ = response.Write([]byte(`{}`))
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/statuses":
			_, _ = response.Write([]byte(`{
				"_embedded": {
					"elements": [
						{"id": 1, "name": "New"},
						{"id": 2, "name": "In development"}
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
	printer.Init(&printer.TestingPrinter{})

	_, err = work_packages.Update("74416", map[work_packages.UpdateOption]string{
		work_packages.UpdateStatus: "in development",
	})
	if err != nil {
		t.Fatal(err)
	}

	links, ok := patchBody["_links"].(map[string]any)
	if !ok {
		t.Fatalf("expected links object, got %#v", patchBody["_links"])
	}
	status, ok := links["status"].(map[string]any)
	if !ok {
		t.Fatalf("expected status link object, got %#v", links["status"])
	}
	if status["href"] != "/api/v3/statuses/2" {
		t.Fatalf("expected status href /api/v3/statuses/2, got %#v", status["href"])
	}
}

func TestUpdateReturnsErrorForUnknownStatus(t *testing.T) {
	mutationCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/work_packages/74416":
			_, _ = response.Write([]byte(`{
				"id": 74416,
				"lockVersion": 7,
				"_links": {
					"self": {"href": "/api/v3/work_packages/74416"},
					"project": {"href": "/api/v3/projects/1482"}
				}
			}`))
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/statuses":
			_, _ = response.Write([]byte(`{"_embedded":{"elements":[{"id": 1, "name": "New"}]}}`))
		case request.Method == http.MethodPost || request.Method == http.MethodPatch:
			mutationCount++
			_, _ = response.Write([]byte(`{}`))
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

	_, err = work_packages.Update("74416", map[work_packages.UpdateOption]string{
		work_packages.UpdateStatus: "Nonexistent",
	})
	if err == nil {
		t.Fatal("Update returned nil, want an unknown status error")
	}
	if !strings.Contains(err.Error(), `no status named "Nonexistent" found`) {
		t.Errorf("error should mention the unresolved status name, got: %v", err)
	}
	if mutationCount != 0 {
		t.Errorf("mutation count = %d, want 0", mutationCount)
	}
}

func TestDryRunUpdateIncludesLegacyFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		switch request.URL.Path {
		case "/api/v3/work_packages/74416":
			_, _ = response.Write([]byte(`{
				"id": 74416,
				"subject": "Old subject",
				"_links": {
					"self": {"href": "/api/v3/work_packages/74416"},
					"project": {"href": "/api/v3/projects/1482"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/7", "title": "Implementation"}
				}
			}`))
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
		case "/api/v3/statuses":
			_, _ = response.Write([]byte(`{
				"_embedded": {
					"elements": [
						{"id": 1, "name": "New"},
						{"id": 2, "name": "In development"}
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
	printer.Init(&printer.TestingPrinter{})

	cases := []struct {
		name    string
		options map[work_packages.UpdateOption]string
		field   string
		want    string
	}{
		{
			name:    "subject only",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateSubject: "Renamed"},
			field:   "subject",
			want:    "Renamed",
		},
		{
			name:    "type by name resolves against project types",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateType: "Feature"},
			field:   "type",
			want:    "Feature",
		},
		{
			name:    "assignee is echoed back",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateAssignee: "42"},
			field:   "assignee",
			want:    "42",
		},
		{
			name:    "status resolves against known statuses",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateStatus: "in development"},
			field:   "status",
			want:    "In development",
		},
		{
			name:    "action surfaces as preview",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateCustomAction: "Claim"},
			field:   "action",
			want:    "Claim",
		},
		{
			name:    "attach surfaces as preview",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateAttachment: "/tmp/f.txt"},
			field:   "attach",
			want:    "/tmp/f.txt",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := work_packages.DryRunUpdate("74416", tc.options)
			if err != nil {
				t.Fatalf("DryRunUpdate returned error: %v", err)
			}
			if !plan.Valid {
				t.Fatalf("expected valid plan, got %+v", plan)
			}
			if plan.WorkPackageID != "74416" {
				t.Errorf("WorkPackageID = %q, want 74416", plan.WorkPackageID)
			}

			marshalled, err := json.Marshal(plan)
			if err != nil {
				t.Fatal(err)
			}
			var unmarshalled map[string]any
			if err := json.Unmarshal(marshalled, &unmarshalled); err != nil {
				t.Fatal(err)
			}
			if unmarshalled[tc.field] != tc.want {
				t.Fatalf("expected %s=%q, got %#v", tc.field, tc.want, unmarshalled[tc.field])
			}
		})
	}
}

func TestDryRunUpdateReturnsErrorForMissingWorkPackage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusNotFound)
		_, _ = response.Write([]byte(`{"message":"not found"}`))
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)
	printer.Init(&printer.TestingPrinter{})

	_, err = work_packages.DryRunUpdate("999999", map[work_packages.UpdateOption]string{
		work_packages.UpdateSubject: "x",
	})
	if err == nil {
		t.Fatal("expected error for missing work package, got nil")
	}
}

func TestDryRunUpdateReturnsErrorForUnresolvedType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		switch request.URL.Path {
		case "/api/v3/work_packages/74416":
			_, _ = response.Write([]byte(`{
				"id": 74416,
				"_links": {
					"self": {"href": "/api/v3/work_packages/74416"},
					"project": {"href": "/api/v3/projects/1482"}
				}
			}`))
		case "/api/v3/projects/1482":
			_, _ = response.Write([]byte(`{
				"_links": {"types": {"href": "/api/v3/projects/1482/types/available"}}
			}`))
		case "/api/v3/projects/1482/types/available":
			_, _ = response.Write([]byte(`{"_embedded":{"elements":[]}}`))
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

	_, err = work_packages.DryRunUpdate("74416", map[work_packages.UpdateOption]string{
		work_packages.UpdateType: "Nonsense",
	})
	if err == nil {
		t.Fatal("expected error for unresolved type, got nil")
	}
	if !strings.Contains(err.Error(), "no unique available type from input") {
		t.Errorf("error should mention the unresolved type, got: %v", err)
	}
}

func TestDryRunUpdateReturnsErrorForUnresolvedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")

		switch request.URL.Path {
		case "/api/v3/work_packages/74416":
			_, _ = response.Write([]byte(`{
				"id": 74416,
				"_links": {
					"self": {"href": "/api/v3/work_packages/74416"},
					"project": {"href": "/api/v3/projects/1482"}
				}
			}`))
		case "/api/v3/statuses":
			_, _ = response.Write([]byte(`{"_embedded":{"elements":[]}}`))
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

	_, err = work_packages.DryRunUpdate("74416", map[work_packages.UpdateOption]string{
		work_packages.UpdateStatus: "In progress",
	})
	if err == nil {
		t.Fatal("expected error for unresolved status, got nil")
	}
	if !strings.Contains(err.Error(), `no status named "In progress" found`) {
		t.Errorf("error should mention the unresolved status, got: %v", err)
	}
}
