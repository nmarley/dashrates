package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// BiboxAPI implements the RateAPI interface and contains info necessary for
// calling to the public Bibox price ticker API.
type BiboxAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewBiboxAPI is a constructor for BiboxAPI.
func NewBiboxAPI() *BiboxAPI {
	return &BiboxAPI{
		BaseAPIURL:          "https://api.bibox.com",
		PriceTickerEndpoint: "/v1/mdata?cmd=market&pair=DASH_BTC",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *BiboxAPI) DisplayName() string {
	return "Bibox"
}

// FetchRate gets the Dash exchange rate from the Bibox API.
//
// This is part of the RateAPI interface implementation.
func (a *BiboxAPI) FetchRate() (*RateInfo, error) {
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
	var res biboxPubTickerResp
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

// biboxPubTickerResp is used in parsing the Bibox API response only.
type biboxPubTickerResp struct {
	Result struct {
		IsHide         int    `json:"is_hide"`
		HighCny        string `json:"high_cny"`
		Amount         string `json:"amount"`
		CoinSymbol     string `json:"coin_symbol"`
		Last           string `json:"last"`
		CurrencySymbol string `json:"currency_symbol"`
		Change         string `json:"change"`
		LowCny         string `json:"low_cny"`
		BaseLastCny    string `json:"base_last_cny"`
		AreaID         int    `json:"area_id"`
		Percent        string `json:"percent"`
		LastCny        string `json:"last_cny"`
		High           string `json:"high"`
		Low            string `json:"low"`
		PairType       int    `json:"pair_type"`
		LastUsd        string `json:"last_usd"`
		Vol24h         string `json:"vol24H"`
		ID             int    `json:"id"`
		HighUsd        string `json:"high_usd"`
		LowUsd         string `json:"low_usd"`
	}
	Cmd string `json:"cmd"`
	Ver string `json:"ver"`
}

// biboxPubTickerData is used in parsing the Bibox API response only.
type biboxPubTickerData struct {
	IsHide         int
	HighCny        float64
	Amount         float64
	CoinSymbol     string
	Last           float64
	CurrencySymbol string
	Change         float64
	LowCny         float64
	BaseLastCny    float64
	AreaID         int
	Percent        string
	LastCny        float64
	High           float64
	Low            float64
	PairType       int
	LastUsd        float64
	Vol24h         float64
	ID             int
	HighUsd        float64
	LowUsd         float64
}

// Normalize parses the fields in biboxPubTickerResp and returns a
// biboxPubTickerData with proper data types.
func (resp *biboxPubTickerResp) Normalize() (*biboxPubTickerData, error) {
	highCny, err := strconv.ParseFloat(resp.Result.HighCny, 64)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseFloat(resp.Result.Amount, 64)
	if err != nil {
		return nil, err
	}

	last, err := strconv.ParseFloat(resp.Result.Last, 64)
	if err != nil {
		return nil, err
	}

	change, err := strconv.ParseFloat(resp.Result.Change, 64)
	if err != nil {
		return nil, err
	}

	lowCny, err := strconv.ParseFloat(resp.Result.LowCny, 64)
	if err != nil {
		return nil, err
	}

	baseLastCny, err := strconv.ParseFloat(resp.Result.BaseLastCny, 64)
	if err != nil {
		return nil, err
	}

	lastCny, err := strconv.ParseFloat(resp.Result.LastCny, 64)
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

	lastUsd, err := strconv.ParseFloat(resp.Result.LastUsd, 64)
	if err != nil {
		return nil, err
	}

	vol24H, err := strconv.ParseFloat(resp.Result.Vol24h, 64)
	if err != nil {
		return nil, err
	}

	highUsd, err := strconv.ParseFloat(resp.Result.HighUsd, 64)
	if err != nil {
		return nil, err
	}

	lowUsd, err := strconv.ParseFloat(resp.Result.LowUsd, 64)
	if err != nil {
		return nil, err
	}

	return &biboxPubTickerData{
		IsHide:         resp.Result.IsHide,
		HighCny:        highCny,
		Amount:         amount,
		CoinSymbol:     resp.Result.CoinSymbol,
		Last:           last,
		CurrencySymbol: resp.Result.CurrencySymbol,
		Change:         change,
		LowCny:         lowCny,
		BaseLastCny:    baseLastCny,
		AreaID:         resp.Result.AreaID,
		Percent:        resp.Result.Percent,
		LastCny:        lastCny,
		High:           high,
		Low:            low,
		PairType:       resp.Result.PairType,
		LastUsd:        lastUsd,
		Vol24h:         vol24H,
		ID:             resp.Result.ID,
		HighUsd:        highUsd,
		LowUsd:         lowUsd,
	}, nil
}
