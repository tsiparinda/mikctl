package cmd

import (
	"mikctl/src/db"
	"mikctl/src/router"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Collect from RouterOS and update information in DB",
	Long: `Collect RouterOS information about one or more routers and update t.devices and t.routers

			Examples:
			mikctl update 
			mikctl update -r mt20
			mikctl update -g CHR
			`,
	Args: cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {
		routers, err := db.GetRouters(false)
		if err != nil {
			return err
		}

		return router.UpdateRouters(
			routers,
		)
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
