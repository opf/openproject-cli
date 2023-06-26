package models

type Notification struct {
	Id              uint64
	ResourceId      uint64
	ResourceSubject string
	Reason          string
	Read            bool
	CreatedAt       string
	UpdatedAt       string
}
