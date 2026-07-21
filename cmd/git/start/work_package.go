package start

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
)

var startWorkPackageCmd = &cobra.Command{
	Use:   "workpackage [id]",
	Short: "Starts the work on a work package",
	Long:  "Creates a branch based on the work package identifier (numeric ID or project-based, e.g. PROJ-123), subject and type. Switches to the branch after creation.",
	RunE:  startWorkPackage,
}

func startWorkPackage(_ *cobra.Command, args []string) error {
	id, err := checkArgumentsForId(args)
	if err != nil {
		return err
	}

	repo, err := git.PlainOpen(".")
	if err != nil {
		printer.ErrorText("No git repository found in the current location.")
		return openerrors.ErrHandled
	}

	head, err := repo.Head()
	if err := handleGitError(err); err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err := handleGitError(err); err != nil {
		return err
	}

	branchName, err := deriveBranchName(id)
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Hash:   head.Hash(),
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName)),
		Create: true,
		Keep:   true,
	})
	if err := handleGitError(err); err != nil {
		return err
	}

	printer.Info(fmt.Sprintf("Switched to a new branch: %s", branchName))
	return nil
}

func checkArgumentsForId(args []string) (string, error) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [id], but got %d", len(args)))
		return "", openerrors.ErrHandled
	}

	if err := work_packages.ValidateIdentifier(args[0]); err != nil {
		printer.ErrorText(err.Error())
		return "", openerrors.ErrHandled
	}

	return args[0], nil
}

func handleGitError(err error) error {
	if err != nil {
		printer.ErrorText(fmt.Sprintf("Fatal error on executing git commands: %s", err.Error()))
		return openerrors.ErrHandled
	}
	return nil
}

func deriveBranchName(id string) (string, error) {
	workPackage, err := work_packages.Lookup(id)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf(
		"%s/%d-%s",
		sanitizeString(workPackage.Type),
		workPackage.Id,
		sanitizeString(workPackage.Subject),
	)

	return name, nil
}

func sanitizeString(str string) string {
	lower := strings.ToLower(str)
	return regexp.MustCompile("[^a-zA-Z0-9-/_.]").ReplaceAllStringFunc(lower, func(s string) string {
		if s == " " {
			return "-"
		}

		return ""
	})
}
