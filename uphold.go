package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// UpholdAPI implements the RateAPI interface and contains info necessary for
// calling to the public Uphold price ticker API.
type UpholdAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewUpholdAPI is a constructor for UpholdAPI.
func NewUpholdAPI() *UpholdAPI {
	return &UpholdAPI{
		BaseAPIURL:          "https://api.uphold.com",
		PriceTickerEndpoint: "/v0/ticker/DASHUSD",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *UpholdAPI) DisplayName() string {
	return "Uphold"
}

// FetchRate gets the Dash exchange rate from the Uphold API.
//
// This is part of the RateAPI interface implementation.
func (a *UpholdAPI) FetchRate() (*RateInfo, error) {
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
	var res upholdPubTickerResp
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
		QuoteCurrency:   "USD",
		LastPrice:       data.Ask,
		BaseAssetVolume: 0,
		FetchTime:       now,
	}

	return &ri, nil
}

// upholdPubTickerResp is used in parsing the Uphold API response only.
type upholdPubTickerResp struct {
		Ask      string `json:"ask"`
		Bid      string `json:"bid"`
		Currency string `json:"currency"`
}

// upholdPubTickerData is used in parsing the Uphold API response only.
type upholdPubTickerData struct {
	Ask      float64
	Bid      float64
	Currency string
}

// Normalize parses the fields in upholdPubTickerResp and returns a
// upholdPubTickerData with proper data types.
func (resp *upholdPubTickerResp) Normalize() (*upholdPubTickerData, error) {
	ask, err := strconv.ParseFloat(resp.Ask, 64)
	if err != nil {
		return nil, err
	}

	bid, err := strconv.ParseFloat(resp.Bid, 64)
	if err != nil {
		return nil, err
	}

	return &upholdPubTickerData{
		Ask:      ask,
		Bid:      bid,
		Currency: resp.Currency,
	}, nil
}
