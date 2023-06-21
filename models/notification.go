package models

type Notification struct {
	Id              int64
	ResourceId      int64
	ResourceSubject string
	Reason          string
	Read            bool
	CreatedAt       string
	UpdatedAt       string
}
