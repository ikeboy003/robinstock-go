package orders

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/robinstock_go"
	"github.com/yourusername/robinstock_go/profiles"
	"github.com/yourusername/robinstock_go/stocks"
	"github.com/yourusername/robinstock_go/urls"
	"github.com/yourusername/robinstock_go/utils"
)

// GetAllStockOrders returns all stock orders for an account.
func GetAllStockOrders(ctx context.Context, client *robinstock_go.Client, accountNumber, startDate *string) ([]map[string]interface{}, error) {
	log.Println("GetAllStockOrders: Fetching stock orders...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := urls.OrdersURL(nil, accountNumber, startDate)
	results, err := client.FetchAllPages(ctx, url, true)
	if err != nil {
		log.Printf("GetAllStockOrders: Error: %v\n", err)
		return nil, err
	}

	log.Printf("GetAllStockOrders: Retrieved %d orders\n", len(results))
	return results, nil
}

// GetAllOpenStockOrders returns all open stock orders.
func GetAllOpenStockOrders(ctx context.Context, client *robinstock_go.Client, accountNumber *string) ([]map[string]interface{}, error) {
	log.Println("GetAllOpenStockOrders: Fetching open stock orders...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := urls.OrdersURL(nil, accountNumber, nil)
	results, err := client.FetchAllPages(ctx, url, true)
	if err != nil {
		log.Printf("GetAllOpenStockOrders: Error: %v\n", err)
		return nil, err
	}

	var openOrders []map[string]interface{}
	for _, order := range results {
		if utils.GetString(order, "cancel") != "" {
			openOrders = append(openOrders, order)
		}
	}

	log.Printf("GetAllOpenStockOrders: Retrieved %d open orders\n", len(openOrders))
	return openOrders, nil
}

// GetStockOrderInfo returns information for a specific stock order.
func GetStockOrderInfo(ctx context.Context, client *robinstock_go.Client, orderID string) (map[string]interface{}, error) {
	log.Printf("GetStockOrderInfo: Fetching order %s...\n", orderID)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := urls.OrdersURL(&orderID, nil, nil)
	resp, err := client.Get(ctx, url, nil, true)
	if err != nil {
		log.Printf("GetStockOrderInfo: Error: %v\n", err)
		return nil, err
	}

	return resp.Data, nil
}

// CancelStockOrder cancels a specific stock order.
func CancelStockOrder(ctx context.Context, client *robinstock_go.Client, orderID string) (map[string]interface{}, error) {
	log.Printf("CancelStockOrder: Cancelling order %s...\n", orderID)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := urls.CancelURL(orderID)
	resp, err := client.Post(ctx, url, nil, true)
	if err != nil {
		log.Printf("CancelStockOrder: Error: %v\n", err)
		return nil, err
	}

	log.Printf("CancelStockOrder: Order %s cancelled\n", orderID)
	return resp.Data, nil
}

// CancelAllStockOrders cancels all open stock orders.
func CancelAllStockOrders(ctx context.Context, client *robinstock_go.Client, accountNumber *string) ([]map[string]interface{}, error) {
	log.Println("CancelAllStockOrders: Cancelling all open stock orders...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	openOrders, err := GetAllOpenStockOrders(ctx, client, accountNumber)
	if err != nil {
		return nil, err
	}

	var cancelledOrders []map[string]interface{}
	for _, order := range openOrders {
		cancelURL := utils.GetString(order, "cancel")
		if cancelURL != "" {
			resp, err := client.Post(ctx, cancelURL, nil, true)
			if err == nil {
				cancelledOrders = append(cancelledOrders, resp.Data)
			}
		}
	}

	log.Printf("CancelAllStockOrders: Cancelled %d orders\n", len(cancelledOrders))
	return cancelledOrders, nil
}

// OrderBuyMarket submits a market buy order.
func OrderBuyMarket(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderBuyMarket: Buying %f shares of %s...\n", quantity, symbol)
	return placeOrder(ctx, client, symbol, quantity, "buy", nil, nil, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderBuyLimit submits a limit buy order.
func OrderBuyLimit(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, limitPrice float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderBuyLimit: Buying %f shares of %s at limit $%f...\n", quantity, symbol, limitPrice)
	return placeOrder(ctx, client, symbol, quantity, "buy", &limitPrice, nil, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderBuyStopLoss submits a stop loss buy order.
func OrderBuyStopLoss(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, stopPrice float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderBuyStopLoss: Buying %f shares of %s at stop $%f...\n", quantity, symbol, stopPrice)
	return placeOrder(ctx, client, symbol, quantity, "buy", nil, &stopPrice, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderBuyStopLimit submits a stop limit buy order.
func OrderBuyStopLimit(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, limitPrice, stopPrice float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderBuyStopLimit: Buying %f shares of %s at limit $%f, stop $%f...\n", quantity, symbol, limitPrice, stopPrice)
	return placeOrder(ctx, client, symbol, quantity, "buy", &limitPrice, &stopPrice, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderSellMarket submits a market sell order.
func OrderSellMarket(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderSellMarket: Selling %f shares of %s...\n", quantity, symbol)
	return placeOrder(ctx, client, symbol, quantity, "sell", nil, nil, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderSellLimit submits a limit sell order.
func OrderSellLimit(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, limitPrice float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderSellLimit: Selling %f shares of %s at limit $%f...\n", quantity, symbol, limitPrice)
	return placeOrder(ctx, client, symbol, quantity, "sell", &limitPrice, nil, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderSellStopLoss submits a stop loss sell order.
func OrderSellStopLoss(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, stopPrice float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderSellStopLoss: Selling %f shares of %s at stop $%f...\n", quantity, symbol, stopPrice)
	return placeOrder(ctx, client, symbol, quantity, "sell", nil, &stopPrice, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderSellStopLimit submits a stop limit sell order.
func OrderSellStopLimit(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, limitPrice, stopPrice float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderSellStopLimit: Selling %f shares of %s at limit $%f, stop $%f...\n", quantity, symbol, limitPrice, stopPrice)
	return placeOrder(ctx, client, symbol, quantity, "sell", &limitPrice, &stopPrice, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderBuyFractionalByQuantity submits a fractional share buy order by quantity.
func OrderBuyFractionalByQuantity(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderBuyFractionalByQuantity: Buying %f fractional shares of %s...\n", quantity, symbol)
	return placeOrder(ctx, client, symbol, quantity, "buy", nil, nil, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderBuyFractionalByPrice submits a fractional share buy order by dollar amount.
func OrderBuyFractionalByPrice(ctx context.Context, client *robinstock_go.Client, symbol string, amountInDollars float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderBuyFractionalByPrice: Buying $%f worth of %s...\n", amountInDollars, symbol)

	if amountInDollars < 1 {
		return nil, fmt.Errorf("fractional share price should meet minimum $1.00")
	}

	priceType := "ask_price"
	prices, err := stocks.GetLatestPrice(ctx, client, &priceType, extendedHours, symbol)
	if err != nil || len(prices) == 0 {
		return nil, fmt.Errorf("failed to get latest price: %w", err)
	}

	price := utils.ParseFloat(prices[0])
	if price == 0 {
		return nil, fmt.Errorf("price is zero, unable to calculate fractional shares")
	}

	fractionalShares := utils.RoundPrice(amountInDollars / price)
	return placeOrder(ctx, client, symbol, fractionalShares, "buy", nil, nil, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderSellFractionalByQuantity submits a fractional share sell order by quantity.
func OrderSellFractionalByQuantity(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderSellFractionalByQuantity: Selling %f fractional shares of %s...\n", quantity, symbol)
	return placeOrder(ctx, client, symbol, quantity, "sell", nil, nil, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderSellFractionalByPrice submits a fractional share sell order by dollar amount.
func OrderSellFractionalByPrice(ctx context.Context, client *robinstock_go.Client, symbol string, amountInDollars float64, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderSellFractionalByPrice: Selling $%f worth of %s...\n", amountInDollars, symbol)

	if amountInDollars < 1 {
		return nil, fmt.Errorf("fractional share price should meet minimum $1.00")
	}

	priceType := "bid_price"
	prices, err := stocks.GetLatestPrice(ctx, client, &priceType, extendedHours, symbol)
	if err != nil || len(prices) == 0 {
		return nil, fmt.Errorf("failed to get latest price: %w", err)
	}

	price := utils.ParseFloat(prices[0])
	if price == 0 {
		return nil, fmt.Errorf("price is zero, unable to calculate fractional shares")
	}

	fractionalShares := utils.RoundPrice(amountInDollars / price)
	return placeOrder(ctx, client, symbol, fractionalShares, "sell", nil, nil, accountNumber, timeInForce, extendedHours, "regular_hours", nil, "")
}

// OrderTrailingStop submits a trailing stop order.
func OrderTrailingStop(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, side string, trailAmount float64, trailType string, accountNumber *string, timeInForce string, extendedHours bool) (map[string]interface{}, error) {
	log.Printf("OrderTrailingStop: %s %f shares of %s with trail %f (%s)...\n", side, quantity, symbol, trailAmount, trailType)
	return placeOrder(ctx, client, symbol, quantity, side, nil, nil, accountNumber, timeInForce, extendedHours, "regular_hours", &trailAmount, trailType)
}

func placeOrder(ctx context.Context, client *robinstock_go.Client, symbol string, quantity float64, side string, limitPrice, stopPrice *float64, accountNumber *string, timeInForce string, extendedHours bool, marketHours string, trailAmount *float64, trailType string) (map[string]interface{}, error) {
	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	orderType := "market"
	trigger := "immediate"

	var priceType string
	if side == "buy" {
		priceType = "ask_price"
	} else {
		priceType = "bid_price"
	}

	var price float64
	if limitPrice != nil && stopPrice != nil {
		price = utils.RoundPrice(*limitPrice)
		*stopPrice = utils.RoundPrice(*stopPrice)
		orderType = "limit"
		trigger = "stop"
	} else if limitPrice != nil {
		price = utils.RoundPrice(*limitPrice)
		orderType = "limit"
	} else if stopPrice != nil {
		*stopPrice = utils.RoundPrice(*stopPrice)
		if side == "buy" {
			price = *stopPrice
		}
		trigger = "stop"
	} else {
		prices, err := stocks.GetLatestPrice(ctx, client, &priceType, extendedHours, symbol)
		if err != nil || len(prices) == 0 {
			return nil, fmt.Errorf("failed to get latest price: %w", err)
		}
		price = utils.RoundPrice(utils.ParseFloat(prices[0]))
	}

	accountURL, err := getAccountURL(ctx, client, accountNumber)
	if err != nil {
		return nil, err
	}

	instrumentURL, err := getInstrumentURL(ctx, client, symbol)
	if err != nil {
		return nil, err
	}

	askPrices, _ := stocks.GetLatestPrice(ctx, client, utils.Address("ask_price"), extendedHours, symbol)
	bidPrices, _ := stocks.GetLatestPrice(ctx, client, utils.Address("bid_price"), extendedHours, symbol)

	payload := map[string]interface{}{
		"account":            accountURL,
		"instrument":         instrumentURL,
		"symbol":             symbol,
		"price":              price,
		"ask_price":          utils.RoundPrice(utils.ParseFloat(askPrices[0])),
		"bid_price":          utils.RoundPrice(utils.ParseFloat(bidPrices[0])),
		"bid_ask_timestamp":  time.Now().Format("2006-01-02 15:04:05.000000"),
		"quantity":           quantity,
		"ref_id":             uuid.NewString(),
		"type":               orderType,
		"time_in_force":      timeInForce,
		"trigger":            trigger,
		"side":               side,
		"market_hours":       marketHours,
		"extended_hours":     extendedHours,
		"order_form_version": 4,
	}

	if stopPrice != nil {
		payload["stop_price"] = *stopPrice
	}

	if orderType == "market" && trigger != "stop" {
		delete(payload, "stop_price")
	}

	if marketHours == "regular_hours" {
		if side == "buy" {
			payload["preset_percent_limit"] = "0.05"
			payload["type"] = "limit"
		} else if orderType == "market" && side == "sell" {
			delete(payload, "price")
		}
	} else if marketHours == "extended_hours" || marketHours == "all_day_hours" {
		payload["type"] = "limit"
		payload["quantity"] = int(quantity)
	}

	if trailAmount != nil {
		stockPrices, _ := stocks.GetLatestPrice(ctx, client, nil, extendedHours, symbol)
		stockPrice := utils.RoundPrice(utils.ParseFloat(stockPrices[0]))

		var margin float64
		var percentage float64

		if trailType == "amount" {
			margin = *trailAmount
		} else {
			margin = stockPrice * (*trailAmount) * 0.01
			percentage = *trailAmount
		}

		var calculatedStopPrice float64
		if side == "buy" {
			calculatedStopPrice = stockPrice + margin
		} else {
			calculatedStopPrice = stockPrice - margin
		}
		calculatedStopPrice = utils.RoundPrice(calculatedStopPrice)

		payload["stop_price"] = calculatedStopPrice
		payload["type"] = "market"
		payload["trigger"] = "stop"

		if side == "buy" {
			payload["price"] = utils.RoundPrice(calculatedStopPrice * 1.05)
		}

		if trailType == "amount" {
			payload["trailing_peg"] = map[string]interface{}{
				"type": "price",
				"price": map[string]interface{}{
					"amount":        *trailAmount,
					"currency_code": "USD",
				},
			}
		} else {
			payload["trailing_peg"] = map[string]interface{}{
				"type":       "percentage",
				"percentage": fmt.Sprintf("%f", percentage),
			}
		}
	}

	url := urls.OrdersURL(nil, accountNumber, nil)
	resp, err := client.Post(ctx, url, payload, true)
	if err != nil {
		log.Printf("placeOrder: Error: %v\n", err)
		return nil, err
	}

	log.Printf("placeOrder: Order placed successfully. Order ID: %s\n", utils.GetString(resp.Data, "id"))
	return resp.Data, nil
}

func getAccountURL(ctx context.Context, client *robinstock_go.Client, accountNumber *string) (string, error) {
	if accountNumber != nil {
		account, err := profiles.GetAccountProfile(ctx, client, *accountNumber)
		if err != nil {
			return "", err
		}
		return account.URL, nil
	}

	accounts, err := profiles.GetAllAccountProfiles(ctx, client)
	if err != nil {
		return "", err
	}

	if len(accounts) == 0 {
		return "", fmt.Errorf("no accounts found")
	}

	return accounts[0].URL, nil
}

func getInstrumentURL(ctx context.Context, client *robinstock_go.Client, symbol string) (string, error) {
	instrument, err := stocks.GetInstrumentBySymbol(ctx, client, symbol)
	if err != nil {
		return "", err
	}
	return instrument.URL, nil
}

