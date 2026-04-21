package cmd

import (
	"fmt"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/cmd/activities"
	"github.com/opf/openproject-cli/cmd/budget"
	"github.com/opf/openproject-cli/cmd/git"
	"github.com/opf/openproject-cli/cmd/notification"
	"github.com/opf/openproject-cli/cmd/project"
	"github.com/opf/openproject-cli/cmd/status"
	"github.com/opf/openproject-cli/cmd/timeentry"
	"github.com/opf/openproject-cli/cmd/user"
	"github.com/opf/openproject-cli/cmd/workpackage"
	"github.com/opf/openproject-cli/cmd/wptype"
	"github.com/opf/openproject-cli/components/configuration"
	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/routes"
)

var Verbose bool
var showVersionFlag bool
var outputFormat string
var profileName string

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "An easy-to-use CLI for the OpenProject APIv3",
	Long: `OpenProject CLI is a fast, reliable and easy-to-use
tool to manage your work packages, notifications and
projects of your OpenProject instance.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		printer.InitRenderer(outputFormat)

		// login and logout manage their own profile and requests setup
		if cmd.Name() == "login" || cmd.Name() == "logout" {
			return nil
		}

		profile, explicit := resolvedProfile(cmd)

		if err := configuration.ValidateProfileName(profile); err != nil {
			printer.Error(err)
			os.Exit(1)
		}

		host, token, err := configuration.ReadConfig(profile)
		if err != nil {
			printer.Error(err)
			os.Exit(1)
		}

		if host == "" && explicit {
			printer.Error(openerrors.Custom(fmt.Sprintf(
				"Profile %q not found. Run 'op login --profile %s' to create it.",
				profile, profile,
			)))
			os.Exit(1)
		}

		parse, _ := url.Parse(host)
		requests.Init(parse, token, Verbose)
		routes.Init(parse)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if showVersionFlag {
			versionText := fmt.Sprintf(
				"%s: %s\n\tcommit: %s\n\tbuilt: %s\n\tbuilt with: %s",
				"OpenProject CLI",
				configuration.CliVersion.Version,
				configuration.CliVersion.Commit,
				configuration.CliVersion.Date.Format(time.UnixDate),
				runtime.Version(),
			)

			printer.Info(printer.Yellow(versionText))

			return
		}

		cmd.Help()
	},
}

// resolvedProfile returns the effective profile name for the current
// invocation, and whether it was explicitly specified by the user (flag or
// OP_CLI_PROFILE env var) vs falling through to the "default" fallback.
func resolvedProfile(cmd *cobra.Command) (profile string, explicit bool) {
	if cmd.Root().PersistentFlags().Changed("profile") {
		return profileName, true
	}
	if env := os.Getenv(configuration.EnvProfile); env != "" {
		return env, true
	}
	return configuration.DefaultProfile, false
}

func Execute(version *configuration.Version) error {
	configuration.Init(version)

	return rootCmd.Execute()
}

func init() {
	activePrinter := &printer.ConsolePrinter{}
	printer.Init(activePrinter)

	rootCmd.Flags().BoolVarP(
		&showVersionFlag,
		"version",
		"",
		false,
		"Show version information of the OpenProject CLI",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&Verbose,
		"verbose",
		"",
		false,
		"Print verbose information of any process that supports this output.",
	)

	rootCmd.PersistentFlags().StringVarP(
		&outputFormat,
		"format",
		"",
		"text",
		`Output format. Accepted values: text, json`,
	)

	rootCmd.PersistentFlags().StringVarP(
		&profileName,
		"profile",
		"",
		configuration.DefaultProfile,
		"Profile name to use (overrides OP_CLI_PROFILE env var)",
	)

	rootCmd.AddCommand(
		loginCmd,
		logoutCmd,
		whoamiCmd,
		// noun-first (new)
		activities.RootCmd,
		budget.RootCmd,
		workpackage.RootCmd,
		project.RootCmd,
		user.RootCmd,
		timeentry.RootCmd,
		wptype.RootCmd,
		status.RootCmd,
		notification.RootCmd,
		git.RootCmd,
	)

	rootCmd.InitDefaultCompletionCmd()
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "completion" {
			_ = cmd.InheritedFlags().MarkHidden("format")
			_ = cmd.InheritedFlags().MarkHidden("verbose")
		}
	}
}
