package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// BittrexAPI implements the RateAPI interface and contains info necessary for
// calling to the public Bittrex price ticker API.
type BittrexAPI struct {
	BaseAPIURL            string
	PriceTickerEndpoint   string
	MarketSummaryEndpoint string
}

// NewBittrexAPI is a constructor for BittrexAPI.
func NewBittrexAPI() *BittrexAPI {
	return &BittrexAPI{
		BaseAPIURL:            "https://api.bittrex.com",
		PriceTickerEndpoint:   "/api/v1.1/public/getticker?market=BTC-DASH",
		MarketSummaryEndpoint: "/api/v1.1/public/getmarketsummary?market=btc-dash",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *BittrexAPI) DisplayName() string {
	return "Bittrex"
}

// FetchRate gets the Dash exchange rate from the Bittrex API.
//
// This is part of the RateAPI interface implementation.
func (a *BittrexAPI) FetchRate() (*RateInfo, error) {
	resp, err := http.Get(a.BaseAPIURL + a.MarketSummaryEndpoint)
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
	var res bittrexPubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "BTC",
		LastPrice:       res.Result[0].Last,
		BaseAssetVolume: res.Result[0].BaseVolume,
		FetchTime:       now,
	}

	return &ri, nil
}

// bittrexPubTickerResp is used in parsing the Bittrex API response only.
type bittrexPubTickerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	// Result BittrexPriceTickerResult `json:"result"`
	Result []bittrexMarketSummaryResult `json:"result"`
}

// bittrexPriceTickerResult is used for parsing the result from Bittrex
// PriceTickerEndpoint
type bittrexPriceTickerResult struct {
	Bid  float64 `json:"bid"`
	Ask  float64 `json:"ask"`
	Last float64 `json:"last"`
}

// bittrexMarketSummaryResult is used for parsing the result from Bittrex
// MarketSummaryEndpoint
//
// Note that like Poloniex, Bittrex gets base and quote volume confused. We
// swap them back here upon parsing so that the data makes sense, and
// explicitly name the fields `BaseVolume` and `QuoteVolume`.
//
// Without extra processing, the json Marshaler can't read the timestamp fields
// that Bittrex spits out, so we comment them b/c it's not worth the ROI to
// worry about this.
type bittrexMarketSummaryResult struct {
	MarketName     string  `json:"MarketName"`
	High           float64 `json:"High"`
	Low            float64 `json:"Low"`
	QuoteVolume    float64 `json:"BaseVolume"`
	Last           float64 `json:"Last"`
	BaseVolume     float64 `json:"Volume"`
	Bid            float64 `json:"Bid"`
	Ask            float64 `json:"Ask"`
	OpenBuyOrders  int     `json:"OpenBuyOrders"`
	OpenSellOrders int     `json:"OpenSellOrders"`
	PrevDay        float64 `json:"PrevDay"`
	//Timestamp      time.Time `json:"TimeStamp"`
	//Created        time.Time `json:"Created"`
}
