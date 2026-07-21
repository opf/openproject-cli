package printer

import (
	"cmp"
	"slices"
	"strings"

	"github.com/opf/openproject-cli/models"
)

func Activities(activities []*models.Activity, users []*models.User) {
	// The renderers resolve actors with a binary search over user IDs; the
	// API guarantees no response order, so sort a copy here.
	sorted := slices.Clone(users)
	slices.SortFunc(sorted, func(a, b *models.User) int {
		return cmp.Compare(a.Id, b.Id)
	})
	activeRenderer.Activities(activities, sorted)
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
