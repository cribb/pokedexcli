package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func goREPL() {
	// fmt.Println("Hello, World!")
	scanner := bufio.NewScanner(os.Stdin)

	// fmt.Println("Foobar, bitches!")
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		text := cleanInput(scanner.Text())
		fmt.Printf("Your command was: %v\n", text[0])
		// fmt.Printf("repl_input: %v", repl_input)
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}
