package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/routes"
	"github.com/opf/openproject-cli/models"
)

func WorkPackages(workPackages []*models.WorkPackage) {
	var maxIdLength = 0
	var maxTypeLength = 0
	var maxStatusLength = 0
	for _, w := range workPackages {
		maxIdLength = common.Max(maxIdLength, idLength(w.Id))
		maxTypeLength = common.Max(maxTypeLength, len(w.Type))
		maxStatusLength = common.Max(maxStatusLength, len(w.Status))
	}

	for _, workPackage := range workPackages {
		printHeadline(workPackage, maxIdLength, maxStatusLength, maxTypeLength)
	}
}

func WorkPackage(workPackage *models.WorkPackage) {
	printHeadline(workPackage, idLength(workPackage.Id), 0, len(workPackage.Type))
	printAttributes(workPackage)
	activePrinter.Println()
	printOpenLink(workPackage)
	activePrinter.Println()
	printDescription(workPackage)
}

func idLength(id uint64) int {
	return len(strconv.FormatUint(id, 10)) + 1
}

func printHeadline(workPackage *models.WorkPackage, maxIdLength, maxStatusLength, maxTypeLength int) {
	var parts []string

	diff := maxIdLength - idLength(workPackage.Id)
	idStr := fmt.Sprintf("%s#%d", indent(diff), workPackage.Id)
	parts = append(parts, Red(idStr))

	diff = maxTypeLength - len(workPackage.Type)
	typeStr := strings.ToUpper(workPackage.Type) + indent(diff)
	parts = append(parts, Green(typeStr))

	if maxStatusLength > 0 {
		diff = maxStatusLength - len(workPackage.Status)
		statusStr := fmt.Sprintf("[%s]%s", Yellow(workPackage.Status), indent(diff))
		parts = append(parts, statusStr)
	}

	parts = append(parts, Cyan(workPackage.Subject))
	activePrinter.Println(strings.Join(parts, " "))
}

func printAttributes(workPackage *models.WorkPackage) {
	activePrinter.Printf("[%s]\n", Yellow(workPackage.Status))

	assigneeStr := workPackage.Assignee.Name
	if len(assigneeStr) == 0 {
		assigneeStr = "-"
	}
	activePrinter.Printf("Assignee: %s\n", assigneeStr)
}

func printOpenLink(workPackage *models.WorkPackage) {
	activePrinter.Printf("Open: %s\n", routes.WorkPackageUrl(workPackage).String())
}

func printDescription(workPackage *models.WorkPackage) {
	lines := splitIntoLines(workPackage.Description, 80)
	for _, line := range lines {
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
	paragraphs := strings.Split(text, "\n")

	var lines []string

	for _, paragraph := range paragraphs {
		splitParagraph := splitWords(paragraph, lineLength)

		if len(splitParagraph) == 0 {
			lines = append(lines, "") // Append empty line
		} else {
			lines = append(lines, splitParagraph...)
		}
	}

	return lines
}
