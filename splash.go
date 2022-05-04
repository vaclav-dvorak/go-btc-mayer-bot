package main

import (
	"fmt"
)

var (
	logoSmall = []string{
		`     ▄ ▄   `,
		`    ▀█▀▀▀▀▄`,
		`     █▄▄▄▄▀`,
		`     █    █`,
		`    ▀▀█▀█▀ `,
	}
	version, date = "(devel)", "now"
)

func printWelcome() {
	for i, line := range logoSmall {
		fmt.Printf("%s%s%s", orange, line, reset)
		if i < (len(logoSmall) - 1) {
			fmt.Print("\n")
		}
	}
	fmt.Printf(" %s%s @%s%s // coinbase bot for mayer multiple driven orders\n\n", green, version, date, reset)
}
