package create

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

func TestCreateWorkPackagePrintsDryRunJSONWithParent(t *testing.T) {
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
	routes.Init(host)

	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	projectId = 0
	parentWorkPackageID = 74316
	shouldOpenWorkPackageInBrowser = false
	printCreatedWorkPackageAsJSON = true
	dryRunCreateWorkPackage = true
	typeFlag = ""
	descriptionFlag = ""
	descriptionFlagChanged = false

	createWorkPackage(nil, []string{"Build reusable skill"})

	expected := "{\"valid\":true,\"operation\":\"create\",\"project_id\":1482,\"parent_id\":74316,\"work_package\":{\"subject\":\"Build reusable skill\",\"type\":\"\",\"description\":\"\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestCreateWorkPackagePrintsDryRunJSONWithDescription(t *testing.T) {
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
	routes.Init(host)

	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	projectId = 0
	parentWorkPackageID = 74316
	shouldOpenWorkPackageInBrowser = false
	printCreatedWorkPackageAsJSON = true
	dryRunCreateWorkPackage = true
	typeFlag = "Implementation"
	descriptionFlag = "Body"
	descriptionFlagChanged = true

	createWorkPackage(nil, []string{"Add explicit work package description support"})

	expected := "{\"valid\":true,\"operation\":\"create\",\"project_id\":1482,\"parent_id\":74316,\"work_package\":{\"subject\":\"Add explicit work package description support\",\"type\":\"Implementation\",\"description\":\"Body\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestCreateWorkPackagePrintsJSONErrorForFlagConflict(t *testing.T) {
	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	projectId = 1482
	parentWorkPackageID = 0
	shouldOpenWorkPackageInBrowser = true
	printCreatedWorkPackageAsJSON = true
	dryRunCreateWorkPackage = false
	typeFlag = ""
	descriptionFlag = ""
	descriptionFlagChanged = false

	createWorkPackage(nil, []string{"Build reusable skill"})

	expected := "{\"error\":{\"code\":\"conflicting_arguments\",\"message\":\"cannot use --open together with --json\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestValidateCreateWorkPackageFlagsRejectsDryRunWithoutJSON(t *testing.T) {
	projectId = 1482
	parentWorkPackageID = 0
	shouldOpenWorkPackageInBrowser = false
	printCreatedWorkPackageAsJSON = false
	dryRunCreateWorkPackage = true
	typeFlag = ""
	descriptionFlag = ""
	descriptionFlagChanged = false

	err := validateCreateWorkPackageFlags()
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCreateWorkPackagePrintsJSONWithoutChildrenQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
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
		case "/api/v3/projects/1482/work_packages":
			_, _ = io.WriteString(w, `{
				"id": 74415,
				"subject": "Build reusable skill",
				"_links": {
					"self": {"href": "/api/v3/work_packages/74415"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/7", "title": "Implementation"},
					"assignee": {"href": null, "title": ""}
				},
				"_embedded": {
					"project": {
						"id": 1482,
						"identifier": "cli",
						"name": "CLI"
					}
				}
			}`)
		case "/api/v3/work_packages/74415":
			_, _ = io.WriteString(w, `{
				"id": 74415,
				"subject": "Build reusable skill",
				"_links": {
					"self": {"href": "/api/v3/work_packages/74415"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/7", "title": "Implementation"},
					"assignee": {"href": null, "title": ""}
				},
				"_embedded": {
					"project": {
						"id": 1482,
						"identifier": "cli",
						"name": "CLI"
					}
				}
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

	projectId = 1482
	parentWorkPackageID = 0
	shouldOpenWorkPackageInBrowser = false
	printCreatedWorkPackageAsJSON = true
	dryRunCreateWorkPackage = false
	typeFlag = "Implementation"

	createWorkPackage(nil, []string{"Build reusable skill"})

	expected := "{\"work_package\":{\"id\":74415,\"subject\":\"Build reusable skill\",\"type\":\"Implementation\",\"status\":\"new\",\"assignee\":\"\",\"description\":\"\",\"parent_id\":null,\"project\":{\"id\":1482,\"identifier\":\"cli\",\"name\":\"CLI\"},\"fields\":{},\"field_labels\":{}},\"children\":[]}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestCreateWorkPackagePrintsInvalidArgumentForValidationError(t *testing.T) {
	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	projectId = 0
	parentWorkPackageID = 0
	shouldOpenWorkPackageInBrowser = false
	printCreatedWorkPackageAsJSON = true
	dryRunCreateWorkPackage = false
	typeFlag = ""
	descriptionFlag = ""
	descriptionFlagChanged = false

	createWorkPackage(nil, []string{"Build reusable skill"})

	expected := "{\"error\":{\"code\":\"invalid_argument\",\"message\":\"either --project or --parent is required\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestCreateWorkPackagePrintsAPIErrorForDryRunParentLookupFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74316":
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
	routes.Init(host)

	activePrinter := &printer.TestingPrinter{}
	printer.Init(activePrinter)

	projectId = 0
	parentWorkPackageID = 74316
	shouldOpenWorkPackageInBrowser = false
	printCreatedWorkPackageAsJSON = true
	dryRunCreateWorkPackage = true
	typeFlag = ""
	descriptionFlag = ""
	descriptionFlagChanged = false

	createWorkPackage(nil, []string{"Build reusable skill"})

	expected := "{\"error\":{\"code\":\"api_error\",\"message\":\"{\\\"message\\\":\\\"boom\\\"}\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}

func TestCreateWorkPackagePrintsPostApplyInspectFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
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
		case "/api/v3/projects/1482/work_packages":
			_, _ = io.WriteString(w, `{
				"id": 74415,
				"subject": "Build reusable skill",
				"_links": {
					"self": {"href": "/api/v3/work_packages/74415"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/7", "title": "Implementation"},
					"assignee": {"href": null, "title": ""}
				}
			}`)
		case "/api/v3/work_packages/74415":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, `{"message":"inspect failed"}`)
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

	projectId = 1482
	parentWorkPackageID = 0
	shouldOpenWorkPackageInBrowser = false
	printCreatedWorkPackageAsJSON = true
	dryRunCreateWorkPackage = false
	typeFlag = "Implementation"
	descriptionFlag = ""
	descriptionFlagChanged = false

	createWorkPackage(nil, []string{"Build reusable skill"})

	expected := "{\"error\":{\"code\":\"post_apply_inspect_failed\",\"message\":\"{\\\"message\\\":\\\"inspect failed\\\"}\"}}\n"
	if activePrinter.Result != expected {
		t.Fatalf("expected %s, got %s", expected, activePrinter.Result)
	}
}
