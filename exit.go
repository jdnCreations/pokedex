package main

import "os"


func commandExit(cfg *Config,_ string) error {
	os.Exit(0)
	return nil
}