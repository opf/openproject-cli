package models

type Status struct {
	Id         uint64
	Name       string
	Color      string
	IsDefault  bool
	IsClosed   bool
	IsReadonly bool
	Position   uint64
}
