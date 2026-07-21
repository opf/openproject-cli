package timeentry

import (
	stderrors "errors"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/time_entries"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var createWorkPackageId string
var createHours float64
var createActivity string
var createSpentOn string
var createUserId uint64
var createComment string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a time entry",
	Long:  "Log time spent on a work package.",
	RunE:  createTimeEntry,
}

func createTimeEntry(cmd *cobra.Command, _ []string) error {
	if err := work_packages.ValidateIdentifier(createWorkPackageId); err != nil {
		printer.ErrorText(err.Error())
		return openerrors.ErrHandled
	}

	options := map[time_entries.CreateOption]string{
		time_entries.CreateWorkPackage: createWorkPackageId,
		time_entries.CreateHours:       strconv.FormatFloat(createHours, 'f', -1, 64),
	}

	if cmd.Flags().Changed("activity") {
		options[time_entries.CreateActivity] = createActivity
	}

	if cmd.Flags().Changed("spent-on") {
		options[time_entries.CreateSpentOn] = createSpentOn
	} else {
		options[time_entries.CreateSpentOn] = time.Now().Format(time.DateOnly)
	}

	if createUserId > 0 {
		options[time_entries.CreateUser] = strconv.FormatUint(createUserId, 10)
	}

	if cmd.Flags().Changed("comment") {
		options[time_entries.CreateComment] = createComment
	}

	entry, err := time_entries.Create(options)
	if err != nil {
		if !stderrors.Is(err, openerrors.ErrHandled) {
			printer.Error(err)
		}
		return openerrors.ErrHandled
	}

	printer.TimeEntry(entry)
	printer.Done()
	return nil
}
