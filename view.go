package main

import "fmt"

func commandViewPokedex(_ *Config, _ string) error {
  viewPokedex()
  return nil
}


func viewPokedex() {
  if len(pokedex.Pokedex) > 0 {
    for _, pokemon := range pokedex.Pokedex {
      fmt.Println(pokemon.Name)
    }
  }

}