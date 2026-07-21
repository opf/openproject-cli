package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/configuration"
	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Removes a stored profile",
	Long:  "Removes the credentials for the given profile (defaults to 'default').",
	RunE:  logout,
}

func logout(cmd *cobra.Command, _ []string) error {
	profile, _ := resolvedProfile(cmd)

	ok, err := confirmRemove(profile)
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}
	if !ok {
		return nil
	}

	if err := configuration.DeleteProfile(profile); err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	printer.Done()
	return nil
}

func confirmRemove(profile string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	printer.Input(fmt.Sprintf("Remove profile %q? [y/N] ", profile))
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	answer := strings.ToLower(strings.TrimSpace(common.SanitizeLineBreaks(input)))
	return answer == "y" || answer == "yes", nil
}
