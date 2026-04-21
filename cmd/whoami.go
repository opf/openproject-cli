package cmd

import (
	stderrors "errors"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/configuration"
	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources/users"
	"github.com/opf/openproject-cli/components/routes"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user and server",
	Long:  "Display the configured OpenProject server and the currently authenticated user.",
	Run:   whoami,
}

func whoami(cmd *cobra.Command, _ []string) {
	profile, explicit := resolvedProfile(cmd)

	if explicit {
		whoamiOne(profile)
		return
	}

	// No profile specified: show all profiles
	profiles, err := configuration.AllProfiles()
	if err != nil {
		printer.Error(err)
		return
	}

	if len(profiles) == 0 {
		printer.ErrorText("No profiles configured. Run `op login` to authenticate.")
		return
	}

	for i, p := range profiles {
		if i > 0 {
			printer.Info("")
		}
		whoamiOne(p.Name)
	}
}

func whoamiOne(profile string) {
	host, token, err := configuration.ReadConfig(profile)
	if err != nil {
		printer.Error(err)
		return
	}

	if host == "" {
		printer.ErrorText("Profile \"" + profile + "\" is not configured. Run `op login --profile " + profile + "` to authenticate.")
		return
	}

	parse, err := url.Parse(host)
	if err != nil {
		printer.Error(err)
		return
	}
	requests.Init(parse, token, Verbose)
	routes.Init(parse)

	user, err := users.Me()
	if err != nil {
		printer.Info("Server: " + host)
		var responseErr *openerrors.ResponseError
		if stderrors.As(err, &responseErr) && responseErr.Status() == http.StatusUnauthorized {
			printer.ErrorText("Invalid or expired token. Run `op login --profile " + profile + "` to re-authenticate.")
		} else {
			printer.Error(err)
		}
		return
	}

	printer.Whoami(profile, host, user)
}
