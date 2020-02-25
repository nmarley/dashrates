package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// WhitebitAPI implements the RateAPI interface and contains info necessary for
// calling to the public Whitebit price ticker API.
type WhitebitAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewWhitebitAPI is a constructor for WhitebitAPI.
func NewWhitebitAPI() *WhitebitAPI {
	return &WhitebitAPI{
		BaseAPIURL:          "https://whitebit.com",
		PriceTickerEndpoint: "/api/v1/public/ticker?market=DASH_USD",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *WhitebitAPI) DisplayName() string {
	return "Whitebit"
}

// FetchRate gets the Dash exchange rate from the Whitebit API.
//
// This is part of the RateAPI interface implementation.
func (a *WhitebitAPI) FetchRate() (*RateInfo, error) {
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
	var res whitebitPubTickerResp
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
		QuoteCurrency:   "USDT",
		LastPrice:       data.Last,
		BaseAssetVolume: data.Volume,
		FetchTime:       now,
	}

	return &ri, nil
}

// whitebitPubTickerResp is used in parsing the Whitebit API response only.
type whitebitPubTickerResp struct {
	Success bool
	Message string
	Result struct {
		Bid            string `json:"bid"`
		Ask            string `json:"ask"`
		Open           string `json:"open"`
		High           string `json:"high"`
		Low            string `json:"low"`
		Last           string `json:"last"`
		Volume         string `json:"volume"`
		Deal           string `json:"deal"`
		Change         string `json:"change"`
	}
}

// whitebitPubTickerData is used in parsing the Whitebit API response only.
type whitebitPubTickerData struct {
	Bid            float64
	Ask            float64
	Open           float64
	High           float64
	Low            float64
	Last           float64
	Volume         float64
	Deal           float64
	Change         float64
}

// Normalize parses the fields in whitebitPubTickerResp and returns a
// whitebitPubTickerData with proper data types.
func (resp *whitebitPubTickerResp) Normalize() (*whitebitPubTickerData, error) {
	bid, err := strconv.ParseFloat(resp.Result.Bid, 64)
	if err != nil {
		return nil, err
	}

	ask, err := strconv.ParseFloat(resp.Result.Ask, 64)
	if err != nil {
		return nil, err
	}

	open, err := strconv.ParseFloat(resp.Result.Open, 64)
	if err != nil {
		return nil, err
	}

	high, err := strconv.ParseFloat(resp.Result.High, 64)
	if err != nil {
		return nil, err
	}

	low, err := strconv.ParseFloat(resp.Result.Low, 64)
	if err != nil {
		return nil, err
	}

	last, err := strconv.ParseFloat(resp.Result.Last, 64)
	if err != nil {
		return nil, err
	}

	volume, err := strconv.ParseFloat(resp.Result.Volume, 64)
	if err != nil {
		return nil, err
	}

	deal, err := strconv.ParseFloat(resp.Result.Deal, 64)
	if err != nil {
		return nil, err
	}

	change, err := strconv.ParseFloat(resp.Result.Change, 64)
	if err != nil {
		return nil, err
	}

	return &whitebitPubTickerData{
		Bid:    bid,
		Ask:    ask,
		Open:   open,
		High:   high,
		Low:    low,
		Last:   last,
		Volume: volume,
		Deal:   deal,
		Change: change,
	}, nil
}
