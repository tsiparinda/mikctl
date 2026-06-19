package cmd

import (
	"mikctl/src/config"
	"mikctl/src/db"
	"mikctl/src/importrsc"

	"github.com/spf13/cobra"
)

var importrscCmd = &cobra.Command{
	Use:   "import <script.rsc>",
	Short: "Import RouterOS script",
	Long: `Import RouterOS configuration script like .*rsc with multiline command to one or more routers.

			Examples:
			mikctl import backup.rsc
			mikctl import backup.rsc -r mt20
			mikctl import backup.rsc -g CHR
			`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {

		scriptName := args[0]

		routers, err := db.GetRouters(false)
		if err != nil {
			return err
		}

		return importrsc.ImportRsc(
			routers,
			scriptName,
			config.Workers,
		)
	},
}

func init() {
	RootCmd.AddCommand(importrscCmd)
}
