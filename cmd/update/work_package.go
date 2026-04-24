package update

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/presenter"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var (
	actionFlag                    string
	assigneeFlag                  uint64
	attachFlag                    string
	dryRunUpdateWorkPackage       bool
	printUpdatedWorkPackageAsJSON bool
	setFlags                      []string
	subjectFlag                   string
	typeFlag                      string
)

var workPackageCmd = &cobra.Command{
	Use:   "workpackage [id]",
	Short: "Updates the work package",
	Long: `Update a work package. Each update
provided by a flag is executed on its own.`,
	Run: updateWorkPackage,
}

func updateWorkPackage(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		printUpdateError("invalid_argument", fmt.Sprintf("Expected 1 argument [id], but got %d", len(args)))
		return
	}

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		printUpdateError("invalid_argument", fmt.Sprintf("'%s' is an invalid work package id. Must be a number.", args[0]))
		return
	}

	if err := validateUpdateWorkPackageFlags(); err != nil {
		printUpdateError("conflicting_arguments", err.Error())
		return
	}

	if len(setFlags) > 0 {
		if dryRunUpdateWorkPackage {
			plan, err := work_packages.DryRunUpdateFields(id, setFlags)
			if err != nil {
				printUpdateError(updateFieldErrorCode(err), err.Error())
				return
			}

			data, err := presenter.MarshalJSON(plan)
			if err != nil {
				printer.Error(err)
				return
			}

			printer.Info(string(data))
			return
		}

		if err := work_packages.UpdateFields(id, setFlags); err != nil {
			printUpdateError(updateFieldErrorCode(err), err.Error())
			return
		}

		if printUpdatedWorkPackageAsJSON {
			payload, err := work_packages.Inspect(id)
			if err != nil {
				printUpdateError("post_apply_inspect_failed", err.Error())
				return
			}

			data, err := presenter.MarshalJSON(payload)
			if err != nil {
				printer.Error(err)
				return
			}

			printer.Info(string(data))
			return
		}

		workPackage, err := work_packages.Lookup(id)
		if err != nil {
			printer.Error(err)
			return
		}

		printer.Info("-- ")
		printer.WorkPackage(workPackage)
		return
	}

	options := updateOptions()
	if len(options) == 0 {
		if printUpdatedWorkPackageAsJSON || dryRunUpdateWorkPackage {
			printUpdateError("invalid_argument", "no updates specified: provide at least one update flag or --set")
			return
		}

		workPackage, err := work_packages.Lookup(id)
		if err != nil {
			printer.Error(err)
			return
		}

		printer.Info("-- ")
		printer.WorkPackage(workPackage)
		return
	}

	if dryRunUpdateWorkPackage {
		plan, err := work_packages.DryRunUpdate(id, options)
		if err != nil {
			printUpdateError("api_error", err.Error())
			return
		}

		data, err := presenter.MarshalJSON(plan)
		if err != nil {
			printer.Error(err)
			return
		}

		printer.Info(string(data))
		return
	}

	if !printUpdatedWorkPackageAsJSON {
		printer.Info("Updating work package ...")
	}

	workPackage, err := work_packages.Update(id, options)
	if err != nil {
		if printUpdatedWorkPackageAsJSON {
			printUpdateError("api_error", err.Error())
			return
		}
		printer.Error(err)
		return
	}

	if printUpdatedWorkPackageAsJSON {
		payload, err := work_packages.Inspect(id)
		if err != nil {
			printUpdateError("post_apply_inspect_failed", err.Error())
			return
		}

		data, err := presenter.MarshalJSON(payload)
		if err != nil {
			printer.Error(err)
			return
		}

		printer.Info(string(data))
		return
	}

	printer.Info("-- ")
	printer.WorkPackage(workPackage)
}

func updateOptions() map[work_packages.UpdateOption]string {
	var options = make(map[work_packages.UpdateOption]string)
	if len(actionFlag) > 0 {
		options[work_packages.UpdateCustomAction] = actionFlag
	}
	if assigneeFlag > 0 {
		options[work_packages.UpdateAssignee] = strconv.FormatUint(assigneeFlag, 10)
	}
	if len(attachFlag) > 0 {
		options[work_packages.UpdateAttachment] = attachFlag
	}
	if len(subjectFlag) > 0 {
		options[work_packages.UpdateSubject] = subjectFlag
	}
	if len(typeFlag) > 0 {
		options[work_packages.UpdateType] = typeFlag
	}

	return options
}

func validateUpdateWorkPackageFlags() error {
	if dryRunUpdateWorkPackage && !printUpdatedWorkPackageAsJSON {
		return fmt.Errorf("cannot use --dry-run without --json")
	}

	if len(setFlags) > 0 && hasLegacyUpdateFlags() {
		return fmt.Errorf("cannot combine --set with %s", strings.Join(activeLegacyUpdateFlags(), ", "))
	}

	return nil
}

func hasLegacyUpdateFlags() bool {
	return len(activeLegacyUpdateFlags()) > 0
}

func activeLegacyUpdateFlags() []string {
	flags := []string{}

	if len(actionFlag) > 0 {
		flags = append(flags, "--action")
	}
	if assigneeFlag > 0 {
		flags = append(flags, "--assignee")
	}
	if len(attachFlag) > 0 {
		flags = append(flags, "--attach")
	}
	if len(subjectFlag) > 0 {
		flags = append(flags, "--subject")
	}
	if len(typeFlag) > 0 {
		flags = append(flags, "--type")
	}

	return flags
}

func printUpdateError(code, message string) {
	if !printUpdatedWorkPackageAsJSON {
		printer.ErrorText(message)
		return
	}

	data, err := presenter.MarshalError(code, message)
	if err != nil {
		printer.Error(err)
		return
	}

	printer.Info(string(data))
}

func updateFieldErrorCode(err error) string {
	switch {
	case errors.Is(err, work_packages.ErrInvalidFieldAssignment):
		return "invalid_argument"
	case errors.Is(err, work_packages.ErrAmbiguousField):
		return "ambiguous_field"
	case errors.Is(err, work_packages.ErrDuplicateField):
		return "duplicate_field"
	case errors.Is(err, work_packages.ErrUnknownField):
		return "unknown_field"
	case errors.Is(err, work_packages.ErrUnsupportedFieldType):
		return "unsupported_field_type"
	case errors.Is(err, work_packages.ErrInvalidFieldValue):
		return "invalid_field_value"
	case errors.Is(err, work_packages.ErrNonWritableField):
		return "non_writable_field"
	default:
		return "api_error"
	}
}
