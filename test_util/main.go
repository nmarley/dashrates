package main

// Test utility to check each exchange API one-by-one. It can help debug
// exchange API routes which are no longer working and should be removed, or
// updated.

import (
	"fmt"
	"os"

	dashrates "github.com/dcginfra/dashrates"
)

func main() {
	// For each exchange rate API, try to pull the rate
	apis := []dashrates.RateAPI{
		dashrates.NewBiboxAPI(),
		dashrates.NewBigONEAPI(),
		dashrates.NewBinanceAPI(),
		dashrates.NewBitbnsAPI(),
		dashrates.NewBitfinexAPI(),
		dashrates.NewBittrexAPI(),
		dashrates.NewBvnexAPI(),
		dashrates.NewCexAPI(),
		dashrates.NewCoinCapAPI(),
		dashrates.NewCoinbaseAPI(),
		dashrates.NewCoinbaseProAPI(),
		dashrates.NewCrex24API(),
		dashrates.NewDigifinexAPI(),
		dashrates.NewExmoAPI(),
		dashrates.NewHitBTCAPI(),
		dashrates.NewHuobiAPI(),
		dashrates.NewIndodaxAPI(),
		dashrates.NewKrakenAPI(),
		dashrates.NewKuCoinAPI(),
		dashrates.NewLiquidAPI(),
		dashrates.NewOKExAPI(),
		dashrates.NewPoloniexAPI(),
		dashrates.NewSouthXchangeAPI(),
		dashrates.NewTrivAPI(),
		dashrates.NewUpholdAPI(),
		dashrates.NewWhiteBITAPI(),
		dashrates.NewYobitAPI(),
	}

	for _, api := range apis {
		_, err := api.FetchRate()
		if err != nil {
			// print err message to stderr
			fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
			// print a <exch> BAD message to stdout
			fmt.Fprintf(os.Stdout, "%v ERROR\n", api.DisplayName())
		} else {
			fmt.Fprintf(os.Stdout, "%v OK\n", api.DisplayName())
		}
	}
}
