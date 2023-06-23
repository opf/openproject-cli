package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/opf/openproject-cli/models"
)

func WorkPackages(workPackages []*models.WorkPackage) {
	for _, workPackage := range workPackages {
		id := fmt.Sprintf("#%d", workPackage.Id)
		fmt.Printf("[%s] %s\n", red(id), cyan(workPackage.Subject))
	}
}

func WorkPackage(workPackage *models.WorkPackage) {
	printHeadline(workPackage, idLength(workPackage.Id), len(workPackage.Type))
	printAttributes(workPackage)
	fmt.Println()
	printDescription(workPackage)
}

func idLength(id int64) int {
	return len(strconv.FormatInt(id, 10)) + 1
}

func printHeadline(workPackage *models.WorkPackage, maxIdLength, maxTypeLength int) {
	diff := maxIdLength - idLength(workPackage.Id)
	idStr := fmt.Sprintf("%s#%d", indent(diff), workPackage.Id)

	diff = maxTypeLength - len(workPackage.Type)
	typeStr := strings.ToUpper(workPackage.Type) + indent(diff)

	fmt.Printf("%s %s %s\n",
		red(idStr),
		green(typeStr),
		cyan(workPackage.Subject))
}

func printAttributes(workPackage *models.WorkPackage) {
	fmt.Printf("[%s]\n", yellow(workPackage.Status))

	assigneeStr := workPackage.Assignee
	if len(assigneeStr) == 0 {
		assigneeStr = "-"
	}
	fmt.Printf("Assignee: %s\n", assigneeStr)
}

func printDescription(workPackage *models.WorkPackage) {
	lines := splitIntoLines(workPackage.Description, 80)
	for _, line := range lines {
		fmt.Printf("%s\n", line)
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
