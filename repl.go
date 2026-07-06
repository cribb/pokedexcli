package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Type definitions occur at "package level" (global scope)
type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func goREPL() {

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

		err := command.callback()

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
	}
	return commands
}

func commandHelp() error {
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	// for INDEX, ELEMENT := range SLICE {}
	commands := getCommands()
	for _, command := range commands {
		fmt.Printf("%v\t%v\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}
