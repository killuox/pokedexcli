package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/killuox/pokedexcli/internal/pokedex"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type ConfigParams struct {
	id string
}

type Config struct {
	Pokedex pokedex.PokedexConfig
	Params  ConfigParams
}

var supportedCommands map[string]cliCommand

func init() {
	supportedCommands = map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Display the next locations area of the pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map back",
			description: "Display the previous locations area of the pokemon world",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "See a list of all the pokemon from a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Thow a pokeball to a desired pokemon",
			callback:    commandCatch,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := Config{
		Pokedex: pokedex.PokedexConfig{
			NextLocation: pokedex.BaseUrl + "/location-area",
		},
		Params: ConfigParams{},
	}
	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		input := scanner.Text()
		cleanedInput := cleanInput(input)

		command := cleanedInput[0]

		if len(cleanedInput) > 1 {
			config.Params.id = cleanedInput[1]
		}

		supportedCommand, ok := supportedCommands[command]
		if ok {
			err := supportedCommand.callback(&config)
			if err != nil {
				fmt.Print(err)
			}
		} else {
			fmt.Print("Unknown command.")
		}
		fmt.Print("Pokedex > ")
	}
}

func commandMap(config *Config) error {
	res, err := pokedex.GetLocations(config.Pokedex.NextLocation)
	if err != nil {
		return fmt.Errorf("Error getting locations: %s", err)
	}

	for _, location := range res.Results {
		fmt.Print(location.Name + "\n")
	}

	config.Pokedex.NextLocation = res.Next
	config.Pokedex.PreviousLocation = res.Previous

	return nil
}

func commandMapBack(config *Config) error {
	res, err := pokedex.GetLocations(config.Pokedex.PreviousLocation)
	if err != nil {
		return fmt.Errorf("Error getting locations: %s", err)
	}

	for _, location := range res.Results {
		fmt.Print(location.Name + "\n")
	}

	config.Pokedex.NextLocation = res.Next
	config.Pokedex.PreviousLocation = res.Previous

	return nil
}

func commandExplore(config *Config) error {
	id := config.Params.id
	if id == "" {
		return fmt.Errorf("Please provide the location name or id to explore")
	}
	fmt.Printf("Exploring %s...\n", id)
	res, err := pokedex.GetLocation(id)
	if err != nil {
		return fmt.Errorf("Error getting location are %s", err)
	}

	fmt.Print("Found Pokemon:\n")

	for _, pokemon := range res.PokemonEncounters {
		fmt.Printf("- %s\n", pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *Config) error {
	id := config.Params.id
	if id == "" {
		return fmt.Errorf("Please provide the location name or id to explore")
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", id)
	res, err := pokedex.GetPokemon(id)
	if err != nil {
		return fmt.Errorf("Error getting location are %s", err)
	}

	return nil
}

func commandHelp(config *Config) error {
	message := "Welcome to the Pokedex!\nUsage:\n\n"

	for _, command := range supportedCommands {
		message += fmt.Sprintf("%s: %s\n", command.name, command.description)
	}

	fmt.Print(message)

	return nil
}

func commandExit(config *Config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
