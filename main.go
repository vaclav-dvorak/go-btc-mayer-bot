package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type config struct {
	Currency  string `koanf:"currency"`
	CCurrency string `koanf:"convert_currency"`
	Orders    []struct {
		Mayer  float64 `koanf:"mayer"`
		Volume float64 `koanf:"volume"`
	} `koanf:"orders"`
	CancelOrders bool `koanf:"cancel_current_orders"`
}

var (
	conf = config{}
	k    = koanf.New(".")
)

func main() {
	_ = godotenv.Load()
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	if err := k.Unmarshal("", &conf); err != nil {
		log.Fatalf("error parsing config: %v", err)
	}
	vol := 0.0
	for _, order := range conf.Orders {
		vol += order.Volume
	}
	if vol > 1 {
		log.Fatalf("sum of all order volumes must be equal or less then 1. yours is %.2f", vol)
	}
	// telegramToken = os.Getenv("TELEGRAM_APITOKEN")

	printWelcome()
	avg := getAvgPrice()
	rate := 1.0
	if conf.Currency == "usd" && conf.CCurrency != "" {
		rate = getConversionRate(conf.CCurrency)
	}

	for _, order := range conf.Orders {
		target := order.Mayer
		msg := fmt.Sprintf("%.2f target = %s%.2f%s(%s%s%s)", target, blue, avg*target, reset, green, conf.Currency, reset)
		if rate != 1.0 {
			msg += fmt.Sprintf(" or converted %s%.2f%s(%s%s%s)\n", blue, avg*target*rate, reset, green, conf.CCurrency, reset)
		} else {
			msg += "\n"
		}
		log.Print(msg)
	}
}
