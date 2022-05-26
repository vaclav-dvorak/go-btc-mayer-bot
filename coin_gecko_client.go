package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	days       = 200
	timeoutSec = 5
)

type avgResponse struct {
	Prices [][]float64 `json:"prices"`
}

type rateResponse struct {
	Rate map[string]float64 `json:"tether"`
}

type priceResponse struct {
	Price map[string]float64 `json:"bitcoin"`
}

func getAvgPrice(cur string) (avg float64, err error) {
	var (
		data avgResponse
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*timeoutSec))
	defer cancel()
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
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	defer func() {
		_ = res.Body.Close()
	}()
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}
	sum := 0.0
	for i := 0; i < len(data.Prices); i++ {
		sum += data.Prices[i][1]
	}
	avg = sum / float64(days)
	log.Infof("%d day moving average: %s\n", days, fmtPrice(avg, cur))
	return
}

func getConversionRate(sourceCur, targetCur string) (rate float64, err error) {
	var (
		data rateResponse
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*timeoutSec))
	defer cancel()
	header := http.Header{
		"Accept":       []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.coingecko.com/api/v3/simple/price", nil)
	req.Header = header
	q := req.URL.Query()
	q.Add("ids", "tether")
	q.Add("vs_currencies", strings.Join([]string{sourceCur, targetCur}, ","))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	defer func() {
		_ = res.Body.Close()
	}()
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}
	rate = data.Rate[targetCur] / data.Rate[sourceCur]
	log.Infof("exchange rate for %s%s%s: %s%.4f%s\n", green, targetCur, reset, blue, rate, reset)
	return
}

func getCurrentPrice(cur string) (price float64, err error) {
	var (
		data priceResponse
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*timeoutSec))
	defer cancel()
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
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	defer func() {
		_ = res.Body.Close()
	}()
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}
	price = data.Price[cur]
	log.Infof("current price of bitcoin: %s\n", fmtPrice(price, cur))
	return
}
