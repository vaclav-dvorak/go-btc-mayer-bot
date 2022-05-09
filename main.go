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
	CalcCurrency string `koanf:"calculation_currency"`
	BuyCurrency  string `koanf:"buy_currency"`
	Coinbase     struct {
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
	avg := getAvgPrice(conf.CalcCurrency)

	price := getCurrentPrice(conf.CalcCurrency)
	curMayer := price / avg
	log.Printf("current price is at %s%.3f%s of mayer multiple", cyan, curMayer, reset)

	rate := 1.0
	if conf.BuyCurrency != "" {
		rate = getConversionRate(conf.CalcCurrency, conf.BuyCurrency)
	}

	for _, order := range conf.Orders {
		target := order.Mayer
		msg := ""
		if target > curMayer {
			msg += fmt.Sprintf("%s%.3f%s target is above current mayer multiple %s%.3f%s it would be best to buy spot price %s", cyan, target, reset, cyan, curMayer, reset, fmtPrice(price, conf.CalcCurrency))
		} else {
			msg += fmt.Sprintf("%s%.3f%s target = %s", cyan, target, reset, fmtPrice(avg*target, conf.CalcCurrency))
			if conf.BuyCurrency != "" {
				msg += fmt.Sprintf(" or converted %s\n", fmtPrice(avg*target*rate, conf.BuyCurrency))
			} else {
				msg += "\n"
			}
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
