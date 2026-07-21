package workpackage

import (
	"net/http"
	"strings"
	"testing"
)

func TestSearchWorkPackagesRendersEmptyJSONArray(t *testing.T) {
	testingPrinter, _ := initWorkPackageTestServer(t, func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"_embedded":{"elements":[]},"total":0,"count":0}`))
	})

	if err := searchWorkPackages(nil, []string{"missing"}); err != nil {
		t.Fatalf("searchWorkPackages returned an error: %v", err)
	}

	if testingPrinter.Result != "[]\n" {
		t.Errorf("stdout = %q, want %q", testingPrinter.Result, "[]\n")
	}
	if !strings.Contains(testingPrinter.ErrResult, "No work package found") {
		t.Errorf("stderr should describe the empty result, got %q", testingPrinter.ErrResult)
	}
}
