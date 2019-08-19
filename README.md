# Dash Exchange Rate Fetchers

> Fetch Dash exchange rate data from various exchanges

## Install

```sh
go get -u github.com/nmarley/dashrates
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/nmarley/dashrates"
)

func main() {
	api := dashrates.NewBinanceAPI()
	rate, err := api.FetchRate()
	if err != nil {
		panic(err)
	}

	fmt.Printf("rate info for %s: %+v\n", api.DisplayName(), rate)
}

// rate info for Binance: &{BaseCurrency:DASH QuoteCurrency:BTC LastPrice:0.008977 BaseAssetVolume:0 FetchTime:2019-08-19 16:03:48.054294 -0300 -03 m=+1.817687680}
```

## Contributing

Feel free to dive in! [Open an issue](https://github.com/nmarley/dashrates/issues/new) or submit PRs.

## License

[ISC](LICENSE)
