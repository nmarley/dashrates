package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// SouthXchangeAPI implements the RateAPI interface and contains info necessary for
// calling to the public SouthXchange price ticker API.
type SouthXchangeAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewSouthXchangeAPI is a constructor for SouthXchangeAPI.
func NewSouthXchangeAPI() *SouthXchangeAPI {
	return &SouthXchangeAPI{
		BaseAPIURL:          "https://www.southxchange.com",
		PriceTickerEndpoint: "/api/price/DASH/BTC",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *SouthXchangeAPI) DisplayName() string {
	return "SouthXchange"
}

// FetchRate gets the Dash exchange rate from the SouthXchange API.
//
// This is part of the RateAPI interface implementation.
func (a *SouthXchangeAPI) FetchRate() (*RateInfo, error) {
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
	var res southxchangePubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "BTC",
		LastPrice:       res.Last,
		BaseAssetVolume: res.Volume24Hr,
		FetchTime:       now,
	}

	return &ri, nil
}

// southxchangePubTickerResp is used in parsing the SouthXchange API response only.
type southxchangePubTickerResp struct {
	Bid           float64 `json:"Bid"`
	Ask           float64 `json:"Ask"`
	Last          float64 `json:"Last"`
	Variation24Hr float64 `json:"Variation24Hr"`
	Volume24Hr    float64 `json:"Volume24Hr"`
}
