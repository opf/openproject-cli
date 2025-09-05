package cmd

import (
	"fmt"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/cmd/create"
	"github.com/opf/openproject-cli/cmd/git"
	"github.com/opf/openproject-cli/cmd/inspect"
	"github.com/opf/openproject-cli/cmd/list"
	"github.com/opf/openproject-cli/cmd/search"
	"github.com/opf/openproject-cli/cmd/update"
	"github.com/opf/openproject-cli/components/configuration"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/routes"
)

var Verbose bool
var showVersionFlag bool

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "An easy-to-use CLI for the OpenProject APIv3",
	Long: `OpenProject CLI is a fast, reliable and easy-to-use
tool to manage your work packages, notifications and
projects of your OpenProject instance.`,
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

			fmt.Println(printer.Yellow(versionText))

			return
		}

		cmd.Help()
	},
}

func Execute(version *configuration.Version) error {
	configuration.Init(version)

	return rootCmd.Execute()
}

func init() {
	activePrinter := &printer.ConsolePrinter{}
	printer.Init(activePrinter)

	host, token, err := configuration.ReadConfig()
	if err != nil {
		printer.Error(err)
		return
	}

	parse, err := url.Parse(host)
	if err != nil {
		printer.Error(err)
	}

	requests.Init(parse, token, Verbose)
	routes.Init(parse)

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

	rootCmd.AddCommand(
		loginCmd,
		list.RootCmd,
		update.RootCmd,
		inspect.RootCmd,
		create.RootCmd,
		search.RootCmd,
		git.RootCmd,
	)
}
