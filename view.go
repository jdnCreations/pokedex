package main

import "fmt"

func commandPokedex(_ *Config, _ string, pokedex Pokedex) error {
  viewPokedex(pokedex)
  return nil
}


func viewPokedex(pokedex Pokedex) {
  if len(pokedex.Pokedex) > 0 {
		fmt.Println("Your Pokedex:")
    for _, pokemon := range pokedex.Pokedex {
      fmt.Printf(" - %s\n", pokemon.Name)
    }
  }
}