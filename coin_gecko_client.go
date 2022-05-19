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

func getAvgPrice(cur string) (float64, error) {
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
		return 0, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}
	_ = res.Body.Close()
	sum := 0.0
	for _, day := range data.Prices {
		sum += day[1]
	}
	avg := sum / float64(days)
	log.Infof("%d day moving average: %s\n", days, fmtPrice(avg, cur))
	return avg, nil
}

func getConversionRate(sourceCur, targetCur string) (float64, error) {
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
		return 0, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}
	_ = res.Body.Close()
	rate := data.Rate[targetCur] / data.Rate[sourceCur]
	log.Infof("exchange rate for %s%s%s: %s%.4f%s\n", green, targetCur, reset, blue, rate, reset)
	return rate, nil
}

func getCurrentPrice(cur string) (float64, error) {
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
		return 0, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}
	_ = res.Body.Close()
	price := data.Price[cur]
	log.Infof("current price of bitcoin: %s\n", fmtPrice(price, cur))
	return price, nil
}
