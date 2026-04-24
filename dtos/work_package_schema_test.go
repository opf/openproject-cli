package dtos

import (
	"encoding/json"
	"testing"
)

func TestWorkPackageSchemaDtoUnmarshalIgnoresNonCustomFields(t *testing.T) {
	body := []byte(`{
		"subject": {"name": "Subject", "type": "String", "writable": true},
		"customField108": {"name": "Requires doc change", "type": "Boolean", "writable": true},
		"customField130": {"name": "Votes", "type": "Integer", "writable": false}
	}`)

	var dto WorkPackageSchemaDto
	if err := json.Unmarshal(body, &dto); err != nil {
		t.Fatal(err)
	}

	if len(dto.Fields) != 2 {
		t.Fatalf("expected two custom fields, got %+v", dto.Fields)
	}

	if _, ok := dto.Fields["subject"]; ok {
		t.Fatalf("expected non-custom fields to be ignored, got %+v", dto.Fields)
	}

	if !dto.Fields["customField108"].Writable {
		t.Fatalf("expected customField108 to be writable, got %+v", dto.Fields["customField108"])
	}

	if dto.Fields["customField130"].Type != "Integer" {
		t.Fatalf("expected customField130 type to be Integer, got %+v", dto.Fields["customField130"])
	}
}
