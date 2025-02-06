package printer

import (
	"fmt"
	"strings"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/models"
)

func TimeEntryList(timeEntry []*models.TimeEntry) {
	var maxIdLength = 0
	var maxActivityLength = 0
	var maxProjectLength = 0
	for _, t := range timeEntry {
		maxIdLength = common.Max(maxIdLength, idLength(t.Id))
		maxActivityLength = common.Max(maxActivityLength, len(t.Activity))
		maxProjectLength = common.Max(maxProjectLength, len(t.Project))
	}

	for _, t := range timeEntry {
		printTimeEntry(t, maxIdLength, maxActivityLength, maxProjectLength)
	}
}

func TimeEntry(timeEntry *models.TimeEntry) {
	printTimeEntry(timeEntry, idLength(timeEntry.Id), len(timeEntry.Activity), len(timeEntry.Project))
}

func printTimeEntry(timeEntry *models.TimeEntry, maxIdLength int, maxActivityLength int, maxProjectLength int) {
	var parts []string

	diff := maxIdLength - idLength(timeEntry.Id)
	idStr := fmt.Sprintf("%s#%d", indent(diff), timeEntry.Id)

	parts = append(parts, Red(idStr))

	if maxActivityLength > 0 {
		diff = maxActivityLength - len(timeEntry.Activity)
		activityStr := Green(strings.ToUpper(timeEntry.Activity)) + indent(diff)
		parts = append(parts, activityStr)
	}

	parts = append(parts, Cyan(timeEntry.SpentOn.Format("Mon Jan _2")))

	hoursStr := fmt.Sprintf("%.2fh", timeEntry.Hours.Hours())
	parts = append(parts, hoursStr)

	if maxProjectLength > 0 {
		diff = maxProjectLength - len(timeEntry.Project)
		projectStr := Yellow(timeEntry.Project) + indent(diff)
		parts = append(parts, projectStr)
	}

	parts = append(parts, Cyan(timeEntry.WorkPackage))

	if len(timeEntry.Comment) > 0 {
		parts = append(parts, timeEntry.Comment)
	}

	if timeEntry.Ongoing {
		parts = append(parts, fmt.Sprintf("(%s)", Yellow("ongoing")))
	}

	activePrinter.Println(strings.Join(parts, " "))
}
