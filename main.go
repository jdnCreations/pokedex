package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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

type Pokedex struct {
  Pokedex map[string]Pokemon
}

type StatInfo struct {
  Name string `json:"name"`
}

type Stat struct {
  BaseStat string `json:"base_stat"`
  StatInfo StatInfo `json:"stat"`
}

type Type struct {
  Name string `json:"name"`
}

type Pokemon struct {
  Name string `json:"name"`
  BaseExperience int `json:"base_experience"`
  Height int `json:"height"`
  Weight int `json:"weight"`
  // Stats []Stat
  // Types []Type
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

type PokemonNameURL struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

type PokemonEncounter struct {
	PokemonNameURL PokemonNameURL `json:"pokemon"`
	_ []any
}

type SingleAreaResponse struct {
  ID int `json:"id"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

var cache *pokecache.Cache
var pokedex = Pokedex{
  Pokedex: make(map[string]Pokemon),
}



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

	data, found := cache.Get(url+area)
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
    fmt.Printf(" - %s\n", encounter.PokemonNameURL.Name)
  }

	return nil
}

func commandViewPokedex(_ *Config, _ string) error {
  viewPokedex()
  return nil
}

func commandCatch(cfg *Config, pokemonName string) error {
  if pokemonName == "" {
    return fmt.Errorf("you didn't provide a pokemon")
  }

  var pokemon Pokemon 
  
  data, found := cache.Get(pokemonName)

  url := "https://pokeapi.co/api/v2/pokemon/"
  if !found {
    res, err := http.Get(url+pokemonName)
    if res.StatusCode == 404 {
			return fmt.Errorf("no such pokemon exists")
		}
    if err != nil {
      return err
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
      return err
    }

    cache.Add(url+pokemonName, body)
    err = json.Unmarshal(body, &pokemon)
    if err != nil {
      return err
    }
  } else {
    err := json.Unmarshal(data, &pokemon)
    if err != nil {
      return err
    }

  }
  caught := attemptCatch(pokemon) 
  if caught {
    fmt.Printf("%s was caught!\n", pokemon.Name)
    pokedex.Pokedex[pokemon.Name] = pokemon
  } else {
    fmt.Printf("%s escaped!\n", pokemon.Name)
  }
  
  return nil
}

func viewPokedex() {
  if len(pokedex.Pokedex) > 0 {
    for _, pokemon := range pokedex.Pokedex {
      fmt.Println(pokemon.Name)
    }
  }

}

func attemptCatch(pokeInfo Pokemon) bool {
  fmt.Printf("Throwing a Pokeball at %s...\n", pokeInfo.Name) 
  src := rand.NewSource(time.Now().UnixNano())
  r := rand.New(src)
  chance := getChance(pokeInfo.BaseExperience)
  if r.Float64() < chance {
    return true
  } else {
    return false
  }
}

func getChance(baseExperience int) float64 {
  switch {
  case baseExperience < 5:
    return 0.95
  case baseExperience < 10:
    return 0.90
  case baseExperience < 20:
    return 0.80
  case baseExperience < 40:
    return 0.60 
  case baseExperience < 60:
    return 0.40
  case baseExperience < 80:
    return 0.20
  case baseExperience < 90:
    return 0.10
  case baseExperience < 95:
    return 0.05
  case baseExperience < 100:
    return 0.04
  case baseExperience < 105:
    return 0.03
  case baseExperience < 110:
    return 0.02
  default:
    return 0.01
  }
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
    "catch": {
      name: "catch",
      description: "Attempt to catch a pokemon",
      callback: commandCatch,
    },
    "view": {
      name: "view",
      description: "View your pokemon",
      callback: commandViewPokedex,
    },
	}
}