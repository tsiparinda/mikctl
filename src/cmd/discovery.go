package cmd

import (
	"mikctl/src/db"
	"mikctl/src/router"

	"github.com/spf13/cobra"
)

var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "Discovery router's neighbors and store their in database",
	Args:  cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {

		routers, err := db.GetRouters(true)
		if err != nil {
			return err
		}

		return router.DiscoveryRouters(routers)
	},
}

func init() {
	RootCmd.AddCommand(discoveryCmd)
}
