package main

import "fmt"

func commandInspect(cfg *Config, pokemon string, pokedex Pokedex) error {
	data, exists := pokedex.Pokedex[pokemon]
	if exists {
		showPokemonInfo(data)
	} else {
		fmt.Println("you have not caught that pokemon")
	}

	return nil
}

func showPokemonInfo(pokemon Pokemon) {
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.StatInfo.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokemonType := range pokemon.PokemonTypes {
		fmt.Printf("  -%s\n", pokemonType.Type.Name)
	}
}