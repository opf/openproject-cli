package models

import "time"

type TimeEntry struct {
	Id          uint64
	Comment     string
	Project     string
	WorkPackage string
	SpentOn     time.Time
	Hours       time.Duration
	Ongoing     bool
	User        string
	Activity    string
	CreatedAt   string
	UpdatedAt   string
}
