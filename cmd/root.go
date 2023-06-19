package cmd

import (
	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/cmd/list"
	"github.com/opf/openproject-cli/components/requests"
)

var rootCmd = &cobra.Command{
	Use:   "op",
	Short: "An easy-to-use CLI for the OpenProject APIv3",
	Long: `OpenProject CLI is a fast, reliable and easy-to-use
tool to manage your work packages, notifications and
projects of your OpenProject instance.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	requests.Init("https://openproject.local")

	rootCmd.AddCommand(testCmd, list.RootCmd)
}
