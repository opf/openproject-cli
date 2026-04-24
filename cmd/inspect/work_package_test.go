package inspect

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/routes"
)

func TestInspectWorkPackagePrintsJSONWithChildren(t *testing.T) {
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
			_, _ = io.WriteString(w, `{"customField130": {"name": "Votes", "type": "Integer", "writable": true}}`)
		case "/api/v3/work_packages":
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
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)
	routes.Init(host)

	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	shouldOpenWorkPackageInBrowser = false
	listAvailableTypes = false
	includeChildrenInJson = true
	printWorkPackageAsJSON = true

	inspectWorkPackage(nil, []string{"74316"})

	expected := "{\"work_package\":{\"id\":74316,\"subject\":\"Expand op CLI to support scripted work package workflows\",\"type\":\"Feature\",\"status\":\"new\",\"assignee\":\"\",\"description\":\"Body\",\"parent_id\":null,\"project\":{\"id\":1482,\"identifier\":\"cli\",\"name\":\"CLI\"},\"fields\":{\"customField130\":3},\"field_labels\":{\"Votes\":[\"customField130\"]}},\"children\":[{\"id\":74413,\"subject\":\"Build a reusable SKILL.md based on OpenProject CLI\",\"type\":\"Implementation\",\"status\":\"new\",\"parent_id\":74316}]}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestInspectWorkPackagePrintsJSONWithoutChildrenQuery(t *testing.T) {
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
			_, _ = io.WriteString(w, `{"customField130": {"name": "Votes", "type": "Integer", "writable": true}}`)
		case "/api/v3/work_packages":
			t.Fatalf("unexpected children query: %s", r.URL.Path)
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
	routes.Init(host)

	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	shouldOpenWorkPackageInBrowser = false
	listAvailableTypes = false
	includeChildrenInJson = false
	printWorkPackageAsJSON = true

	inspectWorkPackage(nil, []string{"74316"})

	expected := "{\"work_package\":{\"id\":74316,\"subject\":\"Expand op CLI to support scripted work package workflows\",\"type\":\"Feature\",\"status\":\"new\",\"assignee\":\"\",\"description\":\"Body\",\"parent_id\":null,\"project\":{\"id\":1482,\"identifier\":\"cli\",\"name\":\"CLI\"},\"fields\":{\"customField130\":3},\"field_labels\":{\"Votes\":[\"customField130\"]}},\"children\":[]}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestValidateInspectWorkPackageFlagsRejectsOpenAndJSON(t *testing.T) {
	shouldOpenWorkPackageInBrowser = true
	listAvailableTypes = false
	includeChildrenInJson = false
	printWorkPackageAsJSON = true

	err := validateInspectWorkPackageFlags()
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateInspectWorkPackageFlagsRejectsTypesAndJSON(t *testing.T) {
	shouldOpenWorkPackageInBrowser = false
	listAvailableTypes = true
	includeChildrenInJson = false
	printWorkPackageAsJSON = true

	err := validateInspectWorkPackageFlags()
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestInspectWorkPackagePrintsJSONErrorForFlagConflict(t *testing.T) {
	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	shouldOpenWorkPackageInBrowser = true
	listAvailableTypes = false
	includeChildrenInJson = false
	printWorkPackageAsJSON = true

	inspectWorkPackage(nil, []string{"74316"})

	expected := "{\"error\":{\"code\":\"conflicting_arguments\",\"message\":\"cannot use --open together with --json\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestInspectWorkPackagePrintsJSONErrorForChildrenWithoutJSON(t *testing.T) {
	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	shouldOpenWorkPackageInBrowser = false
	listAvailableTypes = false
	includeChildrenInJson = true
	printWorkPackageAsJSON = false

	inspectWorkPackage(nil, []string{"42"})

	expected := "\u001b[31m[ERROR]\u001b[0m cannot use --children without --json\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %q, got %q", expected, activePrinter.Result)
	}
}

func TestInspectWorkPackagePrintsJSONErrorForTypesAndJSON(t *testing.T) {
	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	shouldOpenWorkPackageInBrowser = false
	listAvailableTypes = true
	includeChildrenInJson = false
	printWorkPackageAsJSON = true

	inspectWorkPackage(nil, []string{"74316"})

	expected := "{\"error\":{\"code\":\"conflicting_arguments\",\"message\":\"cannot use --types together with --json or --children\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}
