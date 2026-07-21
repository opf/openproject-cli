package activity

import (
	stderrors "errors"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/users"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var listWpId string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists activities",
	Long:  "Get a list of activities, scoped by the provided flag (e.g. --work-package).",
	RunE:  listActivities,
}

func listActivities(_ *cobra.Command, _ []string) error {
	if err := work_packages.ValidateIdentifier(listWpId); err != nil {
		printer.ErrorText(err.Error())
		return openerrors.ErrHandled
	}
	return listWorkPackageActivities(listWpId)
}

func listWorkPackageActivities(wpId string) error {
	acts, err := work_packages.Activities(wpId)
	if err != nil {
		if !stderrors.Is(err, openerrors.ErrHandled) {
			printer.Error(err)
		}
		return openerrors.ErrHandled
	}

	var userIds []uint64
	for _, a := range acts {
		if a.UserId > 0 {
			userIds = append(userIds, a.UserId)
		}
	}

	userList, err := users.ByIds(userIds)
	if err != nil {
		if !stderrors.Is(err, openerrors.ErrHandled) {
			printer.Error(err)
		}
		return openerrors.ErrHandled
	}
	printer.Activities(acts, userList)
	return nil
}
