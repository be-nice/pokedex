package utils

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintColor(s string, c Color) {
	switch c {
	case Green:
		s = color.GreenString(s)
	case Yellow:
		s = color.YellowString(s)
	case Magenta:
		s = color.MagentaString(s)
	case Blue:
		s = color.BlueString(s)
	case Red:
		s = color.RedString(s)
	}

	fmt.Println(s)
}

const Welcome = `
##############################
#                            #
#     Welcome to PoKeDeX     #
#                            #
#> Type "help" for commands <#
##############################
`
