package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jdnCreations/pokedex/internal/pokecache"
)

type cliCommand struct {
	name string
	description string
	callback func(*Config, string) error
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

type Pokemon struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

type PokemonEncounter struct {
	Pokemon Pokemon
	_ []any
}

type SingleAreaResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

var cache *pokecache.Cache



func main() {
	config := &Config{}
	interval := time.Duration(5 * time.Minute)
	cache, _ = pokecache.NewCache(interval)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("pokedex > ")
		scanner.Scan()
		
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		var commandParams string

		commandName := words[0]
		if len(words) > 1 {
			commandParams = words[1]
		}

		if cmd, exists := getCommands()[commandName]; exists {
			if err := cmd.callback(config, commandParams); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}

func commandHelp(cfg *Config, _ string) error {
	fmt.Println()
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage: \n\n")

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandExit(cfg *Config,_ string) error {
	os.Exit(0)
	return nil
}

func commandMap(cfg *Config, _ string) error {
	
	var loc LocationAreaResponse
	var url string

	if cfg.Next == nil {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = *cfg.Next
	}
	data, found := cache.Get(url)
	if !found {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cache.Add(url, body)
		err = json.Unmarshal(body, &loc)
		if err != nil {
			return err
		}
	} else {
		err := json.Unmarshal(data, &loc)
		if err != nil {
			return err
		}
	}

	for _, result := range loc.Results {
		fmt.Println(result.Name)
	}

	cfg.Next = loc.Next
	cfg.Previous = loc.Previous

	return nil

}

func commandMapb(cfg *Config, _ string) error {
	var loc LocationAreaResponse
	
	if cfg.Previous == nil {
		return fmt.Errorf("you're already on the first page")
	}

	data, found := cache.Get(*cfg.Previous)
	if !found {
		res, err := http.Get(*cfg.Previous)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cache.Add(*cfg.Previous, body)
		err = json.Unmarshal(body, &loc)
		if err != nil {
			return err
		}
	} else {
		err := json.Unmarshal(data, &loc)
		if err != nil {
			return err
		}
	}

	for _, result := range loc.Results {
		fmt.Println(result.Name)
	}

	cfg.Next = loc.Next
	cfg.Previous = loc.Previous

	return nil

}

func commandExplore(cfg *Config, area string) error {
	if area == "" {
		return fmt.Errorf("you did not provide an area")
	}

	var areaInfo SingleAreaResponse

	url := "https://pokeapi.co/api/v2/location-area/"

	data, found := cache.Get("area")
	if !found {
		res, err := http.Get(url + area)
		if res.StatusCode == 404 {
			return fmt.Errorf("no such area exists")
		}
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		
		cache.Add(url+area, body)
		err = json.Unmarshal(body, &areaInfo)
		if err != nil {
			return err
		}
	} else {
		err := json.Unmarshal(data, &areaInfo)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Exploring %s...\n", area)
	fmt.Println("Found Pokemon:")
	for _, encounter := range areaInfo.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}

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
		"explore": {
			name: "explore",
			description: "Displays names of pokemon in the area",
			callback: commandExplore,
		},
	}
}