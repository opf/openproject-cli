package printer_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/models"
)

func samplePayload() *models.WorkPackageInspectPayload {
	parentID := uint64(74300)
	return &models.WorkPackageInspectPayload{
		WorkPackage: models.WorkPackageDetails{
			ID:          74316,
			Subject:     "Expand op CLI",
			Type:        "Feature",
			Status:      "New",
			Assignee:    "Jane Doe",
			Description: "Body",
			ParentID:    &parentID,
			Project: models.ProjectRef{
				ID:         1482,
				Identifier: "cli",
				Name:       "CLI",
			},
			Fields: map[string]any{
				"customField130": float64(3),
			},
			FieldLabels: map[string][]string{
				"Votes": {"customField130"},
			},
		},
		Children: []models.WorkPackageSummary{
			{ID: 74413, Subject: "Build a reusable SKILL.md", Type: "Implementation", Status: "New"},
		},
	}
}

func TestWorkPackageDetails_Json_IncludesChildrenAndCustomFields(t *testing.T) {
	if err := printer.InitRenderer("json"); err != nil {
		t.Fatalf("InitRenderer(json) failed: %v", err)
	}
	defer func() {
		_ = printer.InitRenderer("text")
	}()
	testingPrinter.Reset()

	printer.WorkPackageDetails(samplePayload())

	if !strings.Contains(testingPrinter.Result, `"children"`) {
		t.Errorf("expected output to contain \"children\", got: %s", testingPrinter.Result)
	}
	if !strings.Contains(testingPrinter.Result, "customField130") {
		t.Errorf("expected output to contain custom field key customField130, got: %s", testingPrinter.Result)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(testingPrinter.Result), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %q", err, testingPrinter.Result)
	}
}

func TestWorkPackageDetails_Text_ListsChildren(t *testing.T) {
	testingPrinter.Reset()

	printer.WorkPackageDetails(samplePayload())

	if !strings.Contains(testingPrinter.Result, "Children (1):") {
		t.Errorf("expected output to contain \"Children (1):\", got: %s", testingPrinter.Result)
	}
	if !strings.Contains(testingPrinter.Result, "#74413 [New] Build a reusable SKILL.md") {
		t.Errorf("expected output to list child #74413, got: %s", testingPrinter.Result)
	}
	if !strings.Contains(testingPrinter.Result, "Votes: 3") {
		t.Errorf("expected output to contain custom field line \"Votes: 3\", got: %s", testingPrinter.Result)
	}
}

func TestWorkPackageCreatePlan_Text_PrintsDryRunFields(t *testing.T) {
	testingPrinter.Reset()

	parentID := uint64(74300)
	plan := &models.WorkPackageCreatePlan{
		Valid:     true,
		Operation: "create",
		ProjectID: "cli",
		ParentID:  &parentID,
		WorkPackage: models.WorkPackageDraft{
			Subject:     "New task",
			Type:        "Task",
			Description: "",
		},
	}

	printer.WorkPackageCreatePlan(plan)

	if !strings.Contains(testingPrinter.Result, "Dry run — no changes applied.") {
		t.Errorf("expected dry-run banner, got: %s", testingPrinter.Result)
	}
	if !strings.Contains(testingPrinter.Result, "project_id: cli") {
		t.Errorf("expected project_id line, got: %s", testingPrinter.Result)
	}
	if !strings.Contains(testingPrinter.Result, "parent_id: 74300") {
		t.Errorf("expected parent_id line, got: %s", testingPrinter.Result)
	}
	if !strings.Contains(testingPrinter.Result, "subject: New task") {
		t.Errorf("expected subject line, got: %s", testingPrinter.Result)
	}
	if strings.Contains(testingPrinter.Result, "description:") {
		t.Errorf("expected empty description to be omitted, got: %s", testingPrinter.Result)
	}
}

func TestWorkPackageUpdatePlan_Json_EmitsResolvedFields(t *testing.T) {
	if err := printer.InitRenderer("json"); err != nil {
		t.Fatalf("InitRenderer(json) failed: %v", err)
	}
	defer func() {
		_ = printer.InitRenderer("text")
	}()
	testingPrinter.Reset()

	plan := &models.WorkPackageUpdatePlan{
		Valid:         true,
		Operation:     "update",
		WorkPackageID: "74316",
		ResolvedFields: map[string]models.ResolvedField{
			"Story points": {APIField: "customField130", Value: float64(5)},
		},
	}

	printer.WorkPackageUpdatePlan(plan)

	if !strings.Contains(testingPrinter.Result, "customField130") {
		t.Errorf("expected resolved field api name in JSON output, got: %s", testingPrinter.Result)
	}
}

func TestWorkPackageUpdatePlan_Text_SkipsEmptyFields(t *testing.T) {
	testingPrinter.Reset()

	plan := &models.WorkPackageUpdatePlan{
		Valid:         true,
		Operation:     "update",
		WorkPackageID: "74316",
		Status:        "Closed",
	}

	printer.WorkPackageUpdatePlan(plan)

	if !strings.Contains(testingPrinter.Result, "status: Closed") {
		t.Errorf("expected status line, got: %s", testingPrinter.Result)
	}
	if strings.Contains(testingPrinter.Result, "subject:") {
		t.Errorf("expected empty subject to be omitted, got: %s", testingPrinter.Result)
	}
	if strings.Contains(testingPrinter.Result, "assignee:") {
		t.Errorf("expected empty assignee to be omitted, got: %s", testingPrinter.Result)
	}
}
