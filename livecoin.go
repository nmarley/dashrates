package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// LivecoinAPI implements the RateAPI interface and contains info necessary for
// calling to the public Livecoin price ticker API.
type LivecoinAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewLivecoinAPI is a constructor for LivecoinAPI.
func NewLivecoinAPI() *LivecoinAPI {
	return &LivecoinAPI{
		BaseAPIURL:          "https://api.livecoin.net",
		PriceTickerEndpoint: "/exchange/ticker?currencyPair=DASH/USD",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *LivecoinAPI) DisplayName() string {
	return "Livecoin"
}

// FetchRate gets the Dash exchange rate from the Livecoin API.
//
// This is part of the RateAPI interface implementation.
func (a *LivecoinAPI) FetchRate() (*RateInfo, error) {
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
	var res livecoinPubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "USD",
		LastPrice:       res.Last,
		BaseAssetVolume: res.Volume,
		FetchTime:       now,
	}

	return &ri, nil
}

// livecoinPubTickerResp is used in parsing the Livecoin API response only.
type livecoinPubTickerResp struct {
	Currency string  `json:"cur"`
	Symbol   string  `json:"symbol"`
	Last     float64 `json:"last"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Volume   float64 `json:"volume"`
	VWAP     float64 `json:"vwap"`
	MaxBid   float64 `json:"max_bid"`
	MinAsk   float64 `json:"min_ask"`
	BestBid  float64 `json:"best_bid"`
	BestAsk  float64 `json:"best_ask"`
}
