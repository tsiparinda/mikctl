package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	SSHUser    string
	PasswordID int
	Workers    int
	Verbose    bool
	ROS        string
	Site       string
	Router     string
	Group      string
	Main       int
	DBPath     string
)

func InitConfig() error {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/mikctl")

	viper.SetEnvPrefix("MIKCTL")
	viper.AutomaticEnv()

	viper.SetDefault("ssh_user", "admin")
	viper.SetDefault("password_id", 0)
	viper.SetDefault("workers", 10)
	viper.SetDefault("db_path", "./mikrotik.db")
	_ = viper.ReadInConfig()

	SSHUser = viper.GetString("ssh_user")
	PasswordID = viper.GetInt("password_id")
	Workers = viper.GetInt("workers")
	DBPath = viper.GetString("db_path")

	return nil
}

func Dump() {
	fmt.Printf(
		"ssh_user=%s password_id=%d workers=%d\n",
		SSHUser,
		PasswordID,
		Workers,
	)
}

func Verbosef(format string, args ...any) {
	if Verbose {
		fmt.Printf(format, args...)
	}
}
