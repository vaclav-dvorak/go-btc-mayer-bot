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
	Currency string `koanf:"currency"`
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
	// telegramToken = os.Getenv("TELEGRAM_APITOKEN")

	printWelcome()
	avg := getAvgPrice()
	fmt.Printf("0.8 * avg = %.2f in eur (%.2f)\n", avg*0.8, avg*0.8*0.94)
	fmt.Printf("0.79 * avg = %.2f in eur (%.2f)\n", avg*0.79, avg*0.79*0.94)
}
