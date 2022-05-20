package main

import (
	"fmt"
)

var (
	logoSmall = []string{
		`       ▄ ▄   `,
		`      ▀█▀▀▀▀▄`,
		`       █▄▄▄▄▀`,
		`       █    █`,
		`      ▀▀█▀█▀ `,
	}
	version, date = "(devel)", "now"
)

func printWelcome() {
	for i := 0; i < len(logoSmall); i++ {
		fmt.Printf("%s%s%s", orange, logoSmall[i], reset)
		if i < (len(logoSmall) - 1) {
			fmt.Print("\n")
		}
	}
	fmt.Printf(" %s%s @%s%s // coinbase bot for mayer multiple driven orders\n", green, version, date, reset)
	fmt.Printf("all price data are powered by CoinGecko: %shttps://www.coingecko.com/%s\n\n", blue, reset)
}

func fmtPrice(price float64, cur string) (out string) {
	out = fmt.Sprintf("%s%.2f%s(%s%s%s)", blue, price, reset, green, cur, reset)
	return
}
