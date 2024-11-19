package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func commandMap(cfg *Config, _ string, _ Pokedex) error {
	
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