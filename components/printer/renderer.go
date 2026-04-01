package printer

import "github.com/opf/openproject-cli/models"

type Renderer interface {
	WorkPackage(*models.WorkPackage)
	WorkPackages([]*models.WorkPackage)
	Project(*models.Project)
	Projects([]*models.Project)
	User(*models.User)
	Users([]*models.User)
	Types([]*models.Type)
	Status(*models.Status)
	StatusList([]*models.Status)
	TimeEntryList([]*models.TimeEntry)
	TimeEntry(*models.TimeEntry)
	Notification(*models.Notification)
	Notifications([]*models.Notification)
	Activities([]*models.Activity, []*models.User)
	CustomActions([]*models.CustomAction)
	Number(int64)
}

var activeRenderer Renderer = &TextRenderer{}

func InitRenderer(format string) {
	switch format {
	case "json":
		activeRenderer = &JsonRenderer{}
	default:
		activeRenderer = &TextRenderer{}
	}
}
