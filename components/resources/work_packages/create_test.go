package work_packages_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

func TestDryRunCreateWithParentInfersProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74316":
			_, _ = io.WriteString(w, `{
				"id": 74316,
				"subject": "Expand op CLI to support scripted work package workflows",
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
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Feature"},
					"assignee": {"href": null, "title": ""}
				}
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

	plan, err := work_packages.DryRunCreate(0, map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "Build reusable skill",
		work_packages.CreateParent:  "74316",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !plan.Valid || plan.ProjectID != 1482 {
		t.Fatalf("unexpected plan: %+v", plan)
	}

	if plan.ParentID == nil || *plan.ParentID != 74316 {
		t.Fatalf("expected parent id 74316, got %+v", plan.ParentID)
	}
}

func TestDryRunCreateWithTypeResolvesTypeName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74316":
			_, _ = io.WriteString(w, `{
				"id": 74316,
				"subject": "Expand op CLI to support scripted work package workflows",
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
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/6", "title": "Feature"},
					"assignee": {"href": null, "title": ""}
				}
			}`)
		case "/api/v3/projects/1482":
			_, _ = io.WriteString(w, `{
				"id": 1482,
				"identifier": "cli",
				"name": "CLI",
				"_links": {
					"types": {"href": "/api/v3/projects/1482/types/available"}
				}
			}`)
		case "/api/v3/projects/1482/types/available":
			_, _ = io.WriteString(w, `{
				"_embedded": {
					"elements": [
						{
							"id": 7,
							"name": "Implementation",
							"_links": {
								"self": {"href": "/api/v3/types/7"}
							}
						}
					]
				}
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

	plan, err := work_packages.DryRunCreate(0, map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "Build reusable skill",
		work_packages.CreateParent:  "74316",
		work_packages.CreateType:    "Implementation",
	})
	if err != nil {
		t.Fatal(err)
	}

	if plan.WorkPackage.Type != "Implementation" {
		t.Fatalf("expected type Implementation, got %+v", plan.WorkPackage)
	}
}
