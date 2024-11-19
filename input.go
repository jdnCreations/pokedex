package main

import (
	"fmt"
	"strings"
)

func handleInput(input string, config *Config, pokedex Pokedex) {
  words := cleanInput(input)
		if len(words) == 0 {
			return
		}

		var commandParams string

		commandName := words[0]
		if len(words) > 1 {
			commandParams = words[1]
		}

		if cmd, exists := getCommands()[commandName]; exists {
			if err := cmd.callback(config, commandParams, pokedex); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
}


func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}