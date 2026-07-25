package printer

import (
	"fmt"

	"github.com/opf/openproject-cli/models"
)

type Renderer interface {
	Budget(*models.Budget)
	Budgets([]*models.Budget)
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
	Whoami(profile, host string, user *models.User)
	WhoamiList(entries []WhoamiEntry)
}

var activeRenderer Renderer = &TextRenderer{}

func InitRenderer(format string) error {
	switch format {
	case "json":
		activeRenderer = &JsonRenderer{}
	case "text", "":
		activeRenderer = &TextRenderer{}
	default:
		activeRenderer = &TextRenderer{}
		ErrorText(fmt.Sprintf("unknown output format %q", format))
		return fmt.Errorf("unknown output format %q", format)
	}
	return nil
}
