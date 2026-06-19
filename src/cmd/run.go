package cmd

import (
	"mikctl/src/config"
	"mikctl/src/db"
	"mikctl/src/script"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <script>",
	Short: "Run RouterOS script",
	Long: `Run RouterOS configuration script line by line to one or more routers.

			Examples:
			mikctl run backup.txt
			mikctl run backup.txt -r mt20
			mikctl run backup.txt -g CHR
			`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {

		scriptName := args[0]

		routers, err := db.GetRouters(false)
		if err != nil {
			return err
		}

		return script.RunScript(
			routers,
			scriptName,
			config.Workers,
		)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
