package models

type WorkPackageDraft struct {
	Subject string `json:"subject"`
	Type    string `json:"type"`
}

type WorkPackageCreatePlan struct {
	Valid       bool             `json:"valid"`
	Operation   string           `json:"operation"`
	ProjectID   uint64           `json:"project_id"`
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
	WorkPackageID  uint64                   `json:"work_package_id"`
	Subject        string                   `json:"subject,omitempty"`
	Type           string                   `json:"type,omitempty"`
	Assignee       string                   `json:"assignee,omitempty"`
	Description    *string                  `json:"description,omitempty"`
	Action         string                   `json:"action,omitempty"`
	Attach         string                   `json:"attach,omitempty"`
	ResolvedFields map[string]ResolvedField `json:"resolved_fields"`
}
