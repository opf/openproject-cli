package models

type ProjectRef struct {
	ID         uint64 `json:"id"`
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

type WorkPackageSummary struct {
	ID       uint64  `json:"id"`
	Subject  string  `json:"subject"`
	Type     string  `json:"type"`
	Status   string  `json:"status"`
	ParentID *uint64 `json:"parent_id"`
}

type WorkPackageDetails struct {
	ID          uint64              `json:"id"`
	Subject     string              `json:"subject"`
	Type        string              `json:"type"`
	Status      string              `json:"status"`
	Assignee    string              `json:"assignee"`
	Description string              `json:"description"`
	ParentID    *uint64             `json:"parent_id"`
	Project     ProjectRef          `json:"project"`
	Fields      map[string]any      `json:"fields"`
	FieldLabels map[string][]string `json:"field_labels"`
}

type WorkPackageInspectPayload struct {
	WorkPackage WorkPackageDetails   `json:"work_package"`
	Children    []WorkPackageSummary `json:"children"`
}
