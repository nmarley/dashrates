package dashrates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// CoinbaseAPI implements the RateAPI interface and contains info necessary for
// calling to the public Coinbase price ticker API.
type CoinbaseAPI struct {
	BaseAPIURL            string
	PriceTickerEndpoint   string
	ExchangeRatesEndpoint string
}

// NewCoinbaseAPI is a constructor for CoinbaseAPI.
func NewCoinbaseAPI() *CoinbaseAPI {
	return &CoinbaseAPI{
		BaseAPIURL: "https://api.coinbase.com",
		// TODO: Update this when DASH is added to Coinbase
		PriceTickerEndpoint:   "/v2/exchange-rates?currency=LTC",
		ExchangeRatesEndpoint: "/v2/exchange-rates?currency=LTC",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *CoinbaseAPI) DisplayName() string {
	return "Coinbase"
}

// FetchRate gets the Dash exchange rate from the Coinbase API.
//
// This is part of the RateAPI interface implementation.
func (a *CoinbaseAPI) FetchRate() (*RateInfo, error) {
	resp, err := http.Get(a.BaseAPIURL + a.ExchangeRatesEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	now := time.Now()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// parse json and extract Dash rate
	var res coinbaseExchangeRatesResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	rate, ok := res.Data.Rates["USD"]
	if !ok {
		err = fmt.Errorf("oh no, %s does not have %s/USD pair",
			a.DisplayName(),
			res.Data.Currency,
		)
		return nil, err
	}

	price, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    res.Data.Currency,
		QuoteCurrency:   "USD",
		LastPrice:       price,
		BaseAssetVolume: 0,
		FetchTime:       now,
	}

	return &ri, nil
}

// coinbaseExchangeRatesResp is used in parsing the Coinbase API response only.
type coinbaseExchangeRatesResp struct {
	Data struct {
		Currency string            `json:"currency"`
		Rates    map[string]string `json:"rates"`
	} `json:"data"`
}
