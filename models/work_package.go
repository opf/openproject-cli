package models

type WorkPackage struct {
	Id          uint64
	Subject     string
	Type        string
	Assignee    Principal
	Status      string
	Description string
	LockVersion int
}
