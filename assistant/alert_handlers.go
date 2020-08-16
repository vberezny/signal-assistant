package assistant

import (
	"fmt"
	"log"

	"github.com/coinpaprika/coinpaprika-api-go-client/v2/coinpaprika"
)

const (
	DEFAULT_CHANGE_6H  float64 = 2
	DEFAULT_CHANGE_24H float64 = 5
)

var tickerThresholds = map[string][]float64{
	"btc-bitcoin":  []float64{12000, 13000, 14000, 15000},
	"eth-ethereum": []float64{450, 500, 550, 600},
}

// Return an alert message if BTC or ETH either:
// a) break a pre-defined price threshhold.
// c) climb by 2% or more in the last 6 hours.
// d) climb by 5% or more in the last 24 hours.
// e) fall by 2% or more in the last hour.
// f) fall by 5% or more in the last 6 hours.
func cryptoAlertHandler() (message string) {
	paprikaClient := coinpaprika.NewClient(nil)
	for id, thresholds := range tickerThresholds {
		ticker, err := paprikaClient.Tickers.GetByID(id, nil)
		if err != nil {
			log.Print(err)
			return
		}
		if quoteUSD, ok := ticker.Quotes["USD"]; ok {
			price := *quoteUSD.Price
			for _, threshold := range thresholds {
				if price >= threshold {
					message += fmt.Sprintf("%v has crossed %v, it is now %.2f USD\n", *ticker.Symbol, threshold, price)
					break
				}
			}
			change6h := *quoteUSD.PercentChange6h
			change24h := *quoteUSD.PercentChange24h
			if change6h >= DEFAULT_CHANGE_6H {
				message += fmt.Sprintf("%v has gone up by %.2f in the last 6 hours, it is now %.2f USD\n", *ticker.Symbol, change6h, price)
			}
			if change6h <= -1*DEFAULT_CHANGE_6H {
				message += fmt.Sprintf("%v has gone down by %.2f in the last 6 hours, it is now %.2f USD\n", *ticker.Symbol, change6h, price)
			}
			if change24h >= DEFAULT_CHANGE_24H {
				message += fmt.Sprintf("%v has gone up by %.2f in the last 24 hours, it is now %.2f USD\n", *ticker.Symbol, change6h, price)
			}
			if change24h <= -1*DEFAULT_CHANGE_24H {
				message += fmt.Sprintf("%v has gone down by %.2f in the last 24 hours, it is now %.2f USD\n", *ticker.Symbol, change6h, price)
			}
		}
	}
	return
}
