package dashrates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// YobitAPI implements the RateAPI interface and contains info necessary for
// calling to the public Yobit price ticker API.
type YobitAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewYobitAPI is a constructor for YobitAPI.
func NewYobitAPI() *YobitAPI {
	return &YobitAPI{
		BaseAPIURL:          "https://yobit.net",
		PriceTickerEndpoint: "/api/3/ticker/dash_usd",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *YobitAPI) DisplayName() string {
	return "Yobit"
}

// FetchRate gets the Dash exchange rate from the Yobit API.
//
// This is part of the RateAPI interface implementation.
func (a *YobitAPI) FetchRate() (*RateInfo, error) {
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
	var res yobitPubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	data, ok := res["dash_usd"]
	if !ok {
		err = fmt.Errorf("oh no, %s does not have DASH/USD pair", a.DisplayName())
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "USD",
		LastPrice:       data.Last,
		BaseAssetVolume: data.BaseVolume,
		FetchTime:       now,
	}

	return &ri, nil
}

// yobitPubTickerResp is used in parsing the Yobit API response only.
type yobitPubTickerResp map[string]yobitPubTickerData

// yobitPubTickerData is used in parsing the Yobit API response only.
//
// Like Poloniex and Bittrex, Yobit gets base and quote volume flipped.
type yobitPubTickerData struct {
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	Avg         float64 `json:"avg"`
	BaseVolume  float64 `json:"vol_cur"`
	QuoteVolume float64 `json:"vol"`
	Last        float64 `json:"last"`
	Buy         float64 `json:"buy"`
	Sell        float64 `json:"sell"`
	Updated     int     `json:"updated"`
}
