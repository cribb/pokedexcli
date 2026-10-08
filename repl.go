package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io" // for the http handling
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/cribb/pokedexcli/internal/pokeapi"
)

// Type definitions occur at "package level" (global scope)
type cliCommand struct {
	name        string
	description string
	callback    func(*cliConfig, []string) error
}

const baseUrl = "https://pokeapi.co/api/v2/"

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
		if len(commandLine) == 0 {
			continue
		}

		commandName := commandLine[0]
		commandArgs := commandLine[1:]

		// fmt.Printf(" -- _|_ ->%v\n", commandArgs)
		command, ok := commands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(config, commandArgs)

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
		"quit": {
			name:        "quit",
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
		"explore": {
			name:        "explore",
			description: "display Pokemon found in the area",
			callback:    commandExplore,
		},
		"cache": {
			name:        "cache",
			description: "print cache",
			callback:    commandPrintCache,
		},
		"catch": {
			name:        "catch",
			description: "catch pokemon",
			callback:    commandCatch,
		},
		"print": {
			name:        "print",
			description: "print pokedex",
			callback:    commandPrint,
		},
	}
	return commands
}

func commandHelp(config *cliConfig, cmdArgs []string) error {
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	// for INDEX, ELEMENT := range SLICE {}
	commands := getCommands()
	for _, command := range commands {
		fmt.Printf("%v\t%v\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandExit(config *cliConfig, cmdArgs []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandPrintCache(config *cliConfig, cmdArgs []string) error {

	fmt.Println("Currently, the cache contains:")
	config.cache.Print()
	return nil
}

func commandMap(config *cliConfig, cmdArgs []string) error {

	var location_area_list pokeapi.LocationAreaList

	if config.nextUrl == "" {
		config.nextUrl = baseUrl + "location-area/"
	}

	body, err := fetchData(config, config.nextUrl)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &location_area_list); err != nil {
		return err
	}

	for _, result := range location_area_list.Results {
		fmt.Printf("%v\n", result.Name)
	}

	if location_area_list.Next != nil {
		config.nextUrl = *location_area_list.Next
	} else {
		config.nextUrl = ""
	}

	if location_area_list.Previous != nil {
		config.previousUrl = *location_area_list.Previous
	} else {
		config.previousUrl = ""
	}

	return nil
}

func commandMapb(config *cliConfig, cmdArgs []string) error {

	if config.previousUrl == "" {
		fmt.Printf("you're on the first page\n")
		return nil
	}

	var location_area_list pokeapi.LocationAreaList

	body, err := fetchData(config, config.previousUrl)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &location_area_list); err != nil {
		return err
	}

	for _, result := range location_area_list.Results {
		fmt.Printf("%v\n", result.Name)
	}

	if location_area_list.Next != nil {
		config.nextUrl = *location_area_list.Next
	} else {
		config.nextUrl = ""
	}

	if location_area_list.Previous != nil {
		config.nextUrl = *location_area_list.Previous
	} else {
		config.previousUrl = ""
	}

	return nil
}

func commandExplore(config *cliConfig, cmdArgs []string) error {

	area := cmdArgs[0]

	var location_area pokeapi.LocationArea

	url := baseUrl + "location-area/" + area + "/"
	fmt.Printf("Exploring the area: %v\n", area)
	fmt.Printf("Located at: %v\n", url)

	body, err := fetchData(config, url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &location_area); err != nil {
		return err
	}

	for _, encounter := range location_area.PokemonEncounters {
		fmt.Printf("%v\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandPrint(config *cliConfig, cmdArgs []string) error {
	fmt.Println(" === Pokedex contents ===")
	for _, pokemon := range config.pokedex {
		fmt.Printf("  %v\n", pokemon.Name)
	}
	return nil
}

func commandCatch(config *cliConfig, cmdArgs []string) error {

	if len(cmdArgs) == 0 {
		fmt.Println(" No pokemon specified")
		return nil
	}
	pokeTarget := cmdArgs[0]

	var pokemon pokeapi.Pokemon

	url := baseUrl + "pokemon/" + pokeTarget

	// fmt.Printf("Pulling info from %v...\n", url)

	body, err := fetchData(config, url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &pokemon); err != nil {
		return err
	}

	fmt.Printf(" Throwing a Pokeball at %v... \n", pokeTarget)
	throwPokeball(config, pokemon)

	return nil
}

func throwPokeball(config *cliConfig, pokemon pokeapi.Pokemon) {
	success := catchMagic(pokemon.BaseExperience, config.userXp)
	if success {
		config.pokedex[pokemon.Name] = pokemon
		fmt.Printf(" ==SUCCESS== You caught... %v! :)\n", pokemon.Name)
		return
	}
	fmt.Printf(" Oops-- You didn't catch %v. :(\n", pokemon.Name)

}

func catchMagic(baseExp int, userExp int) bool {

	chance := rand.Float64()

	// Goofy catch equation as first pass:
	// the ratio of user to Pokemon experience * random die throw
	score := float64(userExp) / float64(baseExp) * chance

	fmt.Printf(" Success score: %v\n", score)
	if score >= 1 {
		return true
	}
	return false
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))

	return words
}

func fetchData(config *cliConfig, url string) ([]byte, error) {

	var body []byte
	var ok bool
	var err error

	body, ok = config.cache.Get(url)
	if !ok {
		body, err = getDataFromNetwork(url)

		if err != nil {
			return nil, err
		}
		config.cache.Add(url, body)
	}
	return body, nil

}

func getDataFromNetwork(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error in GET functionality.\n")
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	// fmt.Printf(" --> body: %v, err: %v\n", res.Body, err)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		// fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return nil, fmt.Errorf("Response failed with status code: %d", res.StatusCode)
	}

	return body, nil
}
