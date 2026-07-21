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
