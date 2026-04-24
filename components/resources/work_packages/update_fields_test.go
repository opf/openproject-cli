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
