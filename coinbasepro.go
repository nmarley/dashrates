package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// CoinbaseProAPI implements the RateAPI interface and contains info necessary for
// calling to the public CoinbasePro price ticker API.
type CoinbaseProAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewCoinbaseProAPI is a constructor for CoinbaseProAPI.
func NewCoinbaseProAPI() *CoinbaseProAPI {
	return &CoinbaseProAPI{
		BaseAPIURL:          "https://api.pro.coinbase.com",
		PriceTickerEndpoint: "/products/DASH-USD/ticker",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *CoinbaseProAPI) DisplayName() string {
	return "Coinbase Pro"
}

// FetchRate gets the Dash exchange rate from the CoinbasePro API.
//
// This is part of the RateAPI interface implementation.
func (a *CoinbaseProAPI) FetchRate() (*RateInfo, error) {
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
	var res coinbaseProTickerResp
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
		LastPrice:       data.Price,
		BaseAssetVolume: data.Volume,
		FetchTime:       now,
	}

	return &ri, nil
}

// coinbaseProTickerResp is used in parsing the CoinbasePro API response only.
type coinbaseProTickerResp struct {
	TradeID   int    `json:"trade_id"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	Timestamp string `json:"time"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	Volume    string `json:"volume"`
}

// coinbaseProTickerData is used in parsing the CoinbasePro API response only.
type coinbaseProTickerData struct {
	TradeID   int
	Price     float64
	Size      float64
	Timestamp time.Time
	Bid       float64
	Ask       float64
	Volume    float64
}

// Normalize parses the fields in coinbaseProTickerResp and returns a
// coinbaseProTickerData with proper data types.
func (resp *coinbaseProTickerResp) Normalize() (*coinbaseProTickerData, error) {
	bid, err := strconv.ParseFloat(resp.Bid, 64)
	if err != nil {
		return nil, err
	}
	ask, err := strconv.ParseFloat(resp.Ask, 64)
	if err != nil {
		return nil, err
	}
	vol, err := strconv.ParseFloat(resp.Volume, 64)
	if err != nil {
		return nil, err
	}
	price, err := strconv.ParseFloat(resp.Price, 64)
	if err != nil {
		return nil, err
	}
	size, err := strconv.ParseFloat(resp.Size, 64)
	if err != nil {
		return nil, err
	}

	//layout := "2006-01-02T15:04:05.000Z"
	ts, err := time.Parse(time.RFC3339, resp.Timestamp)
	if err != nil {
		return nil, err
	}

	return &coinbaseProTickerData{
		TradeID:   resp.TradeID,
		Price:     price,
		Size:      size,
		Timestamp: ts,
		Bid:       bid,
		Ask:       ask,
		Volume:    vol,
	}, nil
}
