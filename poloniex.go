package dashrates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// PoloniexAPI implements the RateAPI interface and contains info necessary for
// calling to the public Poloniex price ticker API.
type PoloniexAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewPoloniexAPI is a constructor for PoloniexAPI.
func NewPoloniexAPI() *PoloniexAPI {
	return &PoloniexAPI{
		BaseAPIURL:          "https://poloniex.com/public",
		PriceTickerEndpoint: "?command=returnTicker",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *PoloniexAPI) DisplayName() string {
	return "Poloniex"
}

// FetchRate gets the Dash exchange rate from the Poloniex API.
//
// This is part of the RateAPI interface implementation.
func (a *PoloniexAPI) FetchRate() (*RateInfo, error) {
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
	var res poloniexPubTickerResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	// Poloniex gets their base/quotes backwards - BTC is quote, DASH is base
	ticker, ok := res["BTC_DASH"]
	if !ok {
		err = fmt.Errorf("oh no, %s does not have DASH/BTC pair", a.DisplayName())
		return nil, err
	}
	data, err := ticker.Normalize()
	if err != nil {
		return nil, err
	}

	ri := RateInfo{
		BaseCurrency:    "DASH",
		QuoteCurrency:   "BTC",
		LastPrice:       data.Last,
		BaseAssetVolume: data.BaseVolume,
		FetchTime:       now,
	}

	return &ri, nil
}

// poloniexTickerData has the output of the parsed poloniexTickerPair struct
// and has proper data types.
type poloniexTickerData struct {
	ID            int
	Last          float64
	LowestAsk     float64
	HighestBid    float64
	PercentChange float64
	BaseVolume    float64
	QuoteVolume   float64
	IsFrozen      bool
	High24hr      float64
	Low24hr       float64
}

// poloniexPubTickerResp is used in parsing the Poloniex API response only.
type poloniexPubTickerResp map[string]poloniexTickerPair

// poloniexTickerPair is used in parsing the Poloniex API response only.
//
// Note that we flip quoteVolume & baseVolume b/c Poloniex gets that backwards
// too.
type poloniexTickerPair struct {
	ID            int    `json:"id"`
	Last          string `json:"last"`
	LowestAsk     string `json:"lowestAsk"`
	HighestBid    string `json:"highestBid"`
	PercentChange string `json:"percentChange"`
	BaseVolume    string `json:"quoteVolume"`
	QuoteVolume   string `json:"baseVolume"`
	IsFrozen      string `json:"isFrozen"`
	High24hr      string `json:"high24hr"`
	Low24hr       string `json:"low24hr"`
}

// Normalize parses the fields in poloniexTickerPair and returns a
// poloniexTickerData.
func (resp *poloniexTickerPair) Normalize() (*poloniexTickerData, error) {
	last, err := strconv.ParseFloat(resp.Last, 64)
	if err != nil {
		return nil, err
	}

	lowestAsk, err := strconv.ParseFloat(resp.LowestAsk, 64)
	if err != nil {
		return nil, err
	}
	highestBid, err := strconv.ParseFloat(resp.HighestBid, 64)
	if err != nil {
		return nil, err
	}
	percentChange, err := strconv.ParseFloat(resp.PercentChange, 64)
	if err != nil {
		return nil, err
	}
	baseVolume, err := strconv.ParseFloat(resp.BaseVolume, 64)
	if err != nil {
		return nil, err
	}
	quoteVolume, err := strconv.ParseFloat(resp.QuoteVolume, 64)
	if err != nil {
		return nil, err
	}
	high24hr, err := strconv.ParseFloat(resp.High24hr, 64)
	if err != nil {
		return nil, err
	}
	low24hr, err := strconv.ParseFloat(resp.Low24hr, 64)
	if err != nil {
		return nil, err
	}

	isFrozen, err := strconv.ParseInt(resp.IsFrozen, 10, 64)
	if err != nil {
		return nil, err
	}
	isFrozenBool := false
	if isFrozen == 0 {
		isFrozenBool = true
	}

	return &poloniexTickerData{
		ID:            resp.ID,
		Last:          last,
		LowestAsk:     lowestAsk,
		HighestBid:    highestBid,
		PercentChange: percentChange,
		BaseVolume:    baseVolume,
		QuoteVolume:   quoteVolume,
		IsFrozen:      isFrozenBool,
		High24hr:      high24hr,
		Low24hr:       low24hr,
	}, nil
}
