package cmd

import (
	"mikctl/src/db"
	"mikctl/src/router"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show list of routers by filters",
	Args:  cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {

		routers, err := db.GetRouters(false)
		if err != nil {
			return err
		}

		return router.ListRouters(routers)
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
