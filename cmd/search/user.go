package search

import (
	"fmt"
	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/users"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user [searchInput]",
	Short: "Searches for a user",
	Long:  "Searches for a user by id, keyword, or name. Returns a list of possible matches.",
	Run:   searchUser,
}

var keywords = []string{"me"}

func searchUser(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [searchInput], but got %d", len(args)))
		return
	}

	if common.Contains(keywords, args[0]) {
		me, err := users.Me()
		if err != nil {
			printer.Error(err)
		} else {
			printer.User(me)
		}

		return
	}

	collection, err := users.Search(args[0])
	if err != nil {
		printer.Error(err)
		return
	}

	if len(collection) == 0 {
		printer.Info(fmt.Sprintf("No user found for search input %s.", printer.Cyan(args[0])))
	} else {
		printer.Users(collection)
	}
}
