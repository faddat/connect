package coinbase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/oracle/types"
)

// GetSpotPriceEndpoint is the Coinbase endpoint for getting the spot price of a
// currency pair.
func getSpotPriceEndpoint(base, quote string) string {
	return fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-%s/spot", base, quote)
}

// GetSpotPriceHandler returns the spot price of a currency pair. In practice,
// this should not be used because price data should come from an aggregated
// price feed - API that uses a TWAP, TVWAP, or median price.
//
// Response format:
//
//	{
//	  "data": {
//	    "amount": "1020.25",
//	    "currency": "USD"
//	  }
//	}
func getPriceForPair(pair types.CurrencyPair) (*types.TickerPrice, error) {
	url := getSpotPriceEndpoint(pair.Base, pair.Quote)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respMap := make(map[string]map[string]string)
	if err := json.Unmarshal(body, &respMap); err != nil {
		return nil, err
	}

	data, ok := respMap["data"]
	if !ok {
		return nil, fmt.Errorf("failed to parse response")
	}

	amount, ok := data["amount"]
	if !ok {
		return nil, fmt.Errorf("failed to parse response")
	}

	sdkAmount, err := sdk.NewDecFromStr(amount)
	if err != nil {
		return nil, err
	}

	return &types.TickerPrice{
		Price:     sdkAmount,
		Timestamp: time.Now(),
	}, nil
}
