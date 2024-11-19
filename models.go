package main

type cliCommand struct {
	name string
	description string
	callback func(*Config, string, Pokedex) error
}

type Pokedex struct {
  Pokedex map[string]Pokemon
}

type StatInfo struct {
  Name string `json:"name"`
}

type Stat struct {
  BaseStat int `json:"base_stat"`
  StatInfo StatInfo `json:"stat"`
}

type PokemonType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}


type Pokemon struct {
  Name string `json:"name"`
  BaseExperience int `json:"base_experience"`
  Height int `json:"height"`
  Weight int `json:"weight"`
  Stats []Stat `json:"stats"`
  PokemonTypes []PokemonType `json:"types"`
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