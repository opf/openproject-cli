package time_entries

import (
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
)

// An unresolved --activity is already reported inside the resource layer, so
// Create must return ErrHandled (not a fresh error the caller re-prints) and
// must not POST the time entry.
func TestCreateUnknownActivityReturnsErrHandled(t *testing.T) {
	postCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/api/v3/time_entries/activities":
			_, _ = response.Write([]byte(`{"_embedded":{"elements":[{"id":1,"name":"Development"}]}}`))
		case request.Method == http.MethodPost:
			postCount++
			_, _ = response.Write([]byte(`{"id":1}`))
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

	_, err = Create(map[CreateOption]string{
		CreateActivity: "no-such-activity",
	})
	if !stderrors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("Create error = %v, want ErrHandled", err)
	}
	if postCount != 0 {
		t.Errorf("POST count = %d, want 0", postCount)
	}
}
