package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// LiquidAPI implements the RateAPI interface and contains info necessary for
// calling to the public Liquid price ticker API.
type LiquidAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewLiquidAPI is a constructor for LiquidAPI.
func NewLiquidAPI() *LiquidAPI {
	return &LiquidAPI{
		BaseAPIURL:          "https://api.liquid.com",
		PriceTickerEndpoint: "/products/116",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *LiquidAPI) DisplayName() string {
	return "Liquid"
}

// FetchRate gets the Dash exchange rate from the Liquid API.
//
// This is part of the RateAPI interface implementation.
func (a *LiquidAPI) FetchRate() (*RateInfo, error) {
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
	var res liquidPubTickerResp
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
		LastPrice:       data.LastPrice,
		BaseAssetVolume: data.Volume24h,
		FetchTime:       now,
	}

	return &ri, nil
}

// liquidPubTickerData is used in parsing the Liquid API response only.
type liquidPubTickerData struct {
	ID        int64
	LastPrice float64
	Volume24h float64
}

// liquidPubTickerResp is used in parsing the Liquid API response only.
type liquidPubTickerResp struct {
	ID        string `json:"id"`
	LastPrice string `json:"last_traded_price"`
	Volume24h string `json:"volume_24h"`
}

// Normalize parses the fields in liquidPubTickerResp and returns a
// liquidPubTickerData with proper data types.
func (resp *liquidPubTickerResp) Normalize() (*liquidPubTickerData, error) {
	id, err := strconv.ParseInt(resp.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	lastPrice, err := strconv.ParseFloat(resp.LastPrice, 64)
	if err != nil {
		return nil, err
	}

	volume24h, err := strconv.ParseFloat(resp.Volume24h, 64)
	if err != nil {
		return nil, err
	}

	return &liquidPubTickerData{
		ID:        id,
		LastPrice: lastPrice,
		Volume24h: volume24h,
	}, nil
}
