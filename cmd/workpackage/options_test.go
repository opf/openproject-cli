package workpackage

import (
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

func newCmdWithDescriptionFlag(flag *string) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringVar(flag, "description", "", "")
	return cmd
}

func TestCreateOptions_DescriptionOmitted(t *testing.T) {
	createDescriptionFlag = ""
	cmd := newCmdWithDescriptionFlag(&createDescriptionFlag)
	_ = cmd.Flags().Parse([]string{})

	options := createOptions(cmd, "subject")

	if _, ok := options[work_packages.CreateDescription]; ok {
		t.Error("expected CreateDescription to be absent when flag not provided")
	}
}

func TestCreateOptions_DescriptionProvided(t *testing.T) {
	createDescriptionFlag = ""
	cmd := newCmdWithDescriptionFlag(&createDescriptionFlag)
	_ = cmd.Flags().Parse([]string{"--description", "## Hello"})

	options := createOptions(cmd, "subject")

	val, ok := options[work_packages.CreateDescription]
	if !ok {
		t.Error("expected CreateDescription to be present when flag provided")
	}
	if val != "## Hello" {
		t.Errorf("expected %q, got %q", "## Hello", val)
	}
}

func TestUpdateOptions_DescriptionOmitted(t *testing.T) {
	updateDescriptionFlag = ""
	cmd := newCmdWithDescriptionFlag(&updateDescriptionFlag)
	_ = cmd.Flags().Parse([]string{})

	options := updateOptions(cmd)

	if _, ok := options[work_packages.UpdateDescription]; ok {
		t.Error("expected UpdateDescription to be absent when flag not provided")
	}
}

func TestUpdateOptions_DescriptionProvided(t *testing.T) {
	updateDescriptionFlag = ""
	cmd := newCmdWithDescriptionFlag(&updateDescriptionFlag)
	_ = cmd.Flags().Parse([]string{"--description", ""})

	options := updateOptions(cmd)

	if _, ok := options[work_packages.UpdateDescription]; !ok {
		t.Error("expected UpdateDescription to be present when flag explicitly provided with empty string")
	}
}

func TestUpdateWithoutOptionsStopsBeforeRequest(t *testing.T) {
	updateActionFlag = ""
	updateAssigneeFlag = 0
	updateAttachFlag = ""
	updateDescriptionFlag = ""
	updateSubjectFlag = ""
	updateTypeFlag = ""

	cmd := newCmdWithDescriptionFlag(&updateDescriptionFlag)
	if err := cmd.Flags().Parse(nil); err != nil {
		t.Fatal(err)
	}

	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)

	err := updateWorkPackage(cmd, []string{"42"})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("updateWorkPackage error = %v, want ErrHandled", err)
	}
	if !strings.Contains(testingPrinter.ErrResult, "No update options provided") {
		t.Errorf("stderr = %q, want no-options diagnostic", testingPrinter.ErrResult)
	}
	if strings.Contains(testingPrinter.ErrResult, "Cannot execute requests") {
		t.Errorf("command attempted a request: %q", testingPrinter.ErrResult)
	}
}
