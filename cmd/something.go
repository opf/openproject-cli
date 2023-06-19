package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:        "test",
	Short:      "prints some test output",
	Long: `This is a long description text of
how the test command's description should be.`,
	Example:                    "op test",
	Args:                       cobra.MinimumNArgs(0),
	Run:                        test,
}

func test(cmd *cobra.Command, args []string) {
	fmt.Println("I am a test")
}
