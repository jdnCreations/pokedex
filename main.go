package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name string
	description string
	callback func(*Config) error
}

type LocationResult struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

type Config struct {
	Next *string `json:"next"`
	Previous *string `json:"previous"`
}

type LocationAreaResponse struct {
	Count int64 `json:"count"`
	Next *string `json:"next"`
	Previous *string `json:"previous"` 
	Results []LocationResult `json:"results"`
}



func main() {
	config := &Config{}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("pokedex > ")
		scanner.Scan()
		
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		if cmd, exists := getCommands()[commandName]; exists {
			if err := cmd.callback(config); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}

func commandHelp(cfg *Config) error {
	fmt.Println()
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage: \n\n")

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandExit(cfg *Config) error {
	os.Exit(0)
	return nil
}

func commandMap(cfg *Config) error {
	var loc LocationAreaResponse
	var url string

	if cfg.Next == nil {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = *cfg.Next
	}
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &loc)
	if err != nil {
		return err
	}

	for _, result := range loc.Results {
		fmt.Println(result.Name)
	}

	cfg.Next = loc.Next
	cfg.Previous = loc.Previous

	return nil

}

func commandMapb(cfg *Config) error {
	var loc LocationAreaResponse

	if cfg.Previous == nil {
		return fmt.Errorf("you're already on the first page")
	}

	res, err := http.Get(*cfg.Previous)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &loc)
	if err != nil {
		return err
	}

	for _, result := range loc.Results {
		fmt.Println(result.Name)
	}

	cfg.Next = loc.Next
	cfg.Previous = loc.Previous

	return nil

}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"map": {
			name: "map",
			description: "Displays names of 20 location areas in the Pokemon world",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Displays names of the 20 previous location areas in the Pokemon world",
			callback: commandMapb,
		},
	}
}