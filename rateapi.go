package dashrates

import (
	"encoding/json"
	"time"
)

// RateInfo contains information about exchange rates, including a Base and
// Quote pair, the last price, the asset volume in terms of the Base currency,
// and a fetch timestamp. Note that this timestamp is just for fetch time, and
// not an API server timestamp.
type RateInfo struct {
	BaseCurrency    string
	QuoteCurrency   string
	LastPrice       float64
	BaseAssetVolume float64
	FetchTime       time.Time
}

// MarshalBinary is part of the encoding.BinaryMarshaler interface
func (ri *RateInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(ri)
}

// UnmarshalBinary is part of the encoding.BinaryUnmarshaler interface
func (ri *RateInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, ri)
}

// RateAPI is an interface that describes an API which includes the Dash
// cryptocurrency.
type RateAPI interface {
	DisplayName() string
	FetchRate() (*RateInfo, error)
}
