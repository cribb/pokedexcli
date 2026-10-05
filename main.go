package main

import (
	"github.com/cribb/pokedexcli/internal/pokecache"
	"time"
)

type cliConfig struct {
	nextUrl     string
	previousUrl string
	cache       *pokecache.Cache
}

func main() {

	config := cliConfig{
		nextUrl:     "",
		previousUrl: "",
		cache:       pokecache.NewCache(5 * time.Minute),
	}

	goREPL(&config)

}
