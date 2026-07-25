package work_packages_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

func TestCreateRejectsInvalidTypeBeforePosting(t *testing.T) {
	postCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1":
			_, _ = response.Write([]byte(`{"_links":{"types":{"href":"/api/v3/projects/1/types"}}}`))
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1/types":
			_, _ = response.Write([]byte(`{"_embedded":{"elements":[]}}`))
		case request.Method == http.MethodPost:
			postCount++
			_, _ = response.Write([]byte(`{"id":42}`))
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

	_, err = work_packages.Create("1", map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "Subject",
		work_packages.CreateType:    "Missing",
	})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("Create error = %v, want ErrHandled", err)
	}
	if postCount != 0 {
		t.Errorf("POST count = %d, want 0", postCount)
	}
}

// On an unresolved --type, the "available types" hint is guidance and must go
// to stderr, never to stdout, so --format json output stays clean.
func TestCreateInvalidTypeHintGoesToStderr(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1":
			_, _ = response.Write([]byte(`{"_links":{"types":{"href":"/api/v3/projects/1/types"}}}`))
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/projects/1/types":
			_, _ = response.Write([]byte(`{"_embedded":{"elements":[{"id":7,"name":"Task"}]}}`))
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
	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)

	_, err = work_packages.Create("1", map[work_packages.CreateOption]string{
		work_packages.CreateSubject: "Subject",
		work_packages.CreateType:    "Missing",
	})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("Create error = %v, want ErrHandled", err)
	}
	if testingPrinter.Result != "" {
		t.Errorf("stdout = %q, want empty on failed create", testingPrinter.Result)
	}
	if !strings.Contains(testingPrinter.ErrResult, "Task") {
		t.Errorf("stderr = %q, want it to list available type %q", testingPrinter.ErrResult, "Task")
	}
}
