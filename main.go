package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedex/cache"
	"pokedex/cmd"
	"pokedex/utils"
	"time"

	"github.com/fatih/color"
)

func main() {
	utils.PrintColor(utils.Welcome, utils.Blue)
	cache := cache.NewCache(30 * time.Second)
	defer cache.Stop()

	caughtPokis := cmd.CaughtPokis()
	commandMap := cmd.Commands()

	args := &cmd.CommandArgs{Cache: cache, Config: &cmd.Config{}, CaughtPokemon: &caughtPokis}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(color.YellowString("Pokedex > "))

	for scanner.Scan() {
		input := utils.CleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		if cmd, ok := commandMap[input[0]]; ok {
			if len(input) != cmd.Args {
				color.Red("Incorrect amount of arguments! EXPECTED: %d, GOT: %d\n", cmd.Args, len(input))
				fmt.Print("Pokedex > ")
				continue
			}
			if len(input) > 1 {
				args.Argument = &input[1]
			}

			if err := cmd.CallBack(args); err != nil {
				utils.PrintColor(fmt.Sprint(err), utils.Red)
			}
		} else {
			utils.PrintColor(fmt.Sprintf("Unknown command, skipping: %s\n", scanner.Text()), utils.Red)
		}

		fmt.Print(color.YellowString("Pokedex > "))
	}

	fmt.Println("Done")
}
