package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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