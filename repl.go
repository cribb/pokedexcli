package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io" // for the http handling
	"net/http"
	"os"
	"strings"
)

// Type definitions occur at "package level" (global scope)
type cliCommand struct {
	name        string
	description string
	callback    func(*cliConfig) error
}

func goREPL(config *cliConfig) {

	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()

	// fmt.Println("Foobar, bitches!")
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		commandLine := cleanInput(scanner.Text())
		commandName := commandLine[0]

		command, ok := commands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(config)

		if err != nil {
			fmt.Printf("Error '%v' occurred.\n", err)
		}

	}

	// return nil
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "display next location areas in Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "display previous location areas in Pokemon world",
			callback:    commandMapb,
		},
	}
	return commands
}

// func getConfig() cliConfig {

// 	return config
// }

func commandMap(config *cliConfig) error {

	// fmt.Println()

	if config.nextUrl == "" {
		config.nextUrl = "https://pokeapi.co/api/v2/location-area/"
	}

	// fmt.Printf(" --> config: %v\n", config)

	var location_area LocationAreaList

	res, err := http.Get(config.nextUrl)
	if err != nil {
		fmt.Printf("Error in GET functionality.\n")
		// return nil, fmt.Errorf("error creating request: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	// fmt.Printf(" --> body: %v, err: %v\n", res.Body, err)
	defer res.Body.Close()

	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		// log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		return nil //, err
	}

	if err := json.Unmarshal(body, &location_area); err != nil {
		return err
	}

	for _, result := range location_area.Results {
		fmt.Printf("%v\n", result.Name)
	}

	if location_area.Next != nil {
		config.nextUrl = *location_area.Next
	} else {
		config.nextUrl = ""
	}

	if location_area.Previous != nil {
		config.previousUrl = *location_area.Previous
	} else {
		config.previousUrl = ""
	}

	return nil
}

func commandMapb(config *cliConfig) error {

	// fmt.Println()

	if config.previousUrl == "" {
		fmt.Printf("you're on the first page\n")
		return nil
	}

	// fmt.Printf(" --> config: %v\n", config)

	var location_area LocationAreaList

	res, err := http.Get(config.previousUrl)
	if err != nil {
		fmt.Printf("Error in GET functionality.\n")
		// return nil, fmt.Errorf("error creating request: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	// fmt.Printf(" --> body: %v, err: %v\n", res.Body, err)
	defer res.Body.Close()

	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		// log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		return nil //, err
	}

	if err := json.Unmarshal(body, &location_area); err != nil {
		return err
	}

	for _, result := range location_area.Results {
		fmt.Printf("%v\n", result.Name)
	}

	if location_area.Next != nil {
		config.nextUrl = *location_area.Next
	} else {
		config.nextUrl = ""
	}

	if location_area.Previous != nil {
		config.nextUrl = *location_area.Previous
	} else {
		config.previousUrl = ""
	}

	return nil
}

func commandHelp(config *cliConfig) error {
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	// for INDEX, ELEMENT := range SLICE {}
	commands := getCommands()
	for _, command := range commands {
		fmt.Printf("%v\t%v\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandExit(config *cliConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}
