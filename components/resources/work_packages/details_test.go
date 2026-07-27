package work_packages_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

func TestInspectWithChildren(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v3/work_packages/74316":
			_, _ = io.WriteString(w, `{
				"id": 74316,
				"subject": "Expand op CLI to support scripted work package workflows",
				"description": {"raw": "Body"},
				"customField130": 3,
				"_embedded": {
					"project": {
						"id": 1482,
						"identifier": "cli",
						"name": "CLI"
					}
				},
				"_links": {
					"self": {"href": "/api/v3/work_packages/74316"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Feature"},
					"assignee": {"href": null, "title": ""}
				}
			}`)
		case r.URL.Path == "/api/v3/work_packages/schemas/1482-6":
			_, _ = io.WriteString(w, `{
				"customField130": {"name": "Votes", "type": "Integer", "writable": true}
			}`)
		case r.URL.Path == "/api/v3/work_packages":
			if !strings.Contains(r.URL.RawQuery, "parent") {
				t.Fatalf("expected parent filter in query: %s", r.URL.RawQuery)
			}
			_, _ = io.WriteString(w, `{
				"_embedded": {
					"elements": [
						{
							"id": 74413,
							"subject": "Build a reusable SKILL.md based on OpenProject CLI",
							"_links": {
								"type": {"title": "Implementation"},
								"status": {"title": "new"},
								"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
								"parent": {"href": "/api/v3/work_packages/74316", "title": "Expand op CLI to support scripted work package workflows"}
							}
						}
					]
				},
				"_type": "Collection",
				"total": 1,
				"count": 1,
				"pageSize": -1,
				"offset": 1
			}`)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)

	payload, err := work_packages.InspectWithChildren("74316")
	if err != nil {
		t.Fatal(err)
	}

	if payload.WorkPackage.Project.Identifier != "cli" {
		t.Fatalf("expected project identifier cli, got %q", payload.WorkPackage.Project.Identifier)
	}

	if payload.WorkPackage.Fields["customField130"] != float64(3) {
		t.Fatalf("expected custom field value 3, got %#v", payload.WorkPackage.Fields["customField130"])
	}

	labels := payload.WorkPackage.FieldLabels["Votes"]
	if len(labels) != 1 || labels[0] != "customField130" {
		t.Fatalf("expected Votes label mapping, got %#v", payload.WorkPackage.FieldLabels)
	}

	if len(payload.Children) != 1 {
		t.Fatalf("expected one child, got %d", len(payload.Children))
	}

	if payload.Children[0].ID != 74413 {
		t.Fatalf("expected child 74413, got %+v", payload.Children[0])
	}
}

func TestInspectWithChildrenUsesNumericIdForParentFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v3/work_packages/PROJ-9":
			_, _ = io.WriteString(w, `{
				"id": 74316,
				"subject": "Expand op CLI to support scripted work package workflows",
				"description": {"raw": "Body"},
				"_embedded": {
					"project": {
						"id": 1482,
						"identifier": "cli",
						"name": "CLI"
					}
				},
				"_links": {
					"self": {"href": "/api/v3/work_packages/74316"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Feature"},
					"assignee": {"href": null, "title": ""}
				}
			}`)
		case r.URL.Path == "/api/v3/work_packages/schemas/1482-6":
			_, _ = io.WriteString(w, `{}`)
		case r.URL.Path == "/api/v3/work_packages":
			if !strings.Contains(r.URL.RawQuery, "74316") {
				t.Fatalf("expected parent filter to use numeric id 74316: %s", r.URL.RawQuery)
			}
			if strings.Contains(r.URL.RawQuery, "PROJ-9") {
				t.Fatalf("expected parent filter to not use semantic id PROJ-9: %s", r.URL.RawQuery)
			}
			_, _ = io.WriteString(w, `{
				"_embedded": {
					"elements": []
				},
				"_type": "Collection",
				"total": 0,
				"count": 0,
				"pageSize": -1,
				"offset": 1
			}`)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)

	payload, err := work_packages.InspectWithChildren("PROJ-9")
	if err != nil {
		t.Fatal(err)
	}

	if payload.WorkPackage.ID != 74316 {
		t.Fatalf("expected work package id 74316, got %d", payload.WorkPackage.ID)
	}
}

func TestInspectDoesNotQueryChildren(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74316":
			_, _ = io.WriteString(w, `{
				"id": 74316,
				"subject": "Expand op CLI to support scripted work package workflows",
				"description": {"raw": "Body"},
				"customField130": 3,
				"_embedded": {
					"project": {
						"id": 1482,
						"identifier": "cli",
						"name": "CLI"
					}
				},
				"_links": {
					"self": {"href": "/api/v3/work_packages/74316"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Feature"},
					"assignee": {"href": null, "title": ""}
				}
			}`)
		case "/api/v3/work_packages/schemas/1482-6":
			_, _ = io.WriteString(w, `{
				"customField130": {"name": "Votes", "type": "Integer", "writable": true}
			}`)
		case "/api/v3/work_packages":
			t.Fatalf("unexpected children query: %s", r.URL.Path)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)

	payload, err := work_packages.Inspect("74316")
	if err != nil {
		t.Fatal(err)
	}

	if len(payload.Children) != 0 {
		t.Fatalf("expected no children, got %+v", payload.Children)
	}
}

