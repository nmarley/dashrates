package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// BvnexAPI implements the RateAPI interface and contains info necessary for
// calling to the public Bvnex price ticker API.
type BvnexAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewBvnexAPI is a constructor for BvnexAPI.
func NewBvnexAPI() *BvnexAPI {
	return &BvnexAPI{
		BaseAPIURL:          "https://api.bvnex.com",
		PriceTickerEndpoint: "/api/ticker/get?symbol=dash_usdt",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *BvnexAPI) DisplayName() string {
	return "Bvnex"
}

// FetchRate gets the Dash exchange rate from the Bvnex API.
//
// This is part of the RateAPI interface implementation.
func (a *BvnexAPI) FetchRate() (*RateInfo, error) {
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
	var res bvnexPubTickerResp
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
		LastPrice:       data.Last,
		BaseAssetVolume: data.BaseVolume,
		FetchTime:       now,
	}

	return &ri, nil
}

// bvnexPubTickerResp is used in parsing the Bvnex API response only.
type bvnexPubTickerResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Last          string `json:"last"`
		LowestAsk     string `json:"lowestAsk"`
		HighestBid    string `json:"highestBid"`
		PercentChange string `json:"percentChange"`
		BaseVolume    string `json:"baseVolume"`
		QuoteVolume   string `json:"quoteVolume"`
		High24hr      string `json:"high24hr"`
		Low24hr       string `json:"low24hr"`
	} `json:"data"`
}

// bvnexPubTickerData is used in parsing the Bvnex API response only.
type bvnexPubTickerData struct {
	Last          float64
	LowestAsk     float64
	HighestBid    float64
	PercentChange float64
	BaseVolume    float64
	QuoteVolume   float64
	High24hr      float64
	Low24hr       float64
}

// Normalize parses the fields in bvnexPubTickerResp and returns a
// bvnexPubTickerData with proper data types.
func (resp *bvnexPubTickerResp) Normalize() (*bvnexPubTickerData, error) {
	last, err := strconv.ParseFloat(resp.Data.Last, 64)
	if err != nil {
		return nil, err
	}

	lowestAsk, err := strconv.ParseFloat(resp.Data.LowestAsk, 64)
	if err != nil {
		return nil, err
	}

	highestBid, err := strconv.ParseFloat(resp.Data.HighestBid, 64)
	if err != nil {
		return nil, err
	}

	percentChange, err := strconv.ParseFloat(resp.Data.PercentChange, 64)
	if err != nil {
		return nil, err
	}

	baseVolume, err := strconv.ParseFloat(resp.Data.BaseVolume, 64)
	if err != nil {
		return nil, err
	}

	quoteVolume, err := strconv.ParseFloat(resp.Data.QuoteVolume, 64)
	if err != nil {
		return nil, err
	}

	high24hr, err := strconv.ParseFloat(resp.Data.High24hr, 64)
	if err != nil {
		return nil, err
	}

	low24hr, err := strconv.ParseFloat(resp.Data.Low24hr, 64)
	if err != nil {
		return nil, err
	}

	return &bvnexPubTickerData{
		Last:          last,
		LowestAsk:     lowestAsk,
		HighestBid:    highestBid,
		PercentChange: percentChange,
		BaseVolume:    baseVolume,
		QuoteVolume:   quoteVolume,
		High24hr:      high24hr,
		Low24hr:       low24hr,
	}, nil
}
