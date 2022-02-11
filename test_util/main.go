package main

// Test utility to check each exchange API one-by-one. It can help debug
// exchange API routes which are no longer working and should be removed, or
// updated.
//
// TODO: Don't panic() below, instead print a good/bad message so all exchanges
// can be checked in just 1 run.

import (
	"fmt"
	"os"

	dashrates "github.com/dcginfra/dashrates"
)

func main() {
	// 1. Fetch BTC/USD rate
	coinCapAPI := dashrates.NewCoinCapAPI()
	coinCapRate, err := coinCapAPI.FetchRate()
	if err != nil {
		panic(err)
	}
	fmt.Printf("coinCapRate: %+v\n", coinCapRate)

	// now we have BTC/USD rate
	rateBitcoinUSD := coinCapRate.LastPrice
	fmt.Printf("rateBitcoinUSD: %+v\n", rateBitcoinUSD)

	// 2. For each exchange, pull the rate
	apis := []dashrates.RateAPI{
		dashrates.NewCrex24API(),
		dashrates.NewDigifinexAPI(),
		dashrates.NewCexAPI(),
		dashrates.NewBinanceAPI(),
		dashrates.NewBvnexAPI(),
		dashrates.NewUpholdAPI(),
		dashrates.NewKuCoinAPI(),
		dashrates.NewCoinbaseProAPI(),
		dashrates.NewBittrexAPI(),
		dashrates.NewBigONEAPI(),
		dashrates.NewCoinCapAPI(),
		dashrates.NewTrivAPI(),
		dashrates.NewIndodaxAPI(),
		dashrates.NewWhiteBITAPI(),
		dashrates.NewLiquidAPI(),
		dashrates.NewSouthXchangeAPI(),
		dashrates.NewBitfinexAPI(),
		dashrates.NewHuobiAPI(),
		dashrates.NewPoloniexAPI(),
		dashrates.NewBiboxAPI(),
		dashrates.NewExmoAPI(),
		dashrates.NewCoinbaseAPI(),
		dashrates.NewHitBTCAPI(),
		dashrates.NewOKExAPI(),
		dashrates.NewBitbnsAPI(),
		dashrates.NewKrakenAPI(),
		dashrates.NewYobitAPI(),
	}

	for _, api := range apis {
		rate, err := api.FetchRate()
		// TODO: Don't panic() here, print a good or bad message so all
		// exchanges can be checked in just 1 run.
		if err != nil {
			// maybe print this err message to stderr and a "<exch> BAD"
			// message to stdout
			fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
			panic(err)
		}
		fmt.Printf("exch: %v, rate: %+v\n", api.DisplayName(), rate)
		// usdRate, err := getDashRateInUSD(rateBitcoinUSD, api.DisplayName(), rate)
	}
}
