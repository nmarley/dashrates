package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// KrakenAPI implements the RateAPI interface and contains info necessary for
// calling to the public Kraken price ticker API.
type KrakenAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewKrakenAPI is a constructor for KrakenAPI.
func NewKrakenAPI() *KrakenAPI {
	return &KrakenAPI{
		BaseAPIURL:          "https://api.kraken.com",
		PriceTickerEndpoint: "/0/public/Ticker?pair=DASHUSD",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *KrakenAPI) DisplayName() string {
	return "Kraken"
}

// FetchRate gets the Dash exchange rate from the Kraken API.
//
// This is part of the RateAPI interface implementation.
func (a *KrakenAPI) FetchRate() (*RateInfo, error) {
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
	var res krakenTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	// Get a struct w/proper data types
	data, err := res.Result.DashUSDPair.Normalize()
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "USD",
		LastPrice:       data.LastClosed.Price,
		BaseAssetVolume: data.Volume.Today,
		FetchTime:       now,
	}

	return &ri, nil
}

// krakenAPIResult is only used for parsing the Kraken API response.
type krakenAPIResult struct {
	Ask        []string `json:"a"`
	Bid        []string `json:"b"`
	LastClosed []string `json:"c"`
	Volume     []string `json:"v"`
	VWAP       []string `json:"p"`
	Trades     []int    `json:"t"`
	Low        []string `json:"l"`
	High       []string `json:"h"`
	Open       string   `json:"o"`
}

// Normalize ... does the needful.
func (resp *krakenAPIResult) Normalize() (*krakenResult, error) {
	askArr := make([]float64, 3)
	for i := 0; i < 3; i++ {
		x, err := strconv.ParseFloat(resp.Ask[i], 64)
		if err != nil {
			return nil, err
		}
		askArr[i] = x
	}

	bidArr := make([]float64, 3)
	for i := 0; i < 3; i++ {
		x, err := strconv.ParseFloat(resp.Bid[i], 64)
		if err != nil {
			return nil, err
		}
		bidArr[i] = x
	}

	lastClosedArr := make([]float64, 2)
	for i := 0; i < 2; i++ {
		x, err := strconv.ParseFloat(resp.LastClosed[i], 64)
		if err != nil {
			return nil, err
		}
		lastClosedArr[i] = x
	}

	vArr := make([]float64, 2)
	for i := 0; i < 2; i++ {
		x, err := strconv.ParseFloat(resp.Volume[i], 64)
		if err != nil {
			return nil, err
		}
		vArr[i] = x
	}

	vwapArr := make([]float64, 2)
	for i := 0; i < 2; i++ {
		x, err := strconv.ParseFloat(resp.VWAP[i], 64)
		if err != nil {
			return nil, err
		}
		vwapArr[i] = x
	}

	lowArr := make([]float64, 2)
	for i := 0; i < 2; i++ {
		x, err := strconv.ParseFloat(resp.Low[i], 64)
		if err != nil {
			return nil, err
		}
		lowArr[i] = x
	}

	highArr := make([]float64, 2)
	for i := 0; i < 2; i++ {
		x, err := strconv.ParseFloat(resp.High[i], 64)
		if err != nil {
			return nil, err
		}
		highArr[i] = x
	}

	open, err := strconv.ParseFloat(resp.Open, 64)
	if err != nil {
		return nil, err
	}

	return &krakenResult{
		Ask: bidAskArray{
			Price:          askArr[0],
			WholeLotVolume: askArr[1],
			LotVolume:      askArr[2],
		},
		Bid: bidAskArray{
			Price:          bidArr[0],
			WholeLotVolume: bidArr[1],
			LotVolume:      bidArr[2],
		},

		LastClosed: lastTradeClosedArray{
			Price:     lastClosedArr[0],
			LotVolume: lastClosedArr[1],
		},

		Volume: dailyArray{
			Today:       vArr[0],
			Last24Hours: vArr[1],
		},

		VWAP: dailyArray{
			Today:       vwapArr[0],
			Last24Hours: vwapArr[1],
		},

		Trades: dailyTradesArray{
			Today:       resp.Trades[0],
			Last24Hours: resp.Trades[1],
		},

		Low: dailyArray{
			Today:       lowArr[0],
			Last24Hours: lowArr[1],
		},

		High: dailyArray{
			Today:       highArr[0],
			Last24Hours: highArr[1],
		},

		Open: open,
	}, nil
}

// krakenResult is used in parsing the Kraken API response only.
//
// This contains the parsed fields with correct data types. It is created from
// parsing the Kraken API result.
type krakenResult struct {
	Ask        bidAskArray
	Bid        bidAskArray
	LastClosed lastTradeClosedArray
	Volume     dailyArray
	VWAP       dailyArray
	Trades     dailyTradesArray
	Low        dailyArray
	High       dailyArray
	Open       float64
}

// bidAskArray is used in parsing the Kraken API response only.
type bidAskArray struct {
	Price          float64
	WholeLotVolume float64
	LotVolume      float64
}

// lastTradeClosedArray is used in parsing the Kraken API response only.
type lastTradeClosedArray struct {
	Price     float64
	LotVolume float64
}

// dailyArray is used in parsing the Kraken API response only.
type dailyArray struct {
	Today       float64
	Last24Hours float64
}

// dailyTradesArray is used in parsing the Kraken API response only.
type dailyTradesArray struct {
	Today       int
	Last24Hours int
}

// krakenDashResult is only used for parsing the Kraken API response.
type krakenDashResult struct {
	DashUSDPair krakenAPIResult `json:"DASHUSD"`
}

// krakenTickerResp is only used for parsing the Kraken API response.
type krakenTickerResp struct {
	Errors []string         `json:"error"`
	Result krakenDashResult `json:"result"`
}
