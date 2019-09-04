package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// BigONEAPI implements the RateAPI interface and contains info necessary for
// calling to the public BigONE price ticker API.
type BigONEAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewBigONEAPI is a constructor for BigONEAPI.
func NewBigONEAPI() *BigONEAPI {
	return &BigONEAPI{
		BaseAPIURL:          "https://big.one/api/v3",
		PriceTickerEndpoint: "/asset_pairs/DASH-BTC/ticker",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *BigONEAPI) DisplayName() string {
	return "BigONE"
}

// FetchRate gets the Dash exchange rate from the BigONE API.
//
// This is part of the RateAPI interface implementation.
func (a *BigONEAPI) FetchRate() (*RateInfo, error) {
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
	var res bigONEPubTickerResp
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
		LastPrice:       data.Data.Ask.Price,
		BaseAssetVolume: data.Data.Volume,
		FetchTime:       now,
	}

	return &ri, nil
}

// bigONEPubTickerResp is used in parsing the BigONE API response only.
type bigONEPubTickerResp struct {
	Code int `json:"code"`
	Data struct {
		AssetPairName string              `json:"asset_pair_name"`
		Bid           bigONEBidAskStrings `json:"bid"`
		Ask           bigONEBidAskStrings `json:"ask"`
		Open          string              `json:"open"`
		High          string              `json:"high"`
		Low           string              `json:"low"`
		Close         string              `json:"close"`
		Volume        string              `json:"volume"`
		DailyChange   string              `json:"daily_change"`
	} `json:"data"`
}

// bigONEBidAskStrings is used in parsing the BigONE API response only.
type bigONEBidAskStrings struct {
	Price      string `json:"price"`
	OrderCount int    `json:"order_count"`
	Quantity   string `json:"quantity"`
}

// bigONEBidAsk is used in parsing the BigONE API response only.
type bigONEBidAsk struct {
	Price      float64 `json:"price"`
	OrderCount int     `json:"order_count"`
	Quantity   float64 `json:"quantity"`
}

// bigONEPubTickerData is used in parsing the BigONE API response only.
type bigONEPubTickerData struct {
	Code int
	Data bigONEPubTickerInnerDataData
}

// bigONEPubTickerInnerDataData is used in parsing the BigONE API response only.
type bigONEPubTickerInnerDataData struct {
	AssetPairName string
	Bid           bigONEBidAsk
	Ask           bigONEBidAsk
	Open          float64
	High          float64
	Low           float64
	Close         float64
	Volume        float64
	DailyChange   float64
}

// Normalize parses the fields in bigONEPubTickerResp and returns a
// bigONEPubTickerData with proper data types.
func (resp *bigONEPubTickerResp) Normalize() (*bigONEPubTickerData, error) {
	bidPrice, err := strconv.ParseFloat(resp.Data.Bid.Price, 64)
	if err != nil {
		return nil, err
	}
	bidQuantity, err := strconv.ParseFloat(resp.Data.Bid.Quantity, 64)
	if err != nil {
		return nil, err
	}
	askPrice, err := strconv.ParseFloat(resp.Data.Ask.Price, 64)
	if err != nil {
		return nil, err
	}
	askQuantity, err := strconv.ParseFloat(resp.Data.Ask.Quantity, 64)
	if err != nil {
		return nil, err
	}

	_close, err := strconv.ParseFloat(resp.Data.Close, 64)
	if err != nil {
		return nil, err
	}
	high, err := strconv.ParseFloat(resp.Data.High, 64)
	if err != nil {
		return nil, err
	}
	low, err := strconv.ParseFloat(resp.Data.Low, 64)
	if err != nil {
		return nil, err
	}
	open, err := strconv.ParseFloat(resp.Data.Open, 64)
	if err != nil {
		return nil, err
	}

	volume, err := strconv.ParseFloat(resp.Data.Volume, 64)
	if err != nil {
		return nil, err
	}

	dailyChange, err := strconv.ParseFloat(resp.Data.DailyChange, 64)
	if err != nil {
		return nil, err
	}

	return &bigONEPubTickerData{
		Code: resp.Code,
		Data: bigONEPubTickerInnerDataData{
			AssetPairName: resp.Data.AssetPairName,
			Bid: bigONEBidAsk{
				Price:      bidPrice,
				OrderCount: resp.Data.Bid.OrderCount,
				Quantity:   bidQuantity,
			},
			Ask: bigONEBidAsk{
				Price:      askPrice,
				OrderCount: resp.Data.Ask.OrderCount,
				Quantity:   askQuantity,
			},
			Open:        open,
			High:        high,
			Low:         low,
			Close:       _close,
			Volume:      volume,
			DailyChange: dailyChange,
		},
	}, nil
}
