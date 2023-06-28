package printer

import (
	"sort"
	"strings"

	"github.com/opf/openproject-cli/models"
)

func Activities(activities []*models.Activity, users []*models.User) {
	for _, activity := range activities {
		user := &models.User{
			Id:        0,
			Name:      "",
			FirstName: "",
			LastName:  "",
		}
		if activity.UserId > 0 {
			userIndex := sort.Search(len(users)-1, func(i int) bool { return users[i].Id == activity.UserId })
			user = users[userIndex]
		}
		printActivityHeadline(activity, user)
		printActivityBody(activity)
		println("")
	}
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
