package work_packages_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/dtos"
)

func TestSchemaForWithoutSchemaLinkReturnsEmptySchema(t *testing.T) {
	schema, err := work_packages.SchemaFor(&dtos.WorkPackageDto{})
	if err != nil {
		t.Fatal(err)
	}

	if len(schema.Fields) != 0 {
		t.Fatalf("expected empty schema, got %+v", schema.Fields)
	}
}

func TestSchemaForSortsCustomFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/work_packages/schemas/1482-6":
			_, _ = io.WriteString(w, `{
				"customField130": {"name": "Votes", "type": "Integer", "writable": true},
				"customField108": {"name": "Requires doc change", "type": "Boolean", "writable": true},
				"subject": {"name": "Subject", "type": "String", "writable": true}
			}`)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	requests.Init(host, "token", false)

	schema, err := work_packages.SchemaFor(&dtos.WorkPackageDto{
		Links: &dtos.WorkPackageLinksDto{
			Schema: &dtos.LinkDto{Href: "/api/v3/work_packages/schemas/1482-6"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(schema.Fields) != 2 {
		t.Fatalf("expected two custom fields, got %+v", schema.Fields)
	}

	if schema.Fields[0].APIName != "customField108" || schema.Fields[1].APIName != "customField130" {
		t.Fatalf("expected sorted fields, got %+v", schema.Fields)
	}
}
