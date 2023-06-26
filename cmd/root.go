package cmd

import (
	"github.com/opf/openproject-cli/cmd/create"
	"net/url"
	"os"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/cmd/inspect"
	"github.com/opf/openproject-cli/cmd/list"
	"github.com/opf/openproject-cli/cmd/update"
	"github.com/opf/openproject-cli/components/configuration"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/routes"
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "An easy-to-use CLI for the OpenProject APIv3",
	Long: `OpenProject CLI is a fast, reliable and easy-to-use
tool to manage your work packages, notifications and
projects of your OpenProject instance.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	activePrinter := &printer.ConsolePrinter{}
	printer.Init(activePrinter)

	host, token, err := configuration.ReadConfigFile()
	if err != nil {
		printer.Error(err)
	}

	parse, err := url.Parse(host)
	if err != nil {
		printer.Error(err)
	}

	requests.Init(parse, token)
	routes.Init(parse)

	rootCmd.AddCommand(
		loginCmd,
		list.RootCmd,
		update.RootCmd,
		inspect.RootCmd,
		create.RootCmd,
	)
}
