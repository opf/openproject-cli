package models

type WorkPackage struct {
	Id          int64
	Subject     string
	Type        string
	Assignee    string
	Status      string
	Description string
}
