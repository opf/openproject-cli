package routes_test

import (
	"net/url"
	"testing"

	"github.com/opf/openproject-cli/components/routes"
	"github.com/opf/openproject-cli/models"
)

func initHost(t *testing.T) {
	t.Helper()
	host, err := url.Parse("https://op.example.com")
	if err != nil {
		t.Fatal(err)
	}
	routes.Init(host)
}

func TestProjectUrl_UsesIdentifierWhenPresent(t *testing.T) {
	initHost(t)

	got := routes.ProjectUrl(&models.Project{Id: 42, Identifier: "my-project"}).String()
	want := "https://op.example.com/projects/my-project"
	if got != want {
		t.Errorf("ProjectUrl = %q, want %q", got, want)
	}
}

// A project whose identifier is missing (older servers, minimal payloads) must
// still produce a working URL by falling back to the numeric id.
func TestProjectUrl_EmptyIdentifierFallsBackToId(t *testing.T) {
	initHost(t)

	got := routes.ProjectUrl(&models.Project{Id: 42, Identifier: ""}).String()
	want := "https://op.example.com/projects/42"
	if got != want {
		t.Errorf("ProjectUrl = %q, want %q", got, want)
	}
}
