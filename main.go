package main

import (
	"time"

	"github.com/cribb/pokedexcli/internal/pokeapi"
	"github.com/cribb/pokedexcli/internal/pokecache"
)

type cliConfig struct {
	nextUrl     string
	previousUrl string
	cache       *pokecache.Cache
	pokedex     map[string]pokeapi.Pokemon
	userXp      int
}

func main() {

	config := cliConfig{
		nextUrl:     "",
		previousUrl: "",
		cache:       pokecache.NewCache(5 * time.Minute),
		pokedex:     make(map[string]pokeapi.Pokemon),
		userXp:      100,
	}

	goREPL(&config)

}
