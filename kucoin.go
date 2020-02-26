package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// KucoinAPI implements the RateAPI interface and contains info necessary for
// calling to the public Kucoin price ticker API.
type KucoinAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewKucoinAPI is a constructor for KucoinAPI.
func NewKucoinAPI() *KucoinAPI {
	return &KucoinAPI{
		BaseAPIURL:          "https://api.kucoin.com",
		PriceTickerEndpoint: "/api/v1/market/orderbook/level1?symbol=DASH-BTC",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *KucoinAPI) DisplayName() string {
	return "Kucoin"
}

// FetchRate gets the Dash exchange rate from the Kucoin API.
//
// This is part of the RateAPI interface implementation.
func (a *KucoinAPI) FetchRate() (*RateInfo, error) {
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
	var res kucoinPubTickerResp
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
		LastPrice:       data.Price,
		BaseAssetVolume: 0,
		FetchTime:       now,
	}

	return &ri, nil
}

// kucoinPubTickerResp is used in parsing the Kucoin API response only.
type kucoinPubTickerResp struct {
	Code string `json:"code"`
	Data struct {
		Sequence    string `json:"sequence"`
		BestAsk     string `json:"bestAsk"`
		Size        string `json:"size"`
		Price       string `json:"price"`
		BestBidSize string `json:"bestBidSize"`
		Time        int64  `json:"time"`
		BestBid     string `json:"bestBid"`
		BestAskSize string `json:"bestAskSize"`
	}
}

// kucoinTickerData has the output of the parsed KucoinPubTickerResp and
// has proper data types.
type kucoinTickerData struct {
	Sequence    int64
	BestAsk     float64
	Size        float64
	Price       float64
	BestBidSize float64
	Time        int64
	BestBid     float64
	BestAskSize float64
}

// Normalize parses the fields in KucoinPubTickerResp and returns a
// kucoinTickerData with proper data types.
func (resp *kucoinPubTickerResp) Normalize() (*kucoinTickerData, error) {
	sequence, err := strconv.ParseInt(resp.Data.Sequence, 10, 64)
	if err != nil {
		return nil, err
	}

	bestAsk, err := strconv.ParseFloat(resp.Data.BestAsk, 64)
	if err != nil {
		return nil, err
	}

	size, err := strconv.ParseFloat(resp.Data.Size, 64)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(resp.Data.Price, 64)
	if err != nil {
		return nil, err
	}

	bestBidSize, err := strconv.ParseFloat(resp.Data.BestBidSize, 64)
	if err != nil {
		return nil, err
	}

	bestBid, err := strconv.ParseFloat(resp.Data.BestBid, 64)
	if err != nil {
		return nil, err
	}
	bestAskSize, err := strconv.ParseFloat(resp.Data.BestAskSize, 64)
	if err != nil {
		return nil, err
	}

	return &kucoinTickerData{
		Sequence:    sequence,
		BestAsk:     bestAsk,
		Size:        size,
		Price:       price,
		BestBidSize: bestBidSize,
		Time:        resp.Data.Time,
		BestBid:     bestBid,
		BestAskSize: bestAskSize,
	}, nil
}
