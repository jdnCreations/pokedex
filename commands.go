package main

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
		"inspect": {
			name: "inspect",
			description: "Inspect a Pokemon",
			callback: commandInspect,
		},
    "pokedex": {
      name: "pokedex",
      description: "View your pokemon",
      callback: commandPokedex,
    },
	}
}