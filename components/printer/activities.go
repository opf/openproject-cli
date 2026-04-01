package printer

import (
	"strings"

	"github.com/opf/openproject-cli/models"
)

func Activities(activities []*models.Activity, users []*models.User) {
	activeRenderer.Activities(activities, users)
}

func printActivityHeadline(activity *models.Activity, user *models.User) {
	var parts []string
	if len(user.Name) > 0 {
		parts = append(parts, Green(user.Name))
	}
	parts = append(parts, Yellow(activity.UpdatedAt))
	activePrinter.Println(strings.Join(parts, " "))
}

func printActivityBody(activity *models.Activity) {
	var parts []string
	if len(activity.Comment) > 0 {
		parts = append(parts, Yellow(activity.Comment))
		if len(activity.Details) > 0 {
			parts = append(parts, "---")
		}
	}
	var detailsParts []string
	for _, detail := range activity.Details {
		detailsParts = append(detailsParts, *detail)
	}
	parts = append(parts, strings.Join(detailsParts, "\n"))
	activePrinter.Println(strings.Join(parts, "\n \n"))
}
