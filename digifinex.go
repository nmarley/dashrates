package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// DigifinexAPI implements the RateAPI interface and contains info necessary for
// calling to the public Digifinex price ticker API.
type DigifinexAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewDigifinexAPI is a constructor for DigifinexAPI.
func NewDigifinexAPI() *DigifinexAPI {
	return &DigifinexAPI{
		BaseAPIURL:          "https://openapi.digifinex.com",
		PriceTickerEndpoint: "/v3/ticker?symbol=dash_usdt",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *DigifinexAPI) DisplayName() string {
	return "Digifinex"
}

// FetchRate gets the Dash exchange rate from the Digifinex API.
//
// This is part of the RateAPI interface implementation.
func (a *DigifinexAPI) FetchRate() (*RateInfo, error) {
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
	var res digifinexPubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "USDT",
		LastPrice:       res.Ticker[0].Last,
		BaseAssetVolume: res.Ticker[0].BaseVol,
		FetchTime:       now,
	}

	return &ri, nil
}

// digifinexPriceResp is used in parsing the Digifinex API response only.
type digifinexPubTickerResp struct {
	Ticker []digifinexPubTickerData `json:"ticker"`
	Date   int64                    `json:"date"`
	Code   int64                    `json:"code"`
}

// digifinexPubTickerData is used in parsing the Digifinex API response only.
type digifinexPubTickerData struct {
	Vol     float64 `json:"vol"`
	Change  float64 `json:"change"`
	BaseVol float64 `json:"base_vol"`
	Sell    float64 `json:"sell"`
	Last    float64 `json:"last"`
	Symbol  string  `json:"symbol"`
	Low     float64 `json:"low"`
	Buy     float64 `json:"buy"`
	High    float64 `json:"high"`
}
