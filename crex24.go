package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// Crex24API implements the RateAPI interface and contains info necessary for
// calling to the public Crex24 price ticker API.
type Crex24API struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewCrex24API is a constructor for Crex24API.
func NewCrex24API() *Crex24API {
	return &Crex24API{
		BaseAPIURL:          "https://api.crex24.com/v2/public",
		PriceTickerEndpoint: "/tickers?instrument=DASH-BTC",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *Crex24API) DisplayName() string {
	return "CREX24"
}

// FetchRate gets the Dash exchange rate from the Crex24 API.
//
// This is part of the RateAPI interface implementation.
func (a *Crex24API) FetchRate() (*RateInfo, error) {
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
	var res crex24PubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "BTC",
		LastPrice:       res[0].Last,
		BaseAssetVolume: res[0].BaseVolume,
		FetchTime:       now,
	}

	return &ri, nil
}

// crex24PubTickerResp is used in parsing the Crex24 API response only.
type crex24PubTickerResp []crex24PubTickerData

// crex24PubTickerData is used in parsing the Crex24 API Dataonse only.
type crex24PubTickerData struct {
	Instrument    string    `json:"instrument"`
	Last          float64   `json:"last"`
	PercentChange float64   `json:"PercentChange"`
	Low           float64   `json:"low"`
	High          float64   `json:"high"`
	BaseVolume    float64   `json:"baseVolume"`
	QuoteVolume   float64   `json:"quoteVolume"`
	VolumeInBtc   float64   `json:"volumeInBtc"`
	VolumeInUsd   float64   `json:"volumeInUsd"`
	Ask           float64   `json:"ask"`
	Bid           float64   `json:"bid"`
	Timestamp     time.Time `json:"timestamp"`
}
