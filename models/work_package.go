package models

type WorkPackage struct {
	Id          int64
	Subject     string
	Type        string
	Assignee    Principal
	Status      string
	Description string
	LockVersion int
}
