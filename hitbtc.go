package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// HitBTCAPI implements the RateAPI interface and contains info necessary for
// calling to the public HitBTC price ticker API.
type HitBTCAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewHitBTCAPI is a constructor for HitBTCAPI.
func NewHitBTCAPI() *HitBTCAPI {
	return &HitBTCAPI{
		BaseAPIURL:          "https://api.hitbtc.com",
		PriceTickerEndpoint: "/api/2/public/ticker/DASHUSD",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *HitBTCAPI) DisplayName() string {
	return "HitBTC"
}

// FetchRate gets the Dash exchange rate from the HitBTC API.
//
// This is part of the RateAPI interface implementation.
func (a *HitBTCAPI) FetchRate() (*RateInfo, error) {
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
	var res hitBTCPubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	data, err := res.Normalize()

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "USD",
		LastPrice:       data.Last,
		BaseAssetVolume: data.BaseVolume,
		FetchTime:       now,
	}

	return &ri, nil
}

// hitBTCPubTickerData is used in parsing the HitBTC API response only.
//
// This contains the parsed fields and makes sense (correct data types).
type hitBTCPubTickerData struct {
	Symbol      string
	Ask         float64
	Bid         float64
	Last        float64
	High        float64
	Low         float64
	Open        float64
	BaseVolume  float64
	QuoteVolume float64
	Timestamp   time.Time
}

// hitBTCPubTickerResp is used in parsing the HitBTC API response only.
//
// This is the data taken directly from the API response and has all values as
// strings. Because this is the 21st century, we put a fucking man on the moon
// but exchanges still can't figure out correct data types.
type hitBTCPubTickerResp struct {
	Symbol    string `json:"symbol"`
	Ask       string `json:"ask"`
	Bid       string `json:"bid"`
	Last      string `json:"last"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Open      string `json:"open"`
	BaseVol   string `json:"volume"`
	QuoteVol  string `json:"volumeQuote"`
	Timestamp string `json:"timestamp"`
}

// Normalize parses the fields in hitBTCPubTickerResp and returns a
// hitBTCPubTickerData with proper data types.
func (resp *hitBTCPubTickerResp) Normalize() (*hitBTCPubTickerData, error) {
	ask, err := strconv.ParseFloat(resp.Ask, 64)
	if err != nil {
		return nil, err
	}
	bid, err := strconv.ParseFloat(resp.Bid, 64)
	if err != nil {
		return nil, err
	}
	last, err := strconv.ParseFloat(resp.Last, 64)
	if err != nil {
		return nil, err
	}
	high, err := strconv.ParseFloat(resp.High, 64)
	if err != nil {
		return nil, err
	}
	low, err := strconv.ParseFloat(resp.Low, 64)
	if err != nil {
		return nil, err
	}
	open, err := strconv.ParseFloat(resp.Open, 64)
	if err != nil {
		return nil, err
	}
	baseVol, err := strconv.ParseFloat(resp.BaseVol, 64)
	if err != nil {
		return nil, err
	}
	quoteVol, err := strconv.ParseFloat(resp.QuoteVol, 64)
	if err != nil {
		return nil, err
	}

	layout := "2006-01-02T15:04:05.000Z"
	ts, err := time.Parse(layout, resp.Timestamp)
	if err != nil {
		return nil, err
	}

	return &hitBTCPubTickerData{
		Symbol:      resp.Symbol,
		Ask:         ask,
		Bid:         bid,
		Last:        last,
		Open:        open,
		High:        high,
		Low:         low,
		BaseVolume:  baseVol,
		QuoteVolume: quoteVol,
		Timestamp:   ts,
	}, nil
}
