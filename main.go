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

type Config struct {
	Pokedex pokedex.PokedexConfig
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
			Next: pokedex.BaseUrl + "/location-area",
		},
	}
	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		input := scanner.Text()
		command := cleanInput(input)

		supportedCommand, ok := supportedCommands[command[0]]

		if ok {
			err := supportedCommand.callback(&config)
			if err != nil {
				fmt.Print(err)
			}
		} else {
			fmt.Print("Unknown command.")
		}
	}
}

func commandMap(config *Config) error {
	res, err := pokedex.Get(config.Pokedex.Next)
	if err != nil {
		return fmt.Errorf("Error getting locations: %s", err)
	}

	for _, location := range res.Results {
		fmt.Print(location.Name + "\n")
	}

	config.Pokedex.Next = res.Next
	config.Pokedex.Previous = res.Previous

	return nil
}

func commandMapBack(config *Config) error {
	res, err := pokedex.Get(config.Pokedex.Previous)
	if err != nil {
		return fmt.Errorf("Error getting locations: %s", err)
	}

	for _, location := range res.Results {
		fmt.Print(location.Name + "\n")
	}

	config.Pokedex.Next = res.Next
	config.Pokedex.Previous = res.Previous

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
