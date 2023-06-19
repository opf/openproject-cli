package list

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Lists projects",
	Long: `Get a list of visible projects.
The list can get filtered further.`,
	Run: func(cmd *cobra.Command, args []string) {
		count := 0
		for count < 2 {
			fmt.Printf("Project #%d\n", count)
			count++
		}
	},
}
