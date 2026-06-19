package cmd

import (
	"mikctl/src/db"
	"mikctl/src/router"

	"github.com/spf13/cobra"
)

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Show routers count by filters",
	Args:  cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {

		routers, err := db.GetRouters(false)
		if err != nil {
			return err
		}

		return router.CountRouters(routers)
	},
}

func init() {
	RootCmd.AddCommand(countCmd)
}
