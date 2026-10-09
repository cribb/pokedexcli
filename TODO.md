# TODO Pokedex

- refactor repl.go, pulling out "commands" into separate file (commands.go?)
- check on support/utility functions (util.go vs utils embedded in needed commands vs both?)
- implement the up arrow -- bufio.Scanner is insufficient for this. check github.com/chzyer/readline or github.com/manifoldco/promptui for alternatives
- add more (and better) unit testing
- maybe tie this into go-server course and serve current html state of pokedex (using graphics avail by api, etc)
- add 'save pokedex' feature

