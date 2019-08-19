package dashrates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// ExmoAPI implements the RateAPI interface and contains info necessary for
// calling to the public Exmo price ticker API.
type ExmoAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewExmoAPI is a constructor for ExmoAPI.
func NewExmoAPI() *ExmoAPI {
	return &ExmoAPI{
		BaseAPIURL:          "https://api.exmo.com",
		PriceTickerEndpoint: "/v1/ticker/",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *ExmoAPI) DisplayName() string {
	return "Exmo"
}

// FetchRate gets the Dash exchange rate from the Exmo API.
//
// This is part of the RateAPI interface implementation.
func (a *ExmoAPI) FetchRate() (*RateInfo, error) {
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
	var res exmoPubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	pair, ok := res["DASH_USD"]
	if !ok {
		err = fmt.Errorf("oh no, %s does not have DASH/USD pair", a.DisplayName())
		return nil, err
	}
	data, err := pair.Normalize()
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

// exmoPubTickerResp is used in parsing the Exmo API response only.
type exmoPubTickerResp map[string]exmoPubTickerJSON

// exmoPubTickerData is used in parsing the Exmo API response only.
type exmoPubTickerData struct {
	Last        float64
	High        float64
	Low         float64
	Avg         float64
	BaseVolume  float64
	QuoteVolume float64
	BuyPrice    float64
	SellPrice   float64
	Updated     int
}

// exmoPubTickerJSON is used in parsing the Exmo API response only.
//
// This contains the raw data from the API. It will be parsed into a
// exmoPubTickerData with proper data types.
type exmoPubTickerJSON struct {
	Last      string `json:"last_trade"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Avg       string `json:"avg"`
	Vol       string `json:"vol"`
	VolCurr   string `json:"vol_curr"`
	BuyPrice  string `json:"buy_price"`
	SellPrice string `json:"sell_price"`
	Updated   int    `json:"updated"`
}

// Normalize parses the fields in exmoPubTickerJSON and returns a
// exmoPubTickerData with proper data types.
func (resp *exmoPubTickerJSON) Normalize() (*exmoPubTickerData, error) {
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
	avg, err := strconv.ParseFloat(resp.Avg, 64)
	if err != nil {
		return nil, err
	}
	baseVol, err := strconv.ParseFloat(resp.Vol, 64)
	if err != nil {
		return nil, err
	}
	quoteVol, err := strconv.ParseFloat(resp.VolCurr, 64)
	if err != nil {
		return nil, err
	}
	buyPrice, err := strconv.ParseFloat(resp.BuyPrice, 64)
	if err != nil {
		return nil, err
	}
	sellPrice, err := strconv.ParseFloat(resp.SellPrice, 64)
	if err != nil {
		return nil, err
	}

	return &exmoPubTickerData{
		Last:        last,
		High:        high,
		Low:         low,
		Avg:         avg,
		BaseVolume:  baseVol,
		QuoteVolume: quoteVol,
		BuyPrice:    buyPrice,
		SellPrice:   sellPrice,
		Updated:     resp.Updated,
	}, nil
}
