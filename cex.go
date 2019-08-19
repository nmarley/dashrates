package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CexAPI implements the RateAPI interface and contains info necessary for
// calling to the public Cex price ticker API.
type CexAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewCexAPI is a constructor for CexAPI.
func NewCexAPI() *CexAPI {
	return &CexAPI{
		BaseAPIURL:          "https://cex.io",
		PriceTickerEndpoint: "/api/ticker/DASH/USD",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *CexAPI) DisplayName() string {
	return "CEX.IO"
}

// FetchRate gets the Dash exchange rate from the Cex API.
//
// This is part of the RateAPI interface implementation.
func (a *CexAPI) FetchRate() (*RateInfo, error) {
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
	var res cexPubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	data, err := res.Normalize()

	ri := RateInfo{
		BaseCurrency:    data.Pair.Base,
		QuoteCurrency:   data.Pair.Quote,
		LastPrice:       data.Last,
		BaseAssetVolume: data.Volume,
		FetchTime:       now,
	}

	return &ri, nil
}

// cexPubTickerResp is used in parsing the Cex API response only.
type cexPubTickerResp struct {
	Timestamp             string  `json:"timestamp"`
	Low                   string  `json:"low"`
	High                  string  `json:"high"`
	Last                  string  `json:"last"`
	Volume                string  `json:"volume"`
	Volume30d             string  `json:"volume30d"`
	Bid                   float64 `json:"bid"`
	Ask                   float64 `json:"ask"`
	PriceChange           string  `json:"priceChange"`
	PriceChangePercentage string  `json:"priceChangePercentage"`
	Pair                  string  `json:"pair"`
}

// cexPubTickerData is used in parsing the Cex API response only.
type cexPubTickerData struct {
	Timestamp             time.Time
	Low                   float64
	High                  float64
	Last                  float64
	Volume                float64
	Volume30d             float64
	Bid                   float64
	Ask                   float64
	PriceChange           float64
	PriceChangePercentage float64
	Pair                  cexPair
}

// cexPair is used in parsing the Cex API response only.
type cexPair struct {
	Base  string
	Quote string
}

// Normalize parses the fields in cexPubTickerResp and returns a
// cexPubTickerData with proper data types.
func (resp *cexPubTickerResp) Normalize() (*cexPubTickerData, error) {
	tsEpoch, err := strconv.ParseInt(resp.Timestamp, 10, 64)
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

	last, err := strconv.ParseFloat(resp.Last, 64)
	if err != nil {
		return nil, err
	}

	volume, err := strconv.ParseFloat(resp.Volume, 64)
	if err != nil {
		return nil, err
	}

	volume30d, err := strconv.ParseFloat(resp.Volume30d, 64)
	if err != nil {
		return nil, err
	}

	priceChange, err := strconv.ParseFloat(resp.PriceChange, 64)
	if err != nil {
		return nil, err
	}
	priceChangePercentage, err := strconv.ParseFloat(resp.PriceChangePercentage, 64)
	if err != nil {
		return nil, err
	}

	pair := strings.Split(resp.Pair, ":")

	return &cexPubTickerData{
		Timestamp:             time.Unix(tsEpoch, 0),
		Low:                   low,
		High:                  high,
		Last:                  last,
		Volume:                volume,
		Volume30d:             volume30d,
		Bid:                   resp.Bid,
		Ask:                   resp.Ask,
		PriceChange:           priceChange,
		PriceChangePercentage: priceChangePercentage,
		Pair:                  cexPair{Base: pair[0], Quote: pair[1]},
	}, nil
}
