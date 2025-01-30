package cmd

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"pokedex/api"
	"pokedex/cache"
	"pokedex/utils"

	"github.com/fatih/color"
)

type CommandArgs struct {
	Cache          *cache.Cache
	Config         *Config
	Argument       *string
	CaughtPokemon  *map[string]utils.Pokemon
	VisiblePokemon map[string]struct{}
}

type Command struct {
	Name        string
	Description string
	CallBack    func(*CommandArgs) error
	Args        int
}

type Config struct {
	Next *string
	Prev *string
}

func CaughtPokis() map[string]utils.Pokemon {
	return make(map[string]utils.Pokemon)
}

func Commands() map[string]Command {
	return map[string]Command{
		"exit":    {"exit", "Exit the pokedex", exit, 1},
		"help":    {"Help", "Display help", help, 1},
		"map":     {"map", "Displays the next 20 location areas", mapCmd, 1},
		"mapb":    {"mapb", "Displays the previous 20 location areas", mapBackCmd, 1},
		"explore": {"explore <city name>", "Explore pokemons in a location", exploreCmd, 2},
		"catch":   {"catch <pokemon name>", "Attempt to catch pokemon", catchCmd, 2},
		"inspect": {"inspect <pokemon name>", "Show details about pokemon", inspectCmd, 2},
		"pokedex": {"pokedex", "Display all your caught pokemons", pokedexCmd, 1},
	}
}

func pokedexCmd(args *CommandArgs) error {
	utils.PrintColor("Pokemons in your pokedex", utils.Magenta)
	for k := range *args.CaughtPokemon {
		utils.PrintColor("#  "+k, utils.Green)
	}

	return nil
}

func inspectCmd(args *CommandArgs) error {
	if p, ok := (*args.CaughtPokemon)[*args.Argument]; !ok {
		utils.PrintColor(fmt.Sprintf("%s is not in your pokedex, find him and capture to see that stats", *args.Argument), utils.Red)
	} else {
		blue := color.New(color.FgBlue).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Printf("%s: %s\n", blue("Name"), green(p.Name))
		fmt.Printf("%s: %s\n", blue("Height"), green(fmt.Sprintf("%d", p.Height)))
		fmt.Printf("%s: %s\n", blue("Weight"), green(fmt.Sprintf("%d", p.Weight)))

		fmt.Println(blue("Stats:"))
		for _, stat := range p.Stats {
			fmt.Printf("  - %s: %s\n", blue(stat.Stat.Name), green(fmt.Sprintf("%d", stat.BaseStat)))
		}

		fmt.Println(blue("Types:"))
		for _, t := range p.Types {
			fmt.Printf("  - %s\n", green(t.Type.Name))
		}
	}
	return nil
}

func catchHelper(res utils.Pokemon) bool {
	xp := res.BaseExperience

	chance := 100 - xp/2
	if rand.IntN(101) > chance {
		return true
	}

	return false
}

func catchCmd(args *CommandArgs) error {
	if _, ok := (*args.CaughtPokemon)[*args.Argument]; ok {
		return fmt.Errorf("%s is allready caught and in your pokedex", *args.Argument)
	}
	if args.VisiblePokemon == nil {
		return fmt.Errorf("There are no pokemon here")
	}

	if _, ok := args.VisiblePokemon[*args.Argument]; !ok {
		return fmt.Errorf("%s is not in current location", *args.Argument)
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + *args.Argument
	response, err := api.FetchPokemon(&url)
	if err != nil {
		return fmt.Errorf("Failed to catch pokemon")
	}

	caught := catchHelper(*response)

	utils.PrintColor(fmt.Sprintf("Throwing pokeball at %s...\n", response.Name), utils.Magenta)
	if caught {
		(*args.CaughtPokemon)[response.Name] = *response
		utils.PrintColor(fmt.Sprintf("%s is captured\n", response.Name), utils.Green)
	} else {
		utils.PrintColor(fmt.Sprintf("%s escaped\n", response.Name), utils.Red)
	}
	return nil
}

func exploreCmd(args *CommandArgs) error {
	url := "https://pokeapi.co/api/v2/location-area/" + *args.Argument
	response, err := api.FetchLocationAreas(&url, args.Cache)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %v", err)
	}

	pokis := make(map[string]struct{})
	utils.PrintColor(fmt.Sprintf("Pokemons in %s", *args.Argument), utils.Magenta)
	for _, location := range response.PokemonEncounters {
		if _, ok := (*args.CaughtPokemon)[location.Pokemon.Name]; !ok {
			pokis[location.Pokemon.Name] = struct{}{}
			color.Green("#  " + location.Pokemon.Name)
		}
	}

	args.VisiblePokemon = pokis
	return nil
}

func mapCmd(args *CommandArgs) error {
	args.VisiblePokemon = nil
	if args.Config.Next == nil {
		url := "https://pokeapi.co/api/v2/location-area/"
		args.Config.Next = &url
	}

	response, err := api.FetchLocationAreas(args.Config.Next, args.Cache)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %v", err)
	}

	color.Magenta("Cities:")
	for _, location := range response.Results {
		utils.PrintColor("#  "+location.Name, utils.Green)
	}

	if response.Next != nil {
		args.Config.Next = response.Next
	} else {
		args.Config.Next = nil
	}

	if response.Prev != nil {
		args.Config.Prev = response.Prev
	} else {
		args.Config.Prev = nil
	}

	return nil
}

func mapBackCmd(args *CommandArgs) error {
	args.VisiblePokemon = nil
	if args.Config.Prev == nil {
		utils.PrintColor("You're on the first page.", utils.Blue)

		return nil
	}

	response, err := api.FetchLocationAreas(args.Config.Prev, args.Cache)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %v", err)
	}

	color.Magenta("Cities:")
	for _, location := range response.Results {
		utils.PrintColor("#  "+location.Name, utils.Green)
	}

	if response.Next != nil {
		args.Config.Next = response.Next
	} else {
		args.Config.Next = nil
	}

	if response.Prev != nil {
		args.Config.Prev = response.Prev
	} else {
		args.Config.Prev = nil
	}

	return nil
}

func exit(args *CommandArgs) error {
	utils.PrintColor("Exiting pokedex, bye bye!", utils.Blue)
	os.Exit(0)
	return errors.New("Impossible error")
}

func help(args *CommandArgs) error {
	utils.PrintColor("Welcome to pokedex help", utils.Blue)
	utils.PrintColor("Available commands", utils.Magenta)
	for _, v := range Commands() {
		utils.PrintColor(v.Name+": "+v.Description, utils.Green)
	}
	return nil
}
