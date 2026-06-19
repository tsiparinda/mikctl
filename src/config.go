package main

type Config struct {
	Verbose bool
	Workers int
	Timeout int

	DBPath string
}

