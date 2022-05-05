package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type config struct {
	Currency  string `koanf:"currency"`
	CCurrency string `koanf:"convert_currency"`
	Coinbase  struct {
		Key    string `koanf:"key"`
		Secret string `koanf:"secret"`
	} `koanf:"coinbase"`
	Orders []struct {
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
	_ = k.Load(env.Provider("MAYERBOT_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "MAYERBOT_")), "_", ".", -1)
	}), nil)
	if err := k.Unmarshal("", &conf); err != nil {
		log.Fatalf("error parsing config: %v", err)
	}
	if err := validateConfig(conf); err != nil {
		log.Fatalf("error validating config: %v", err)
	}

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

func validateConfig(conf config) error {
	if conf.Coinbase.Key == "" {
		return fmt.Errorf("value of COINBASE_KEY is not set")
	}
	if conf.Coinbase.Secret == "" {
		return fmt.Errorf("value of COINBASE_SECRET is not set")
	}
	vol := 0.0
	for _, order := range conf.Orders {
		vol += order.Volume
	}
	if vol > 1 {
		return fmt.Errorf("sum of all order volumes must be equal or less then 1. yours is %.2f", vol)
	}
	return nil
}
