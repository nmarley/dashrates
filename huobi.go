package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// HuobiAPI implements the RateAPI interface and contains info necessary for
// calling to the public Huobi price ticker API.
type HuobiAPI struct {
	BaseAPIURL           string
	MarketDetailEndpoint string
	LastTradeEndpoint    string
}

// NewHuobiAPI is a constructor for HuobiAPI.
func NewHuobiAPI() *HuobiAPI {
	return &HuobiAPI{
		BaseAPIURL:           "https://api.huobi.pro",
		MarketDetailEndpoint: "/market/detail/merged?symbol=dashbtc",
		LastTradeEndpoint:    "/market/trade?symbol=dashbtc",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *HuobiAPI) DisplayName() string {
	return "Huobi"
}

// FetchRate gets the Dash exchange rate from the Huobi API.
//
// This is part of the RateAPI interface implementation.
func (a *HuobiAPI) FetchRate() (*RateInfo, error) {
	now := time.Now()

	// parse json and extract Dash rate
	//marketDetail, err := a.fetchMarketDetail()
	//if err != nil {
	//	return nil, err
	//}
	lastTradePrice, err := a.fetchLastTrade()
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:  "DASH",
		QuoteCurrency: "BTC",
		LastPrice:     lastTradePrice,
		//BaseAssetVolume: marketDetail.Tick.Volume,
		BaseAssetVolume: 0,
		FetchTime:       now,
	}

	return &ri, nil
}

// huobiPubTickerResp is used in parsing the Huobi API response only.
type huobiPubTickerResp struct {
	Status    string `json:"status"`
	Channel   string `json:"ch"`
	Timestamp int64  `json:"ts"`
	Tick      struct {
		ID      int64     `json:"id"`
		Close   float64   `json:"close"`
		Open    float64   `json:"open"`
		High    float64   `json:"high"`
		Low     float64   `json:"low"`
		Amount  float64   `json:"amount"`
		Count   int       `json:"count"`
		Version int64     `json:"version"`
		Volume  float64   `json:"vol"`
		Ask     []float64 `json:"ask"`
		Bid     []float64 `json:"bid"`
	} `json:"tick"`
}

// huobiLastTradeResp is used in parsing the Huobi API response only.
type huobiLastTradeResp struct {
	Status    string `json:"status"`
	Channel   string `json:"ch"`
	Timestamp int64  `json:"ts"`
	Tick      struct {
		// ID        int64 `json:"id"`
		Timestamp int64 `json:"ts"`
		Data      []struct {
			Amount    float64 `json:"amount"`
			Timestamp int64   `json:"ts"`
			// ID        int64   `json:"id"`
			Price     float64 `json:"price"`
			Direction string  `json:"direction"`
		} `json:"data"`
	} `json:"tick"`
}

// fetchLastTrade gets the Dash exchange rate from the Huobi API.
func (a *HuobiAPI) fetchLastTrade() (float64, error) {
	// Get last trade
	resp, err := http.Get(a.BaseAPIURL + a.LastTradeEndpoint)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// parse json and extract Dash rate
	var res huobiLastTradeResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return 0, err
	}

	return res.Tick.Data[0].Price, nil
}

// fetchMarketDetail gets the Dash market detail from the Huobi API.
func (a *HuobiAPI) fetchMarketDetail() (*huobiPubTickerResp, error) {
	resp, err := http.Get(a.BaseAPIURL + a.MarketDetailEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res huobiPubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
