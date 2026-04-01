package printer

import (
	"fmt"
	"strings"

	"github.com/opf/openproject-cli/models"
)

func TimeEntry(t *models.TimeEntry) {
	activeRenderer.TimeEntry(t)
}

func TimeEntryList(entries []*models.TimeEntry) {
	activeRenderer.TimeEntryList(entries)
}

func printTimeEntry(t *models.TimeEntry, maxIdLength, maxActivityLength, maxProjectLength int) {
	var parts []string

	diff := maxIdLength - idLength(t.Id)
	idStr := fmt.Sprintf("%s#%d", indent(diff), t.Id)
	parts = append(parts, Red(idStr))

	if maxActivityLength > 0 {
		diff = maxActivityLength - len(t.Activity)
		activityStr := Green(strings.ToUpper(t.Activity)) + indent(diff)
		parts = append(parts, activityStr)
	}

	parts = append(parts, Cyan(t.SpentOn.Format("Mon Jan _2")))
	parts = append(parts, fmt.Sprintf("%.2fh", t.Hours.Hours()))

	if maxProjectLength > 0 {
		diff = maxProjectLength - len(t.Project)
		projectStr := Yellow(t.Project) + indent(diff)
		parts = append(parts, projectStr)
	}

	parts = append(parts, Cyan(t.WorkPackage))

	if len(t.Comment) > 0 {
		parts = append(parts, t.Comment)
	}
	if t.Ongoing {
		parts = append(parts, fmt.Sprintf("(%s)", Yellow("ongoing")))
	}

	activePrinter.Println(strings.Join(parts, " "))
}
