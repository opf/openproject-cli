package workpackage

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/projects"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var searchProjectId string

var searchCmd = &cobra.Command{
	Use:   "search <query>...",
	Short: "Searches for work packages",
	Long: `Searches for work packages by subject, type, status, project name, or identifier.
Multiple words are ANDed: all terms must match. Returns up to 100 results.`,
	Run: searchWorkPackages,
}

func searchWorkPackages(_ *cobra.Command, args []string) {
	query := strings.Join(args, " ")
	if strings.TrimSpace(query) == "" {
		printer.ErrorText("Search query cannot be blank")
		return
	}

	isProjectScoped := len(searchProjectId) > 0
	if isProjectScoped {
		if err := projects.ValidateIdentifier(searchProjectId); err != nil {
			printer.ErrorText(fmt.Sprintf("--project: %s", err.Error()))
			return
		}
	}

	collection, err := work_packages.Search(query, searchProjectId)
	if err != nil {
		if isNotFound(err) && isProjectScoped {
			printer.ErrorText(fmt.Sprintf("--project: no project found with identifier or ID '%s'", searchProjectId))
		} else {
			printer.Error(err)
		}
		return
	}

	if len(collection) == 0 {
		printer.Info(fmt.Sprintf("No work package found for search input %s.", printer.Cyan(query)))
	} else {
		printer.WorkPackages(collection)
	}
}