func TestInspectReturnsErrorForMissingWorkPackage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"message":"not found"}`)
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)

	_, err = work_packages.Inspect("999999")
	if err == nil {
		t.Fatal("expected error for missing work package, got nil")
	}
}

func TestInspectReturnsErrorWhenSchemaFetchFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74316":
			_, _ = io.WriteString(w, `{
				"id": 74316,
				"subject": "Any",
				"_links": {
					"self": {"href": "/api/v3/work_packages/74316"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Feature"},
					"assignee": {"href": null, "title": ""}
				}
			}`)
		case "/api/v3/work_packages/schemas/1482-6":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, `{"message":"boom"}`)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)

	_, err = work_packages.Inspect("74316")
	if err == nil {
		t.Fatal("expected error when schema endpoint fails, got nil")
	}
}

func TestInspectPreservesDuplicateFieldLabels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74316":
			_, _ = io.WriteString(w, `{
				"id": 74316,
				"subject": "Any",
				"_links": {
					"self": {"href": "/api/v3/work_packages/74316"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Feature"},
					"assignee": {"href": null, "title": ""}
				}
			}`)
		case "/api/v3/work_packages/schemas/1482-6":
			_, _ = io.WriteString(w, `{
				"customField17": {"name": "KPI", "type": "Integer", "writable": true},
				"customField22": {"name": "KPI", "type": "Integer", "writable": true}
			}`)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)

	payload, err := work_packages.Inspect("74316")
	if err != nil {
		t.Fatal(err)
	}

	labels := payload.WorkPackage.FieldLabels["KPI"]
	sort.Strings(labels)
	if len(labels) != 2 || labels[0] != "customField17" || labels[1] != "customField22" {
		t.Fatalf("expected duplicate KPI mappings, got %#v", payload.WorkPackage.FieldLabels)
	}
}

func TestInspectNormalizesFormattableCustomFieldsToRawStrings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74172":
			_, _ = io.WriteString(w, `{
				"id": 74172,
				"subject": "Epic",
				"customField401": {
					"format": "markdown",
					"html": "<p>Body</p>",
					"raw": "Body"
				},
				"_links": {
					"self": {"href": "/api/v3/work_packages/74172"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Epic"},
					"assignee": {"href": null, "title": ""}
				}
			}`)
		case "/api/v3/work_packages/schemas/1482-6":
			_, _ = io.WriteString(w, `{
				"customField401": {"name": "Acceptance criteria", "type": "Formattable", "writable": true}
			}`)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)

	payload, err := work_packages.Inspect("74172")
	if err != nil {
		t.Fatal(err)
	}

	if payload.WorkPackage.Fields["customField401"] != "Body" {
		t.Fatalf("expected Formattable field raw string, got %#v", payload.WorkPackage.Fields["customField401"])
	}
}

func TestInspectPassesThroughNonFormattableDictCustomField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74173":
			_, _ = io.WriteString(w, `{
				"id": 74173,
				"subject": "Epic",
				"customField402": {
					"id": 9,
					"name": "v1.0"
				},
				"_links": {
					"self": {"href": "/api/v3/work_packages/74173"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Epic"},
					"assignee": {"href": null, "title": ""}
				}
			}`)
		case "/api/v3/work_packages/schemas/1482-6":
			_, _ = io.WriteString(w, `{
				"customField402": {"name": "Release", "type": "Version", "writable": true}
			}`)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)

	payload, err := work_packages.Inspect("74173")
	if err != nil {
		t.Fatal(err)
	}

	value, ok := payload.WorkPackage.Fields["customField402"].(map[string]any)
	if !ok || value["name"] != "v1.0" {
		t.Fatalf("expected non-Formattable dict value untouched, got %#v", payload.WorkPackage.Fields["customField402"])
	}
}

func TestInspectPassesThroughFormattableFieldAlreadyAString(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74174":
			_, _ = io.WriteString(w, `{
				"id": 74174,
				"subject": "Epic",
				"customField401": "Already raw",
				"_links": {
					"self": {"href": "/api/v3/work_packages/74174"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Epic"},
					"assignee": {"href": null, "title": ""}
				}
			}`)
		case "/api/v3/work_packages/schemas/1482-6":
			_, _ = io.WriteString(w, `{
				"customField401": {"name": "Acceptance criteria", "type": "Formattable", "writable": true}
			}`)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)

	payload, err := work_packages.Inspect("74174")
	if err != nil {
		t.Fatal(err)
	}

	if payload.WorkPackage.Fields["customField401"] != "Already raw" {
		t.Fatalf("expected plain string Formattable field untouched, got %#v", payload.WorkPackage.Fields["customField401"])
	}
}

func TestInspectHandlesEmptyCustomFieldsMap(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74175":
			_, _ = io.WriteString(w, `{
				"id": 74175,
				"subject": "Epic",
				"_links": {
					"self": {"href": "/api/v3/work_packages/74175"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Epic"},
					"assignee": {"href": null, "title": ""}
				}
			}`)
		case "/api/v3/work_packages/schemas/1482-6":
			_, _ = io.WriteString(w, `{
				"customField401": {"name": "Acceptance criteria", "type": "Formattable", "writable": true}
			}`)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)

	payload, err := work_packages.Inspect("74175")
	if err != nil {
		t.Fatal(err)
	}

	if len(payload.WorkPackage.Fields) != 0 {
		t.Fatalf("expected empty fields map, got %#v", payload.WorkPackage.Fields)
	}
}
