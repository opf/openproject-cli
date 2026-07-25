package models

type WorkPackageDraft struct {
	Subject     string `json:"subject"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type WorkPackageCreatePlan struct {
	Valid       bool             `json:"valid"`
	Operation   string           `json:"operation"`
	ProjectID   string           `json:"project_id"`
	ParentID    *uint64          `json:"parent_id"`
	WorkPackage WorkPackageDraft `json:"work_package"`
}

type ResolvedField struct {
	APIField string `json:"api_field"`
	Value    any    `json:"value"`
}

type WorkPackageUpdatePlan struct {
	Valid          bool                     `json:"valid"`
	Operation      string                   `json:"operation"`
	WorkPackageID  string                   `json:"work_package_id"`
	Subject        string                   `json:"subject,omitempty"`
	Type           string                   `json:"type,omitempty"`
	Assignee       string                   `json:"assignee,omitempty"`
	Status         string                   `json:"status,omitempty"`
	Description    *string                  `json:"description,omitempty"`
	Action         string                   `json:"action,omitempty"`
	Attach         string                   `json:"attach,omitempty"`
	ResolvedFields map[string]ResolvedField `json:"resolved_fields"`
}
