package work_packages_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

func TestUpdatePatchIncludesSubjectAndDescription(t *testing.T) {
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
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}

				description, ok := body["description"].(map[string]any)
				if !ok {
					t.Fatalf("expected description object, got %#v", body["description"])
				}

				if body["subject"] != "New subject" {
					t.Fatalf("expected subject New subject, got %#v", body["subject"])
				}

				if description["raw"] != "Body" {
					t.Fatalf("expected description raw Body, got %#v", description["raw"])
				}

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
	printer.Init(&printer.TestingPrinter{})

	_, err = work_packages.Update(74416, map[work_packages.UpdateOption]string{
		work_packages.UpdateSubject:     "New subject",
		work_packages.UpdateDescription: "Body",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDryRunUpdateIncludesAllLegacyFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74416":
			_, _ = io.WriteString(w, `{
				"id": 74416,
				"subject": "Old subject",
				"_links": {
					"self": {"href": "/api/v3/work_packages/74416"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"},
					"status": {"href": "/api/v3/statuses/1", "title": "new"},
					"type": {"href": "/api/v3/types/7", "title": "Implementation"},
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
							"_links": {"self": {"href": "/api/v3/types/7"}}
						},
						{
							"id": 6,
							"name": "Feature",
							"_links": {"self": {"href": "/api/v3/types/6"}}
						}
					]
				}
			}`)
		case "/api/v3/statuses":
			_, _ = io.WriteString(w, `{
				"_embedded": {
					"elements": [
						{"id": 1, "name": "New"},
						{"id": 2, "name": "In development"}
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
	printer.Init(&printer.TestingPrinter{})

	cases := []struct {
		name    string
		options map[work_packages.UpdateOption]string
		check   func(t *testing.T, plan any)
	}{
		{
			name:    "subject only",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateSubject: "Renamed"},
			check: func(t *testing.T, plan any) {
				assertPlanField(t, plan, "subject", "Renamed")
			},
		},
		{
			name:    "type by name resolves against project types",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateType: "Feature"},
			check: func(t *testing.T, plan any) {
				assertPlanField(t, plan, "type", "Feature")
			},
		},
		{
			name:    "assignee is echoed back",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateAssignee: "42"},
			check: func(t *testing.T, plan any) {
				assertPlanField(t, plan, "assignee", "42")
			},
		},
		{
			name:    "description round-trips",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateDescription: "New body"},
			check: func(t *testing.T, plan any) {
				assertPlanField(t, plan, "description", "New body")
			},
		},
		{
			name:    "status resolves against known statuses",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateStatus: "in development"},
			check: func(t *testing.T, plan any) {
				assertPlanField(t, plan, "status", "In development")
			},
		},
		{
			name:    "action surfaces as preview",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateCustomAction: "Claim"},
			check: func(t *testing.T, plan any) {
				assertPlanField(t, plan, "action", "Claim")
			},
		},
		{
			name:    "attach surfaces as preview",
			options: map[work_packages.UpdateOption]string{work_packages.UpdateAttachment: "/tmp/f.txt"},
			check: func(t *testing.T, plan any) {
				assertPlanField(t, plan, "attach", "/tmp/f.txt")
			},
		},
		{
			name: "combined legacy fields all appear",
			options: map[work_packages.UpdateOption]string{
				work_packages.UpdateSubject:     "S",
				work_packages.UpdateDescription: "D",
				work_packages.UpdateAssignee:    "7",
			},
			check: func(t *testing.T, plan any) {
				assertPlanField(t, plan, "subject", "S")
				assertPlanField(t, plan, "description", "D")
				assertPlanField(t, plan, "assignee", "7")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := work_packages.DryRunUpdate(74416, tc.options)
			if err != nil {
				t.Fatalf("DryRunUpdate returned error: %v", err)
			}

			marshalled, err := json.Marshal(plan)
			if err != nil {
				t.Fatal(err)
			}

			var unmarshalled map[string]any
			if err := json.Unmarshal(marshalled, &unmarshalled); err != nil {
				t.Fatal(err)
			}

			if unmarshalled["valid"] != true {
				t.Fatalf("expected valid:true, got %#v", unmarshalled["valid"])
			}

			tc.check(t, unmarshalled)
		})
	}
}

func TestDryRunUpdateReturnsErrorForMissingWorkPackage(t *testing.T) {
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
	printer.Init(&printer.TestingPrinter{})

	_, err = work_packages.DryRunUpdate(999999, map[work_packages.UpdateOption]string{
		work_packages.UpdateSubject: "x",
	})
	if err == nil {
		t.Fatal("expected error for missing work package, got nil")
	}
}

func TestDryRunUpdateReturnsErrorForUnresolvedType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74416":
			_, _ = io.WriteString(w, `{
				"id": 74416,
				"_links": {
					"self": {"href": "/api/v3/work_packages/74416"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"}
				}
			}`)
		case "/api/v3/projects/1482":
			_, _ = io.WriteString(w, `{
				"id": 1482,
				"_links": {"types": {"href": "/api/v3/projects/1482/types/available"}}
			}`)
		case "/api/v3/projects/1482/types/available":
			_, _ = io.WriteString(w, `{"_embedded":{"elements":[]}}`)
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
	printer.Init(&printer.TestingPrinter{})

	_, err = work_packages.DryRunUpdate(74416, map[work_packages.UpdateOption]string{
		work_packages.UpdateType: "Nonsense",
	})
	if err == nil {
		t.Fatal("expected error for unresolved type, got nil")
	}
}

func TestDryRunUpdateReturnsErrorForUnresolvedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/74416":
			_, _ = io.WriteString(w, `{
				"id": 74416,
				"_links": {
					"self": {"href": "/api/v3/work_packages/74416"},
					"project": {"href": "/api/v3/projects/1482", "title": "CLI"}
				}
			}`)
		case "/api/v3/statuses":
			_, _ = io.WriteString(w, `{"_embedded":{"elements":[]}}`)
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
	printer.Init(&printer.TestingPrinter{})

	_, err = work_packages.DryRunUpdate(74416, map[work_packages.UpdateOption]string{
		work_packages.UpdateStatus: "In progress",
	})
	if err == nil {
		t.Fatal("expected error for unresolved status, got nil")
	}
}

func TestUpdatePatchIncludesStatus(t *testing.T) {
	var patchBody map[string]any

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
						"status": {"href": "/api/v3/statuses/1", "title": "New"},
						"type": {"href": "/api/v3/types/7", "title": "Implementation"},
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
		case "/api/v3/statuses":
			_, _ = io.WriteString(w, `{
				"_embedded": {
					"elements": [
						{"id": 1, "name": "New"},
						{"id": 2, "name": "In development"}
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
	printer.Init(&printer.TestingPrinter{})

	_, err = work_packages.Update(74416, map[work_packages.UpdateOption]string{
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

func assertPlanField(t *testing.T, plan any, field string, expected string) {
	t.Helper()
	m, ok := plan.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", plan)
	}
	if m[field] != expected {
		t.Fatalf("expected %s=%q, got %#v", field, expected, m[field])
	}
}
