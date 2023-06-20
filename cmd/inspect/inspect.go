package inspect

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "inspect [type] [id]",
	Short: "Show details about an object",
	Long:  "Show detailed information of an object of a specific type refereced by it's ID.",
}

func init() {
	RootCmd.AddCommand(inspectProjectCmd)
}
