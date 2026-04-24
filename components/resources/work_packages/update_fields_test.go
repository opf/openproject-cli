package work_packages_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

func TestDryRunUpdateFieldsResolvesLabels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74316":
			_, _ = io.WriteString(w, `{
				"id": 74316,
				"subject": "Expand op CLI to support scripted work package workflows",
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

	plan, err := work_packages.DryRunUpdateFields(74316, []string{"Votes=3"})
	if err != nil {
		t.Fatal(err)
	}

	field := plan.ResolvedFields["Votes"]
	if field.APIField != "customField130" || field.Value != int64(3) {
		t.Fatalf("unexpected plan: %+v", plan)
	}
}

func TestUpdateFieldsPatchesFormattableCustomFieldsAsLongText(t *testing.T) {
	var patchBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74172":
			switch r.Method {
			case http.MethodGet:
				_, _ = io.WriteString(w, `{
					"id": 74172,
					"subject": "Epic",
					"lockVersion": 5,
					"_links": {
						"self": {"href": "/api/v3/work_packages/74172"},
						"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
						"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
						"status": {"href": "/api/v3/statuses/1", "title": "new"},
						"type": {"href": "/api/v3/types/6", "title": "Feature"},
						"assignee": {"href": null, "title": ""}
					}
				}`)
			case http.MethodPatch:
				if err := json.NewDecoder(r.Body).Decode(&patchBody); err != nil {
					t.Fatal(err)
				}
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, `{}`)
			default:
				t.Fatalf("unexpected method %s", r.Method)
			}
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

	if err := work_packages.UpdateFields(74172, []string{"Acceptance criteria=Native HTML5 DnD was rejected.\nDragula was rejected."}); err != nil {
		t.Fatal(err)
	}

	value, ok := patchBody["customField401"].(map[string]any)
	if !ok {
		t.Fatalf("expected long-text patch object, got %#v", patchBody["customField401"])
	}
	if value["raw"] != "Native HTML5 DnD was rejected.\nDragula was rejected." {
		t.Fatalf("expected raw long text, got %#v", value["raw"])
	}
}
