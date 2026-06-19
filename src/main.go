package main

import (
	"log"
	"mikctl/src/cmd"
	"mikctl/src/config"
	"mikctl/src/db"
)

func main() {

	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	conn, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	db.DB = conn

	defer conn.Close()

	cmd.Execute()

}
