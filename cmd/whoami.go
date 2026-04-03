package cmd

import (
	stderrors "errors"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/configuration"
	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/users"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user and server",
	Long:  "Display the configured OpenProject server and the currently authenticated user.",
	Run:   whoami,
}

func whoami(_ *cobra.Command, _ []string) {
	host, _, err := configuration.ReadConfig()
	if err != nil {
		printer.Error(err)
		return
	}

	if host == "" {
		printer.ErrorText("Not logged in. Run `op login` to authenticate.")
		return
	}

	user, err := users.Me()
	if err != nil {
		printer.Info("Server: " + host)
		var responseErr *openerrors.ResponseError
		if stderrors.As(err, &responseErr) && responseErr.Status() == http.StatusUnauthorized {
			printer.ErrorText("Invalid or expired token. Run `op login` to re-authenticate.")
		} else {
			printer.Error(err)
		}
		return
	}

	printer.Whoami(host, user)
}
