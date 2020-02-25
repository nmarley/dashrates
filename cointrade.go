package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// CointradeAPI implements the RateAPI interface and contains info necessary for
// calling to the public Cointrade price ticker API.
type CointradeAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewCointradeAPI is a constructor for CointradeAPI.
func NewCointradeAPI() *CointradeAPI {
	return &CointradeAPI{
		BaseAPIURL:          "https://api.cointradecx.com",
		PriceTickerEndpoint: "/public/ticker?market=DASH_BTC",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *CointradeAPI) DisplayName() string {
	return "Cointrade"
}

// FetchRate gets the Dash exchange rate from the Cointrade API.
//
// This is part of the RateAPI interface implementation.
func (a *CointradeAPI) FetchRate() (*RateInfo, error) {
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
	var res cointradePubTickerResp
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
		BaseAssetVolume: data.Vol24h,
		FetchTime:       now,
	}

	return &ri, nil
}

// cointradePubTickerResp is used in parsing the Cointrade API response only.
type cointradePubTickerResp struct {
	Success bool
	Message string
	Result []struct {
		Timestamp	 int64  `json:"timestamp"`
		Market       string `json:"market"`
		Ask          string `json:"ask"`
		Bid          string `json:"bid"`
		Last         string `json:"last"`
		Spread       string `json:"spread"`
		Low24h       string `json:"low24h"`
		High24h      string `json:"high24h"`
		Vol24h       string `json:"vol24h"`
		QuoteVolume  string `json:"quoteVolume"`
		IsFrozen     int    `json:"isFrozen"`
	}
}

// cointradePubTickerData is used in parsing the Cointrade API response only.
type cointradePubTickerData struct {
		Timestamp	   int64 
		Market         string
		Ask            float64
		Bid            float64
		Last           float64
		Spread         float64
		Low24h         float64
		High24h        float64
		Vol24h         float64
		QuoteVolume    float64
		IsFrozen       bool
}

// Normalize parses the fields in cointradePubTickerResp and returns a
// cointradePubTickerData with proper data types.
func (resp *cointradePubTickerResp) Normalize() (*cointradePubTickerData, error) {
	ask, err := strconv.ParseFloat(resp.Result[0].Ask, 64)
	if err != nil {
		return nil, err
	}

	bid, err := strconv.ParseFloat(resp.Result[0].Bid, 64)
	if err != nil {
		return nil, err
	}

	last, err := strconv.ParseFloat(resp.Result[0].Last, 64)
	if err != nil {
		return nil, err
	}

	spread, err := strconv.ParseFloat(resp.Result[0].Spread, 64)
	if err != nil {
		return nil, err
	}

	low24h, err := strconv.ParseFloat(resp.Result[0].Low24h, 64)
	if err != nil {
		return nil, err
	}

	high24h, err := strconv.ParseFloat(resp.Result[0].High24h, 64)
	if err != nil {
		return nil, err
	}

	vol24h, err := strconv.ParseFloat(resp.Result[0].Vol24h, 64)
	if err != nil {
		return nil, err
	}

	quoteVolume, err := strconv.ParseFloat(resp.Result[0].QuoteVolume, 64)
	if err != nil {
		return nil, err
	}

	isFrozenBool := false
	if resp.Result[0].IsFrozen != 0 {
		isFrozenBool = true
	}

	return &cointradePubTickerData{
		Timestamp:	    resp.Result[0].Timestamp ,
		Market:         resp.Result[0].Market,
		Ask:            ask,
		Bid:            bid,
		Last:           last,
		Spread:         spread,
		Low24h:         low24h,
		High24h:        high24h,
		Vol24h:         vol24h,
		QuoteVolume:    quoteVolume,
		IsFrozen:       isFrozenBool,
	}, nil
}
