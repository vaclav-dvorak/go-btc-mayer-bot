package main

import "github.com/joho/godotenv"

func main() {
	_ = godotenv.Load()
	// telegramToken = os.Getenv("TELEGRAM_APITOKEN")

	printWelcome()
	_ = getAvgPrice()
}
