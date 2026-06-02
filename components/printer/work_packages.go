package printer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/opf/openproject-cli/components/routes"
	"github.com/opf/openproject-cli/models"
)

func WorkPackage(wp *models.WorkPackage) {
	activeRenderer.WorkPackage(wp)
}

func WorkPackages(wps []*models.WorkPackage) {
	activeRenderer.WorkPackages(wps)
}

func idLength(id uint64) int {
	return len(strconv.FormatUint(id, 10)) + 1
}

// formatDisplayId returns the display identifier: "SJF-13" for semantic IDs,
// "#42" for numeric-only systems (where displayId equals the numeric id).
func formatDisplayId(wp *models.WorkPackage) string {
	if wp.DisplayId == strconv.FormatUint(wp.Id, 10) {
		return fmt.Sprintf("#%d", wp.Id)
	}
	return wp.DisplayId
}

func displayIdLength(wp *models.WorkPackage) int {
	return utf8.RuneCountInString(formatDisplayId(wp))
}

func printHeadline(workPackage *models.WorkPackage, maxIdLength, maxStatusLength, maxTypeLength int) {
	var parts []string

	diff := maxIdLength - displayIdLength(workPackage)
	idStr := fmt.Sprintf("%s%s", indent(diff), formatDisplayId(workPackage))
	parts = append(parts, Red(idStr))

	diff = maxTypeLength - utf8.RuneCountInString(workPackage.Type)
	typeStr := strings.ToUpper(workPackage.Type) + indent(diff)
	parts = append(parts, Green(typeStr))

	if maxStatusLength > 0 {
		diff = maxStatusLength - utf8.RuneCountInString(workPackage.Status)
		statusStr := fmt.Sprintf("[%s]%s", Yellow(workPackage.Status), indent(diff))
		parts = append(parts, statusStr)
	}

	parts = append(parts, Cyan(workPackage.Subject))
	activePrinter.Println(strings.Join(parts, " "))
}

func printAttributes(workPackage *models.WorkPackage) {
	activePrinter.Printf("[%s]\n", Yellow(workPackage.Status))

	assigneeStr := workPackage.Assignee
	if len(assigneeStr) == 0 {
		assigneeStr = "-"
	}
	activePrinter.Printf("Assignee: %s\n", assigneeStr)
}

func printOpenLink(workPackage *models.WorkPackage) {
	activePrinter.Printf("Open: %s\n", routes.WorkPackageUrl(workPackage).String())
}

func printDescription(workPackage *models.WorkPackage) {
	for _, line := range splitIntoLines(workPackage.Description, 80) {
		activePrinter.Printf("%s\n", line)
	}
}

func splitWords(text string, lineLength int) []string {
	words := strings.Fields(text)
	var lines []string
	var line string
	for _, word := range words {
		if len(line)+len(word)+1 > lineLength {
			lines = append(lines, line)
			line = ""
		}
		if len(line) > 0 {
			line += " "
		}
		line += word
	}
	if len(line) > 0 {
		lines = append(lines, line)
	}
	return lines
}

func splitIntoLines(text string, lineLength int) []string {
	var lines []string
	for _, paragraph := range strings.Split(text, "\n") {
		split := splitWords(paragraph, lineLength)
		if len(split) == 0 {
			lines = append(lines, "")
		} else {
			lines = append(lines, split...)
		}
	}
	return lines
}
