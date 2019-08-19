package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// CoinCapAPI implements the RateAPI interface and contains info necessary for
// calling to the public CoinCap price ticker API.
type CoinCapAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewCoinCapAPI is a constructor for CoinCapAPI.
func NewCoinCapAPI() *CoinCapAPI {
	return &CoinCapAPI{
		BaseAPIURL:          "https://api.coincap.io",
		PriceTickerEndpoint: "/v2/rates/bitcoin",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *CoinCapAPI) DisplayName() string {
	return "CoinCap"
}

// FetchRate gets the Dash exchange rate from the CoinCap API.
//
// This is part of the RateAPI interface implementation.
func (a *CoinCapAPI) FetchRate() (*RateInfo, error) {
	resp, err := http.Get(a.BaseAPIURL + a.PriceTickerEndpoint)
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
	var res coinCapPubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	rateUSD, err := res.GetRateUSD()
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "BTC",
		QuoteCurrency:   "USD",
		LastPrice:       rateUSD,
		BaseAssetVolume: 0,
		FetchTime:       now,
	}

	return &ri, nil
}

// coinCapPubTickerResp is used in parsing the CoinCap API response only.
type coinCapPubTickerResp struct {
	Data struct {
		ID             string `json:"id"`
		Symbol         string `json:"symbol"`
		CurrencySymbol string `json:"currencySymbol"`
		Type           string `json:"type"`
		RateUSD        string `json:"rateUsd"`
	} `json:"data"`
	Timestamp int64 `json:"timestamp"`
}

// GetRateUSD returns the USD Rate without bothering to re-write an entire
// struct with mostly strings.
func (resp *coinCapPubTickerResp) GetRateUSD() (float64, error) {
	rate, err := strconv.ParseFloat(resp.Data.RateUSD, 64)
	if err != nil {
		return 0, err
	}
	return rate, nil
}
