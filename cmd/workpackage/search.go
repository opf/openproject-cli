package workpackage

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>...",
	Short: "Searches for work packages",
	Long: `Searches for work packages by subject, type, status, project name, or identifier.
Multiple words are ANDed: all terms must match. Returns up to 100 results.`,
	Run: searchWorkPackages,
}

func searchWorkPackages(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		printer.ErrorText("Expected at least 1 argument [searchInput], but got 0")
		return
	}

	collection, err := work_packages.Search(strings.Join(args, " "))
	if err != nil {
		printer.Error(err)
		return
	}

	if len(collection) == 0 {
		printer.Info(fmt.Sprintf("No work package found for search input %s.", printer.Cyan(args[0])))
	} else {
		printer.WorkPackages(collection)
	}
}
