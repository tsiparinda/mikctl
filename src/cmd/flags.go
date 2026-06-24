package cmd

import (
	"fmt"
	"mikctl/src/config"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

var RootCmd = &cobra.Command{
	Use:   "mikctl",
	Short: "MikroTik management utility",
	Long: `mikctl is a utility for managing MikroTik routers.

			Features:
			- Run scripts
			- Import configurations
			- Discover neighbors
			- Ask RouterOS and update information in DB
		
Author: Oleg Tsyparynda
			`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	RootCmd.Version = Version
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	RootCmd.PersistentFlags().
		BoolVarP(
			&config.Verbose,
			"verbose",
			"v",
			false,
			"verbose output",
		)

	RootCmd.PersistentFlags().
		IntVarP(
			&config.Workers,
			"workers",
			"w",
			10,
			"parallel workers",
		)

	RootCmd.PersistentFlags().
		StringVarP(
			&config.ROS,
			"ROS",
			"R",
			"",
			"ROS major version (LIKE ... 6 or 7)",
		)

	RootCmd.PersistentFlags().
		StringVarP(
			&config.Site,
			"site",
			"s",
			"",
			"Site",
		)
	RootCmd.PersistentFlags().
		StringVarP(
			&config.Group,
			"group",
			"g",
			"",
			"Group",
		)

	RootCmd.PersistentFlags().
		StringVarP(
			&config.Router,
			"router",
			"r",
			"",
			"Name of router",
		)

	RootCmd.PersistentFlags().
		IntVarP(
			&config.Main,
			"main",
			"m",
			1,
			"1=main, 0=slave",
		)
}
