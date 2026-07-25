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
	RunE:  whoami,
}

func whoami(cmd *cobra.Command, _ []string) error {
	// OP_CLI_HOST/OP_CLI_TOKEN override every stored profile, so listing file
	// profiles would just query the same server under different labels.
	if configuration.HasEnvironmentConfig() {
		return whoamiSingle("environment")
	}

	profile, explicit := resolvedProfile(cmd)

	if explicit {
		return whoamiSingle(profile)
	}

	// No profile specified: show all profiles
	profiles, err := configuration.AllProfiles()
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	if len(profiles) == 0 {
		printer.ErrorText("No profiles configured. Run `op login` to authenticate.")
		return openerrors.ErrHandled
	}

	var entries []printer.WhoamiEntry
	var failed bool
	for _, p := range profiles {
		entry, err := whoamiEntry(p.Name)
		if err != nil {
			failed = true
			continue
		}
		entries = append(entries, entry)
	}

	if len(entries) > 0 {
		printer.WhoamiList(entries)
	}
	if failed {
		return openerrors.ErrHandled
	}
	return nil
}

// whoamiSingle renders exactly one profile as a list, so JSON output stays a
// single-element array consistent with the all-profiles path.
func whoamiSingle(profile string) error {
	entry, err := whoamiEntry(profile)
	if err != nil {
		return err
	}
	printer.WhoamiList([]printer.WhoamiEntry{entry})
	return nil
}

func whoamiEntry(profile string) (printer.WhoamiEntry, error) {
	var entry printer.WhoamiEntry

	host, token, err := configuration.ReadConfig(profile)
	if err != nil {
		printer.Error(err)
		return entry, openerrors.ErrHandled
	}

	if host == "" {
		printer.ErrorText("Profile \"" + profile + "\" is not configured. Run `op login --profile " + profile + "` to authenticate.")
		return entry, openerrors.ErrHandled
	}

	parse, err := url.Parse(host)
	if err != nil {
		printer.Error(err)
		return entry, openerrors.ErrHandled
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
		return entry, openerrors.ErrHandled
	}

	return printer.WhoamiEntry{Profile: profile, Host: host, User: user}, nil
}
