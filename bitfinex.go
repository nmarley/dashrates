package dashrates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// BitfinexAPI implements the RateAPI interface and contains info necessary for
// calling to the public Bitfinex price ticker API.
type BitfinexAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewBitfinexAPI is a constructor for BitfinexAPI.
func NewBitfinexAPI() *BitfinexAPI {
	return &BitfinexAPI{
		BaseAPIURL:          "https://api.bitfinex.com",
		PriceTickerEndpoint: "/v1/pubticker/dshusd",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *BitfinexAPI) DisplayName() string {
	return "Bitfinex"
}

// FetchRate gets the Dash exchange rate from the Bitfinex API.
//
// This is part of the RateAPI interface implementation.
func (a *BitfinexAPI) FetchRate() (*RateInfo, error) {
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
	var res bitfinexPubTickerResp
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
		LastPrice:       data.LastPrice,
		BaseAssetVolume: data.Volume,
		FetchTime:       now,
	}

	return &ri, nil
}

// bitfinexTickerData has the output of the parsed BitfinexPubTickerResp and
// has proper data types.
type bitfinexTickerData struct {
	Mid       float64
	Bid       float64
	Ask       float64
	LastPrice float64
	Low       float64
	High      float64
	Volume    float64
	Timestamp time.Time
}

// bitfinexPubTickerResp is used in parsing the Bitfinex API response only.
type bitfinexPubTickerResp struct {
	Mid       string `json:"mid"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	LastPrice string `json:"last_price"`
	Low       string `json:"low"`
	High      string `json:"high"`
	Volume    string `json:"volume"`
	Timestamp string `json:"timestamp"`
}

// Normalize parses the fields in BitfinexPubTickerResp and returns a
// bitfinexTickerData with proper data types.
func (resp *bitfinexPubTickerResp) Normalize() (*bitfinexTickerData, error) {
	mid, err := strconv.ParseFloat(resp.Mid, 64)
	if err != nil {
		return nil, err
	}

	bid, err := strconv.ParseFloat(resp.Bid, 64)
	if err != nil {
		return nil, err
	}

	ask, err := strconv.ParseFloat(resp.Ask, 64)
	if err != nil {
		return nil, err
	}

	last, err := strconv.ParseFloat(resp.LastPrice, 64)
	if err != nil {
		return nil, err
	}

	low, err := strconv.ParseFloat(resp.Low, 64)
	if err != nil {
		return nil, err
	}
	high, err := strconv.ParseFloat(resp.High, 64)
	if err != nil {
		return nil, err
	}
	volume, err := strconv.ParseFloat(resp.Volume, 64)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(resp.Timestamp, ".")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid timestamp %v", resp.Timestamp)
		return nil, err
	}
	sec, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	nsec, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}

	ts := time.Unix(sec, nsec)

	return &bitfinexTickerData{
		Mid:       mid,
		Bid:       bid,
		Ask:       ask,
		LastPrice: last,
		Low:       low,
		High:      high,
		Volume:    volume,
		Timestamp: ts,
	}, nil
}
