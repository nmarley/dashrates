package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// BinanceAPI implements the RateAPI interface and contains info necessary for
// calling to the public Binance price ticker API.
type BinanceAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewBinanceAPI is a constructor for BinanceAPI.
func NewBinanceAPI() *BinanceAPI {
	return &BinanceAPI{
		BaseAPIURL:          "https://api.binance.com",
		PriceTickerEndpoint: "/api/v3/ticker/price?symbol=DASHBTC",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *BinanceAPI) DisplayName() string {
	return "Binance"
}

// FetchRate gets the Dash exchange rate from the Binance API.
//
// This is part of the RateAPI interface implementation.
func (a *BinanceAPI) FetchRate() (*RateInfo, error) {
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
	var res binancePriceResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	price, err := strconv.ParseFloat(res.Price, 64)
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "BTC",
		LastPrice:       price,
		BaseAssetVolume: 0,
		FetchTime:       now,
	}

	return &ri, nil
}

// binancePriceResp is used in parsing the Binance API response only.
type binancePriceResp struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}
