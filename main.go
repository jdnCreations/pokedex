package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/jdnCreations/pokedex/internal/pokecache"
)

var cache *pokecache.Cache
var pokedex = Pokedex{
  Pokedex: make(map[string]Pokemon),
}

func main() {
	config := &Config{}
  cache = pokecache.InitializeCache(5*time.Minute)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("pokedex > ")
		scanner.Scan()
		handleInput(scanner.Text(), config)
	}
}
