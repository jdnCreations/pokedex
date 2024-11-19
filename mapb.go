package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func commandMapb(cfg *Config, _ string, _ Pokedex) error {
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