package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/configuration"
	"github.com/opf/openproject-cli/components/printer"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Removes a stored profile",
	Long:  "Removes the credentials for the given profile (defaults to 'default').",
	Run:   logout,
}

func logout(cmd *cobra.Command, _ []string) {
	profile, _ := resolvedProfile(cmd)

	ok, err := confirmRemove(profile)
	if err != nil {
		printer.Error(err)
		return
	}
	if !ok {
		return
	}

	if err := configuration.DeleteProfile(profile); err != nil {
		printer.Error(err)
		return
	}

	printer.Done()
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
