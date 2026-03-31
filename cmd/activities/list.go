package activities

import (
	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/users"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var listWpId uint64

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists activities",
	Long:  "Get a list of activities, scoped by the provided flag (e.g. --wp).",
	Run:   listActivities,
}

func listActivities(_ *cobra.Command, _ []string) {
	if listWpId > 0 {
		listWorkPackageActivities(listWpId)
		return
	}

	printer.ErrorText("Please specify a scope. Example: --wp [id]")
}

func listWorkPackageActivities(wpId uint64) {
	acts, err := work_packages.Activities(wpId)
	if err != nil {
		printer.ErrorText(err.Error())
		return
	}

	var userIds []uint64
	for _, a := range acts {
		if a.UserId > 0 {
			userIds = append(userIds, a.UserId)
		}
	}

	userList := users.ByIds(userIds)
	printer.Activities(acts, userList)
}
