package list

import (
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/types"
	"github.com/spf13/cobra"
)

var typesCmd = &cobra.Command{
	Use:   "types",
	Short: "Lists work package types",
	Long:  "Get a list of all work package types of the instance.",
	Run:   listTypes,
}

func listTypes(_ *cobra.Command, _ []string) {
	if all, err := types.All(); err == nil {
		printer.Types(all)
	} else {
		printer.Error(err)
	}
}
