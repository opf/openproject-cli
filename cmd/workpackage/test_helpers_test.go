package workpackage

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/routes"
)

func initWorkPackageTestServer(t *testing.T, handler http.HandlerFunc) (*printer.TestingPrinter, *int) {
	t.Helper()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		requestCount++
		handler(response, request)
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)
	routes.Init(host)

	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)
	if err := printer.InitRenderer("json"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = printer.InitRenderer("text")
	})

	return testingPrinter, &requestCount
}
