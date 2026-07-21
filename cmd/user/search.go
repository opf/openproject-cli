package user

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/common"
	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/users"
)

var keywords = []string{"me"}

var searchCmd = &cobra.Command{
	Use:   "search [searchInput]",
	Short: "Searches for a user",
	Long:  "Searches for a user by id, keyword, or name. Returns a list of possible matches.",
	RunE:  searchUser,
}

func searchUser(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [searchInput], but got %d", len(args)))
		return openerrors.ErrHandled
	}

	if common.Contains(keywords, args[0]) {
		me, err := users.Me()
		if err != nil {
			printer.Error(err)
			return openerrors.ErrHandled
		}
		printer.User(me)
		return nil
	}

	collection, err := users.Search(args[0])
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	if len(collection) == 0 {
		printer.Info(fmt.Sprintf("No user found for search input %s.", printer.Cyan(args[0])))
	} else {
		printer.Users(collection)
	}
	return nil
}
