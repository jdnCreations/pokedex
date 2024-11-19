package main

import "os"


func commandExit(cfg *Config,_ string, _ Pokedex) error {
	os.Exit(0)
	return nil
}