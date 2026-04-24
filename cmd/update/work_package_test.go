package update

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/routes"
)

func TestUpdateWorkPackagePrintsDryRunJSON(t *testing.T) {
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
	routes.Init(host)

	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	actionFlag = ""
	assigneeFlag = 0
	attachFlag = ""
	subjectFlag = ""
	typeFlag = ""
	setFlags = []string{"Votes=3"}
	printUpdatedWorkPackageAsJSON = true
	dryRunUpdateWorkPackage = true

	updateWorkPackage(nil, []string{"74316"})

	expected := "{\"valid\":true,\"operation\":\"update\",\"work_package_id\":74316,\"resolved_fields\":{\"Votes\":{\"api_field\":\"customField130\",\"value\":3}}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestValidateUpdateWorkPackageFlagsRejectsDryRunWithoutJSON(t *testing.T) {
	actionFlag = ""
	assigneeFlag = 0
	attachFlag = ""
	subjectFlag = ""
	typeFlag = ""
	setFlags = []string{"Votes=3"}
	printUpdatedWorkPackageAsJSON = false
	dryRunUpdateWorkPackage = true

	err := validateUpdateWorkPackageFlags()
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateUpdateWorkPackageFlagsRejectsMixedSetAndLegacyFlags(t *testing.T) {
	actionFlag = ""
	assigneeFlag = 0
	attachFlag = ""
	subjectFlag = "new subject"
	typeFlag = ""
	setFlags = []string{"Votes=3"}
	printUpdatedWorkPackageAsJSON = true
	dryRunUpdateWorkPackage = false

	err := validateUpdateWorkPackageFlags()
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestUpdateWorkPackagePrintsJSONWithoutChildrenQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74316":
			switch r.Method {
			case http.MethodGet:
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
						"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
						"status": {"href": "/api/v3/statuses/1", "title": "new"},
						"type": {"href": "/api/v3/types/6", "title": "Feature"},
						"assignee": {"href": null, "title": ""}
					}
				}`)
			case http.MethodPatch:
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, `{}`)
			default:
				t.Fatalf("unexpected method %s", r.Method)
			}
		case "/api/v3/work_packages/schemas/1482-6":
			_, _ = io.WriteString(w, `{
				"customField130": {"name": "Votes", "type": "Integer", "writable": true}
			}`)
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

	actionFlag = ""
	assigneeFlag = 0
	attachFlag = ""
	subjectFlag = ""
	typeFlag = ""
	setFlags = []string{"Votes=3"}
	printUpdatedWorkPackageAsJSON = true
	dryRunUpdateWorkPackage = false

	updateWorkPackage(nil, []string{"74316"})

	expected := "{\"work_package\":{\"id\":74316,\"subject\":\"Expand op CLI to support scripted work package workflows\",\"type\":\"Feature\",\"status\":\"new\",\"assignee\":\"\",\"description\":\"\",\"parent_id\":null,\"project\":{\"id\":1482,\"identifier\":\"cli\",\"name\":\"CLI\"},\"fields\":{},\"field_labels\":{\"Votes\":[\"customField130\"]}},\"children\":[]}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestUpdateWorkPackagePrintsAmbiguousFieldJSONError(t *testing.T) {
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
	routes.Init(host)

	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	actionFlag = ""
	assigneeFlag = 0
	attachFlag = ""
	subjectFlag = ""
	typeFlag = ""
	setFlags = []string{"KPI=3"}
	printUpdatedWorkPackageAsJSON = true
	dryRunUpdateWorkPackage = true

	updateWorkPackage(nil, []string{"74316"})

	expected := "{\"error\":{\"code\":\"ambiguous_field\",\"message\":\"ambiguous field: \\\"KPI\\\"\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestUpdateWorkPackagePrintsDuplicateFieldJSONError(t *testing.T) {
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
	routes.Init(host)

	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	actionFlag = ""
	assigneeFlag = 0
	attachFlag = ""
	subjectFlag = ""
	typeFlag = ""
	setFlags = []string{"Votes=3", "customField130=4"}
	printUpdatedWorkPackageAsJSON = true
	dryRunUpdateWorkPackage = true

	updateWorkPackage(nil, []string{"74316"})

	expected := "{\"error\":{\"code\":\"duplicate_field\",\"message\":\"duplicate field: \\\"customField130\\\"\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestUpdateWorkPackagePrintsAPIErrorForSetDryRunFetchFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, `{"message":"boom"}`)
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

	actionFlag = ""
	assigneeFlag = 0
	attachFlag = ""
	subjectFlag = ""
	typeFlag = ""
	setFlags = []string{"Votes=3"}
	printUpdatedWorkPackageAsJSON = true
	dryRunUpdateWorkPackage = true

	updateWorkPackage(nil, []string{"74316"})

	expected := "{\"error\":{\"code\":\"api_error\",\"message\":\"{\\\"message\\\":\\\"boom\\\"}\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestUpdateWorkPackagePrintsAPIErrorWhenPatchFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74416":
			switch r.Method {
			case http.MethodGet:
				_, _ = io.WriteString(w, `{
					"id": 74416,
					"subject": "Old subject",
					"lockVersion": 7,
					"_links": {
						"self": {"href": "/api/v3/work_packages/74416"},
						"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
						"status": {"href": "/api/v3/statuses/1", "title": "new"},
						"type": {"href": "/api/v3/types/7", "title": "Implementation"},
						"assignee": {"href": null, "title": ""}
					}
				}`)
			case http.MethodPatch:
				w.WriteHeader(http.StatusConflict)
				_, _ = io.WriteString(w, `{"message":"conflict"}`)
			default:
				t.Fatalf("unexpected method %s", r.Method)
			}
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

	actionFlag = ""
	assigneeFlag = 0
	attachFlag = ""
	subjectFlag = "New subject"
	typeFlag = ""
	setFlags = nil
	printUpdatedWorkPackageAsJSON = true
	dryRunUpdateWorkPackage = false

	updateWorkPackage(nil, []string{"74416"})

	expected := "{\"error\":{\"code\":\"api_error\",\"message\":\"{\\\"message\\\":\\\"conflict\\\"}\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestUpdateWorkPackageWithoutFlagsPrintsCurrentWorkPackage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74316":
			_, _ = io.WriteString(w, `{
				"id": 74316,
				"subject": "Expand op CLI to support scripted work package workflows",
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
	routes.Init(host)

	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	actionFlag = ""
	assigneeFlag = 0
	attachFlag = ""
	subjectFlag = ""
	typeFlag = ""
	setFlags = nil
	printUpdatedWorkPackageAsJSON = false
	dryRunUpdateWorkPackage = false

	updateWorkPackage(nil, []string{"74316"})

	if activePrinter.Result == "" {
		t.Fatal("expected work package output, got empty result")
	}
}

func TestUpdateWorkPackagePrintsHumanProgressForNonJSONUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74416":
			switch r.Method {
			case http.MethodGet:
				_, _ = io.WriteString(w, `{
					"id": 74416,
					"subject": "Old subject",
					"lockVersion": 7,
					"_links": {
						"self": {"href": "/api/v3/work_packages/74416"},
						"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
						"status": {"href": "/api/v3/statuses/1", "title": "new"},
						"type": {"href": "/api/v3/types/7", "title": "Implementation"},
						"assignee": {"href": null, "title": ""}
					}
				}`)
			case http.MethodPatch:
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, `{}`)
			default:
				t.Fatalf("unexpected method %s", r.Method)
			}
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

	actionFlag = ""
	assigneeFlag = 0
	attachFlag = ""
	subjectFlag = "New subject"
	typeFlag = ""
	setFlags = nil
	printUpdatedWorkPackageAsJSON = false
	dryRunUpdateWorkPackage = false

	updateWorkPackage(nil, []string{"74416"})

	if !strings.Contains(activePrinter.Result, "Updating work package ...") {
		t.Fatalf("expected human progress output, got %q", activePrinter.Result)
	}
}
