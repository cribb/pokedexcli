package main

type cliConfig struct {
	nextUrl     string
	previousUrl string
}

func main() {

	config := cliConfig{
		nextUrl:     "",
		previousUrl: "",
	}

	goREPL(&config)

}
