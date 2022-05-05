package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type response struct {
	Prices       [][]float64 `json:"prices"`
	MarketCaps   [][]float64 `json:"market_caps"`
	TotalVolumes [][]float64 `json:"total_volumes"`
}

func getAvgPrice() float64 {
	var (
		data response
		cur  = os.Getenv("CURRENCY")
	)

	ctx := context.Background()
	header := http.Header{
		"Accept":       []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.coingecko.com/api/v3/coins/bitcoin/market_chart", nil)
	req.Header = header
	q := req.URL.Query()
	q.Add("days", "200")
	q.Add("interval", "daily")
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
	count := 0
	for _, day := range data.Prices {
		sum += day[1]
		count++
	}
	var avg float64
	avg = sum / float64(count)
	fmt.Printf("200 day moving avarage: %s%.2f%s(%s)\n", blue, avg, reset, cur)
	return avg
}
