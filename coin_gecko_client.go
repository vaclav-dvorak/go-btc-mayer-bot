package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const days = 200

type avgResponse struct {
	Prices [][]float64 `json:"prices"`
}

type rateResponse struct {
	Rate map[string]float64 `json:"tether"`
}

type priceResponse struct {
	Price map[string]float64 `json:"bitcoin"`
}

func getAvgPrice() (avg float64) {
	var (
		data avgResponse
		cur  = conf.Currency
	)

	ctx := context.Background()
	header := http.Header{
		"Accept":       []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.coingecko.com/api/v3/coins/bitcoin/market_chart", nil)
	req.Header = header
	q := req.URL.Query()
	q.Add("days", strconv.Itoa(days-1))
	q.Add("vs_currency", cur)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Print(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		log.Print(err)
	}
	_ = res.Body.Close()
	sum := 0.0
	for _, day := range data.Prices {
		sum += day[1]
	}
	avg = sum / float64(days)
	log.Printf("%d day moving avarage: %s%.2f%s(%s%s%s)\n", days, blue, avg, reset, green, cur, reset)
	return
}

func getConversionRate(cur string) (rate float64) {
	var (
		data rateResponse
	)

	ctx := context.Background()
	header := http.Header{
		"Accept":       []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.coingecko.com/api/v3/simple/price", nil)
	req.Header = header
	q := req.URL.Query()
	q.Add("ids", "tether")
	q.Add("vs_currencies", cur)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Print(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		log.Print(err)
	}
	_ = res.Body.Close()
	rate = data.Rate[cur]
	log.Printf("exchange rate for %s%s%s: %s%.4f%s\n", green, cur, reset, blue, rate, reset)
	return
}

func getCurrentPrice(cur string) (price float64) {
	var (
		data priceResponse
	)

	ctx := context.Background()
	header := http.Header{
		"Accept":       []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.coingecko.com/api/v3/simple/price", nil)
	req.Header = header
	q := req.URL.Query()
	q.Add("ids", "bitcoin")
	q.Add("vs_currencies", cur)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Print(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		log.Print(err)
	}
	_ = res.Body.Close()
	price = data.Price[cur]
	log.Printf("current price of bitcoin: %s%.2f%s(%s%s%s)\n", blue, price, reset, green, cur, reset)
	return
}
