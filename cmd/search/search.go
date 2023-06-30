package search

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "search [resource]",
	Short: "Searches for the specific resource",
	Long: `Execute a search on the given resource, with an search
input text. This input can be an id, a key word
like 'me', or a name.`,
}

func init() {
	RootCmd.AddCommand(userCmd)
}
