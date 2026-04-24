package work_packages

import (
	"bytes"
	"encoding/json"

	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/models"
)

func DryRunUpdateFields(id uint64, assignments []string) (*models.WorkPackageUpdatePlan, error) {
	workPackage, err := fetch(id)
	if err != nil {
		return nil, err
	}

	schema, err := SchemaFor(workPackage)
	if err != nil {
		return nil, err
	}

	resolved, err := resolveFieldAssignments(schema, assignments)
	if err != nil {
		return nil, err
	}

	return &models.WorkPackageUpdatePlan{
		Valid:          true,
		Operation:      "update",
		WorkPackageID:  id,
		ResolvedFields: resolved,
	}, nil
}

func UpdateFields(id uint64, assignments []string) error {
	workPackage, err := fetch(id)
	if err != nil {
		return err
	}

	schema, err := SchemaFor(workPackage)
	if err != nil {
		return err
	}

	resolved, err := resolveFieldAssignments(schema, assignments)
	if err != nil {
		return err
	}

	patch := map[string]any{
		"lockVersion": workPackage.LockVersion,
	}
	for _, field := range resolved {
		patch[field.APIField] = field.Value
	}

	body, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	_, err = requests.Patch(workPackage.Links.Self.Href, &requests.RequestData{
		ContentType: "application/json",
		Body:        bytes.NewReader(body),
	})
	return err
}
