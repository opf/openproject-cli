package dtos

import (
	"encoding/json"
	"testing"
)

func TestWorkPackageDtoUnmarshalCapturesInspectFields(t *testing.T) {
	body := []byte(`{
		"id": 74316,
		"subject": "Expand op CLI to support scripted work package workflows",
		"description": {"raw": "Body"},
		"customField130": 3,
		"customField108": false,
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
			"parent": {"href": "/api/v3/work_packages/70000", "title": "Umbrella"},
			"schema": {"href": "/api/v3/work_packages/schemas/1482-6"},
			"status": {"href": "/api/v3/statuses/1", "title": "new"},
			"type": {"href": "/api/v3/types/6", "title": "Feature"},
			"assignee": {"href": null, "title": ""}
		}
	}`)

	var dto WorkPackageDto
	if err := json.Unmarshal(body, &dto); err != nil {
		t.Fatal(err)
	}

	if dto.Embedded == nil || dto.Embedded.Project == nil {
		t.Fatalf("expected embedded project, got %+v", dto.Embedded)
	}

	if dto.Embedded.Project.Identifier != "cli" {
		t.Fatalf("expected project identifier cli, got %q", dto.Embedded.Project.Identifier)
	}

	if dto.Links == nil || dto.Links.Parent == nil || dto.Links.Parent.Href != "/api/v3/work_packages/70000" {
		t.Fatalf("expected parent link to be captured, got %+v", dto.Links)
	}

	if dto.Links == nil || dto.Links.Schema == nil || dto.Links.Schema.Href != "/api/v3/work_packages/schemas/1482-6" {
		t.Fatalf("expected schema link to be captured, got %+v", dto.Links)
	}

	if dto.CustomFields["customField130"] != float64(3) {
		t.Fatalf("expected customField130 to be captured, got %#v", dto.CustomFields["customField130"])
	}

	if dto.CustomFields["customField108"] != false {
		t.Fatalf("expected customField108 to be captured, got %#v", dto.CustomFields["customField108"])
	}
}
