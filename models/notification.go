package models

import "github.com/opf/openproject-cli/models/types"

type Notification struct {
	Id              int64
	ResourceId      int64
	ResourceSubject string
	Reason          types.Reason
	Read            bool
	CreatedAt       string
	UpdatedAt       string
}
