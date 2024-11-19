package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

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
