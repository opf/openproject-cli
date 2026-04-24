package work_packages

import (
	"errors"
	"testing"
)

func TestResolveFieldAssignmentsByLabel(t *testing.T) {
	schema := &Schema{
		Fields: []SchemaField{
			{APIName: "customField130", Label: "Votes", Type: "Integer", Writable: true},
		},
	}

	resolved, err := resolveFieldAssignments(schema, []string{"Votes=3"})
	if err != nil {
		t.Fatal(err)
	}

	field := resolved["Votes"]
	if field.APIField != "customField130" || field.Value != int64(3) {
		t.Fatalf("unexpected resolution: %+v", field)
	}
}

func TestResolveFieldAssignmentsRejectsAmbiguousLabels(t *testing.T) {
	schema := &Schema{
		Fields: []SchemaField{
			{APIName: "customField17", Label: "KPI", Type: "Integer", Writable: true},
			{APIName: "customField22", Label: "KPI", Type: "Integer", Writable: true},
		},
	}

	if _, err := resolveFieldAssignments(schema, []string{"KPI=42"}); err == nil {
		t.Fatal("expected ambiguous field error")
	} else if !errors.Is(err, ErrAmbiguousField) {
		t.Fatalf("expected ErrAmbiguousField, got %v", err)
	}
}

func TestResolveFieldAssignmentsRejectsNonWritableFields(t *testing.T) {
	schema := &Schema{
		Fields: []SchemaField{
			{APIName: "customField130", Label: "Votes", Type: "Integer", Writable: false},
		},
	}

	if _, err := resolveFieldAssignments(schema, []string{"Votes=3"}); err == nil {
		t.Fatal("expected non-writable field error")
	} else if !errors.Is(err, ErrNonWritableField) {
		t.Fatalf("expected ErrNonWritableField, got %v", err)
	}
}

func TestResolveFieldAssignmentsRejectsDuplicateAPIFields(t *testing.T) {
	schema := &Schema{
		Fields: []SchemaField{
			{APIName: "customField130", Label: "Votes", Type: "Integer", Writable: true},
		},
	}

	if _, err := resolveFieldAssignments(schema, []string{"Votes=3", "customField130=4"}); err == nil {
		t.Fatal("expected duplicate field error")
	} else if !errors.Is(err, ErrDuplicateField) {
		t.Fatalf("expected ErrDuplicateField, got %v", err)
	}
}
