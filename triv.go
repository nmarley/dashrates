package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// TrivAPI implements the RateAPI interface and contains info necessary for
// calling to the public Triv price ticker API.
type TrivAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewTrivAPI is a constructor for TrivAPI.
func NewTrivAPI() *TrivAPI {
	return &TrivAPI{
		BaseAPIURL:          "https://triv.id",
		PriceTickerEndpoint: "/api/v1/config/ticker?pair=USD",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *TrivAPI) DisplayName() string {
	return "Triv"
}

// FetchRate gets the Dash exchange rate from the Triv API.
//
// This is part of the RateAPI interface implementation.
func (a *TrivAPI) FetchRate() (*RateInfo, error) {
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
	var res []*TrivPriceResp

	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	var x2 []*TrivPriceResp
	for _, v := range res {
		if v.Code == "DASH" {
			x2 = append(x2, v)
		}
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "USD",
		LastPrice:       x2[0].Buy,
		BaseAssetVolume: 0,
		FetchTime:       now,
	}

	return &ri, nil
}

// TrivPriceResp is used in parsing the Triv API response only.
type TrivPriceResp struct {
	Code string  `json:"code"`
	Name string  `json:"name"`
	Sell float64 `json:"sell"`
	Buy  float64 `json:"buy"`
}
