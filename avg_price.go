package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const days = 200

type response struct {
	Prices [][]float64 `json:"prices"`
}

func getAvgPrice() (avg float64) {
	var (
		data response
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
	fmt.Printf("%d day moving avarage: %s%.2f%s(%s)\n", days, blue, avg, reset, cur)
	return
}
