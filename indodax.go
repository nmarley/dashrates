package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// IndodaxAPI implements the RateAPI interface and contains info necessary for
// calling to the public Indodax price ticker API.
type IndodaxAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewIndodaxAPI is a constructor for IndodaxAPI.
func NewIndodaxAPI() *IndodaxAPI {
	return &IndodaxAPI{
		BaseAPIURL:          "https://indodax.com",
		PriceTickerEndpoint: "/api/drk_btc/ticker",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *IndodaxAPI) DisplayName() string {
	return "Indodax"
}

// FetchRate gets the Dash exchange rate from the Indodax API.
//
// This is part of the RateAPI interface implementation.
func (a *IndodaxAPI) FetchRate() (*RateInfo, error) {
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
	var res indodaxPubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	data, err := res.Normalize()
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "BTC",
		LastPrice:       data.Last,
		BaseAssetVolume: data.VolDrk,
		FetchTime:       now,
	}

	return &ri, nil
}

// indodaxPubTickerResp is used in parsing the Indodax API response only.
type indodaxPubTickerResp struct {
	Ticker struct {
		High       string `json:"high"`
		Low        string `json:"low"`
		VolDrk     string `json:"vol_drk"`
		VolBtc     string `json:"vol_btc"`
		Last       string `json:"last"`
		Buy        string `json:"buy"`
		Sell       string `json:"sell"`
		ServerTime int64  `json:"server_time"`
	}
}

// indodaxPubTickerData is used in parsing the Indodax API response only.
type indodaxPubTickerData struct {
	High       float64
	Low        float64
	VolDrk     float64
	VolBtc     float64
	Last       float64
	Buy        float64
	Sell       float64
	ServerTime time.Time
}

// Normalize parses the fields in indodaxPubTickerResp and returns a
// indodaxPubTickerData with proper data types.
func (resp *indodaxPubTickerResp) Normalize() (*indodaxPubTickerData, error) {
	high, err := strconv.ParseFloat(resp.Ticker.High, 64)
	if err != nil {
		return nil, err
	}

	low, err := strconv.ParseFloat(resp.Ticker.Low, 64)
	if err != nil {
		return nil, err
	}

	volDrk, err := strconv.ParseFloat(resp.Ticker.VolDrk, 64)
	if err != nil {
		return nil, err
	}

	volBtc, err := strconv.ParseFloat(resp.Ticker.VolBtc, 64)
	if err != nil {
		return nil, err
	}

	last, err := strconv.ParseFloat(resp.Ticker.Last, 64)
	if err != nil {
		return nil, err
	}

	buy, err := strconv.ParseFloat(resp.Ticker.Buy, 64)
	if err != nil {
		return nil, err
	}

	sell, err := strconv.ParseFloat(resp.Ticker.Sell, 64)
	if err != nil {
		return nil, err
	}

	return &indodaxPubTickerData{
		High:       high,
		Low:        low,
		VolDrk:     volDrk,
		VolBtc:     volBtc,
		Last:       last,
		Buy:        buy,
		Sell:       sell,
		ServerTime: time.Unix(resp.Ticker.ServerTime, 0),
	}, nil
}
