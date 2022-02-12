# Dash Exchange Rate Fetchers

> Fetch Dash exchange rate data from various exchanges


## Contributing

Please run `make` to check code lint before committing

```sh
make
```

## Install

```sh
go get -u github.com/dcginfra/dashrates
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/dcginfra/dashrates"
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

## Test Utility

You can debug if exchanges are working or not by using the `test_util`:

```sh
cd test_util/
go build
./test_util 2>err | tee out
```

Sample output from test util:

```
Binance OK
Bitfinex OK
Coinbase OK
Cointrade ERROR
Kraken OK
Livecoin ERROR
Yobit OK
```

View the STDERR stream to see error detail:

```sh
# err file from above example
$ cat err

error fetching Cointrade: Market is not available.
error fetching Livecoin: Get "https://api.livecoin.net/exchange/ticker?currencyPair=DASH/USD": dial tcp: lookup api.livecoin.net: no such host
```

## Contributing

Feel free to dive in! [Open an issue](https://github.com/dcginfra/dashrates/issues/new) or submit PRs.

## License

[ISC](LICENSE)
