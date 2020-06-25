package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// BitbnsAPI implements the RateAPI interface and contains info necessary for
// calling to the public Bitbns price ticker API.
type BitbnsAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewBitbnsAPI is a constructor for BitbnsAPI.
func NewBitbnsAPI() *BitbnsAPI {
	return &BitbnsAPI{
		BaseAPIURL:          "https://bitbns.com",
		PriceTickerEndpoint: "/order/getTickerWithVolume/",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *BitbnsAPI) DisplayName() string {
	return "Bitbns"
}

// FetchRate gets the Dash exchange rate from the Bitbns API.
//
// This is part of the RateAPI interface implementation.
func (a *BitbnsAPI) FetchRate() (*RateInfo, error) {
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
	var res bitbnsPriceResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "USD",
		LastPrice:       res.Dashusdt.LastTradedPrice,
		BaseAssetVolume: 0,
		FetchTime:       now,
	}

	return &ri, nil
}

// bitbnsPriceResp is used in parsing the Bitbns API response only.
type bitbnsPriceResp struct {
	Dashusdt struct {
		HighestBuyBid   float64 `json:"highest_buy_bid"`
		LowestSellBid   float64 `json:"lowest_sell_bid"`
		LastTradedPrice float64 `json:"last_traded_price"`
		YesPrice        float64 `json:"yes_price"`
		InrPrice        float64 `json:"inr_price"`
		Volume          struct{}
	} `json:"DASHUSDT"`
}
