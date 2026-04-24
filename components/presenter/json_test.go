package presenter_test

import (
	"testing"

	"github.com/opf/openproject-cli/components/presenter"
	"github.com/opf/openproject-cli/models"
)

func TestInspectPayloadJSON(t *testing.T) {
	parentID := uint64(74316)

	payload := models.WorkPackageInspectPayload{
		WorkPackage: models.WorkPackageDetails{
			ID:          74316,
			Subject:     "Expand op CLI to support scripted work package workflows",
			Type:        "Feature",
			Status:      "new",
			Assignee:    "",
			Description: "Body",
			Project: models.ProjectRef{
				ID:         1482,
				Identifier: "cli",
				Name:       "CLI",
			},
			Fields: map[string]any{
				"customField130": float64(3),
			},
			FieldLabels: map[string][]string{
				"Votes": []string{"customField130"},
			},
		},
		Children: []models.WorkPackageSummary{
			{
				ID:       74413,
				Subject:  "Build a reusable SKILL.md based on OpenProject CLI",
				Type:     "Implementation",
				Status:   "new",
				ParentID: &parentID,
			},
		},
	}

	got, err := presenter.MarshalJSON(payload)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"work_package":{"id":74316,"subject":"Expand op CLI to support scripted work package workflows","type":"Feature","status":"new","assignee":"","description":"Body","parent_id":null,"project":{"id":1482,"identifier":"cli","name":"CLI"},"fields":{"customField130":3},"field_labels":{"Votes":["customField130"]}},"children":[{"id":74413,"subject":"Build a reusable SKILL.md based on OpenProject CLI","type":"Implementation","status":"new","parent_id":74316}]}`

	if string(got) != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestErrorJSON(t *testing.T) {
	got, err := presenter.MarshalError("conflicting_arguments", "cannot use --open together with --json")
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"error":{"code":"conflicting_arguments","message":"cannot use --open together with --json"}}`
	if string(got) != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}
