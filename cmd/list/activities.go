package list

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/users"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var activitiesCmd = &cobra.Command{
	Use:     "activities [workPackageId]",
	Aliases: []string{"ac"},
	Short:   "Lists activities for work package",
	Long:    `Get a list of activities for a work package.`,
	Run:     listActivities,
}

func listActivities(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [workPackageId], but got %d", len(args)))
		return
	}

	wpId, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		printer.ErrorText(err.Error())
	}
	activities, err := work_packages.Activities(wpId)
	if err != nil {
		printer.ErrorText(err.Error())
	}

	var userIds []uint64
	for _, a := range activities {
		if a.UserId > 0 {
			userIds = append(userIds, a.UserId)
			continue
		}
	}

	userList := users.ByIds(userIds)
	printer.Activities(activities, userList)
}
