package dashrates

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// OKExAPI implements the RateAPI interface and contains info necessary for
// calling to the public OKEx price ticker API.
type OKExAPI struct {
	BaseAPIURL          string
	PriceTickerEndpoint string
}

// NewOKExAPI is a constructor for OKExAPI.
func NewOKExAPI() *OKExAPI {
	return &OKExAPI{
		BaseAPIURL:          "https://www.okex.com",
		PriceTickerEndpoint: "/api/spot/v3/instruments/DASH-BTC/ticker",
	}
}

// DisplayName returns the exchange display name. It is part of the RateAPI
// interface implementation.
func (a *OKExAPI) DisplayName() string {
	return "OKEx"
}

// FetchRate gets the Dash exchange rate from the OKEx API.
//
// This is part of the RateAPI interface implementation.
func (a *OKExAPI) FetchRate() (*RateInfo, error) {
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
	var res okexPubTickerResp
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
		BaseAssetVolume: data.QuoteVolume24h,
		FetchTime:       now,
	}

	return &ri, nil
}

// okexTickerData has the output of the parsed OKExPubTickerResp and
// has proper data types.
type okexTickerData struct {
	BestAsk        float64
	BestBid        float64
	InstrumentID   string
	ProductID      string
	Last           float64
	LastQty        float64
	Ask            float64
	BestAskSize    float64
	Bid            float64
	BestBidSize    float64
	Open24h        float64
	High24h        float64
	Low24h         float64
	BaseVolume24h  float64
	Timestamp      time.Time
	QuoteVolume24h float64
}

// okexPubTickerResp is used in parsing the OKEx API response only.
type okexPubTickerResp struct {
	BestAsk        string `json:"best_ask"`
	BestBid        string `json:"best_bid"`
	InstrumentID   string `json:"instrument_id"`
	ProductID      string `json:"product_id"`
	Last           string `json:"last"`
	LastQty        string `json:"last_qty"`
	Ask            string `json:"ask"`
	BestAskSize    string `json:"best_ask_size"`
	Bid            string `json:"bid"`
	BestBidSize    string `json:"best_bid_size"`
	Open24h        string `json:"open_24h"`
	High24h        string `json:"high_24h"`
	Low24h         string `json:"low_24h"`
	BaseVolume24h  string `json:"base_volume_24h"`
	Timestamp      string `json:"timestamp"`
	QuoteVolume24h string `json:"quote_volume_24h"`
}

// Normalize parses the fields in OKExPubTickerResp and returns a
// okexTickerData with proper data types.
func (resp *okexPubTickerResp) Normalize() (*okexTickerData, error) {
	bestAsk, err := strconv.ParseFloat(resp.BestAsk, 64)
	if err != nil {
		return nil, err
	}

	bestBid, err := strconv.ParseFloat(resp.BestBid, 64)
	if err != nil {
		return nil, err
	}

	last, err := strconv.ParseFloat(resp.Last, 64)
	if err != nil {
		return nil, err
	}

	lastQty, err := strconv.ParseFloat(resp.LastQty, 64)
	if err != nil {
		return nil, err
	}

	ask, err := strconv.ParseFloat(resp.Ask, 64)
	if err != nil {
		return nil, err
	}

	bestAskSize, err := strconv.ParseFloat(resp.BestAskSize, 64)
	if err != nil {
		return nil, err
	}

	bid, err := strconv.ParseFloat(resp.Bid, 64)
	if err != nil {
		return nil, err
	}

	bestBidSize, err := strconv.ParseFloat(resp.BestBidSize, 64)
	if err != nil {
		return nil, err
	}

	open24h, err := strconv.ParseFloat(resp.Open24h, 64)
	if err != nil {
		return nil, err
	}

	high24h, err := strconv.ParseFloat(resp.High24h, 64)
	if err != nil {
		return nil, err
	}

	low24h, err := strconv.ParseFloat(resp.Low24h, 64)
	if err != nil {
		return nil, err
	}

	baseVolume24h, err := strconv.ParseFloat(resp.BaseVolume24h, 64)
	if err != nil {
		return nil, err
	}

	quoteVolume24h, err := strconv.ParseFloat(resp.QuoteVolume24h, 64)
	if err != nil {
		return nil, err
	}

	ts, err := time.Parse(time.RFC3339, resp.Timestamp)
	if err != nil {
		return nil, err
	}

	return &okexTickerData{
		BestAsk:        bestAsk,
		BestBid:        bestBid,
		InstrumentID:   resp.InstrumentID,
		ProductID:      resp.ProductID,
		Last:           last,
		LastQty:        lastQty,
		Ask:            ask,
		BestAskSize:    bestAskSize,
		Bid:            bid,
		BestBidSize:    bestBidSize,
		Open24h:        open24h,
		High24h:        high24h,
		Low24h:         low24h,
		BaseVolume24h:  baseVolume24h,
		Timestamp:      ts,
		QuoteVolume24h: quoteVolume24h,
	}, nil
}
