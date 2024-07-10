package git

import (
	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/cmd/git/start"
)

var RootCmd = &cobra.Command{
	Use:   "git [command]",
	Short: "Executes a git related command.",
	Long: `This is a command prefix that enables a nested
set of executable commands, related to a git context.`,
}

func init() {
	RootCmd.AddCommand(start.RootCmd)
}
