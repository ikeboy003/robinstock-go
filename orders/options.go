package orders

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/yourusername/robinstock_go"
	"github.com/yourusername/robinstock_go/urls"
	"github.com/yourusername/robinstock_go/utils"
)

// GetAllOptionOrders returns all option orders for an account.
func GetAllOptionOrders(ctx context.Context, client *robinstock_go.Client, accountNumber, startDate *string) ([]map[string]interface{}, error) {
	log.Println("GetAllOptionOrders: Fetching option orders...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := urls.OptionOrdersURL(nil, accountNumber, startDate)
	results, err := client.FetchAllPages(ctx, url, true)
	if err != nil {
		log.Printf("GetAllOptionOrders: Error: %v\n", err)
		return nil, err
	}

	log.Printf("GetAllOptionOrders: Retrieved %d orders\n", len(results))
	return results, nil
}

// GetAllOpenOptionOrders returns all open option orders.
func GetAllOpenOptionOrders(ctx context.Context, client *robinstock_go.Client, accountNumber *string) ([]map[string]interface{}, error) {
	log.Println("GetAllOpenOptionOrders: Fetching open option orders...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := urls.OptionOrdersURL(nil, accountNumber, nil)
	results, err := client.FetchAllPages(ctx, url, true)
	if err != nil {
		log.Printf("GetAllOpenOptionOrders: Error: %v\n", err)
		return nil, err
	}

	var openOrders []map[string]interface{}
	for _, order := range results {
		if utils.GetString(order, "cancel_url") != "" {
			openOrders = append(openOrders, order)
		}
	}

	log.Printf("GetAllOpenOptionOrders: Retrieved %d open orders\n", len(openOrders))
	return openOrders, nil
}

// GetOptionOrderInfo returns information for a specific option order.
func GetOptionOrderInfo(ctx context.Context, client *robinstock_go.Client, orderID string) (map[string]interface{}, error) {
	log.Printf("GetOptionOrderInfo: Fetching order %s...\n", orderID)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := urls.OptionOrdersURL(&orderID, nil, nil)
	resp, err := client.Get(ctx, url, nil, true)
	if err != nil {
		log.Printf("GetOptionOrderInfo: Error: %v\n", err)
		return nil, err
	}

	return resp.Data, nil
}

// CancelOptionOrder cancels a specific option order.
func CancelOptionOrder(ctx context.Context, client *robinstock_go.Client, orderID string) (map[string]interface{}, error) {
	log.Printf("CancelOptionOrder: Cancelling order %s...\n", orderID)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := urls.OptionCancelURL(orderID)
	resp, err := client.Post(ctx, url, nil, true)
	if err != nil {
		log.Printf("CancelOptionOrder: Error: %v\n", err)
		return nil, err
	}

	log.Printf("CancelOptionOrder: Order %s cancelled\n", orderID)
	return resp.Data, nil
}

// CancelAllOptionOrders cancels all open option orders.
func CancelAllOptionOrders(ctx context.Context, client *robinstock_go.Client, accountNumber *string) ([]map[string]interface{}, error) {
	log.Println("CancelAllOptionOrders: Cancelling all open option orders...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	openOrders, err := GetAllOpenOptionOrders(ctx, client, accountNumber)
	if err != nil {
		return nil, err
	}

	var cancelledOrders []map[string]interface{}
	for _, order := range openOrders {
		cancelURL := utils.GetString(order, "cancel_url")
		if cancelURL != "" {
			resp, err := client.Post(ctx, cancelURL, nil, true)
			if err == nil {
				cancelledOrders = append(cancelledOrders, resp.Data)
			}
		}
	}

	log.Printf("CancelAllOptionOrders: Cancelled %d orders\n", len(cancelledOrders))
	return cancelledOrders, nil
}

// GetAllOptionPositions returns all option positions ever held.
func GetAllOptionPositions(ctx context.Context, client *robinstock_go.Client, accountNumber *string) ([]map[string]interface{}, error) {
	log.Println("GetAllOptionPositions: Fetching all option positions...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := buildOptionPositionsURL(accountNumber)
	results, err := client.FetchAllPages(ctx, url, true)
	if err != nil {
		log.Printf("GetAllOptionPositions: Error: %v\n", err)
		return nil, err
	}

	log.Printf("GetAllOptionPositions: Retrieved %d positions\n", len(results))
	return results, nil
}

// GetOpenOptionPositions returns all open option positions.
func GetOpenOptionPositions(ctx context.Context, client *robinstock_go.Client, accountNumber *string) ([]map[string]interface{}, error) {
	log.Println("GetOpenOptionPositions: Fetching open option positions...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := buildOptionPositionsURL(accountNumber)
	params := map[string]string{"nonzero": "true"}

	resp, err := client.Get(ctx, url, params, true)
	if err != nil {
		log.Printf("GetOpenOptionPositions: Error: %v\n", err)
		return nil, err
	}

	var positions []map[string]interface{}
	if results, ok := resp.Data["results"].([]interface{}); ok {
		for _, item := range results {
			if pos, ok := item.(map[string]interface{}); ok {
				positions = append(positions, pos)
			}
		}
	}

	log.Printf("GetOpenOptionPositions: Retrieved %d positions\n", len(positions))
	return positions, nil
}

// GetOptionChains returns the option chain for a symbol.
func GetOptionChains(ctx context.Context, client *robinstock_go.Client, symbol string) (map[string]interface{}, error) {
	log.Printf("GetOptionChains: Fetching chains for %s...\n", symbol)

	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	url := fmt.Sprintf("https://api.robinhood.com/options/chains/?equity_instrument_ids=%s", symbol)
	resp, err := client.Get(ctx, url, nil, false)
	if err != nil {
		log.Printf("GetOptionChains: Error: %v\n", err)
		return nil, err
	}

	return resp.Data, nil
}

// GetOptionInstruments returns option instruments based on filters.
func GetOptionInstruments(ctx context.Context, client *robinstock_go.Client, chainID, expirationDate, strikePrice, optionType *string) ([]map[string]interface{}, error) {
	log.Println("GetOptionInstruments: Fetching option instruments...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := "https://api.robinhood.com/options/instruments/"
	params := make(map[string]string)

	if chainID != nil {
		params["chain_id"] = *chainID
	}
	if expirationDate != nil {
		params["expiration_dates"] = *expirationDate
	}
	if strikePrice != nil {
		params["strike_price"] = *strikePrice
	}
	if optionType != nil {
		params["type"] = *optionType
	}

	results, err := client.FetchAllPages(ctx, url, true)
	if err != nil {
		log.Printf("GetOptionInstruments: Error: %v\n", err)
		return nil, err
	}

	log.Printf("GetOptionInstruments: Retrieved %d instruments\n", len(results))
	return results, nil
}

// OrderOptionBuyLimit places a limit buy order for an option.
func OrderOptionBuyLimit(ctx context.Context, client *robinstock_go.Client, positionEffect, creditOrDebit string, price float64, symbol string, quantity int, expirationDate, strike, optionType string, accountNumber *string, timeInForce string) (map[string]interface{}, error) {
	log.Printf("OrderOptionBuyLimit: Buying %d %s %s $%s %s...\n", quantity, symbol, expirationDate, strike, optionType)
	return placeOptionOrder(ctx, client, "buy", positionEffect, creditOrDebit, price, 0, symbol, quantity, expirationDate, strike, optionType, accountNumber, timeInForce)
}

// OrderOptionSellLimit places a limit sell order for an option.
func OrderOptionSellLimit(ctx context.Context, client *robinstock_go.Client, positionEffect, creditOrDebit string, price float64, symbol string, quantity int, expirationDate, strike, optionType string, accountNumber *string, timeInForce string) (map[string]interface{}, error) {
	log.Printf("OrderOptionSellLimit: Selling %d %s %s $%s %s...\n", quantity, symbol, expirationDate, strike, optionType)
	return placeOptionOrder(ctx, client, "sell", positionEffect, creditOrDebit, price, 0, symbol, quantity, expirationDate, strike, optionType, accountNumber, timeInForce)
}

// OrderOptionSpread places an option spread order.
func OrderOptionSpread(ctx context.Context, client *robinstock_go.Client, direction string, price float64, symbol string, quantity int, spread []map[string]interface{}, accountNumber *string, timeInForce string) (map[string]interface{}, error) {
	log.Printf("OrderOptionSpread: Placing %s spread for %s...\n", direction, symbol)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	accountURL, err := getAccountURL(ctx, client, accountNumber)
	if err != nil {
		return nil, err
	}

	var legs []map[string]interface{}
	for _, leg := range spread {
		optionID, err := getOptionID(ctx, client, symbol,
			utils.GetString(leg, "expirationDate"),
			utils.GetString(leg, "strike"),
			utils.GetString(leg, "optionType"))
		if err != nil {
			return nil, err
		}

		legs = append(legs, map[string]interface{}{
			"position_effect": utils.GetString(leg, "effect"),
			"side":            utils.GetString(leg, "action"),
			"ratio_quantity":  utils.GetInt(leg, "ratio_quantity"),
			"option":          fmt.Sprintf("https://api.robinhood.com/options/instruments/%s/", optionID),
		})
	}

	payload := map[string]interface{}{
		"account":                   accountURL,
		"direction":                 direction,
		"time_in_force":             timeInForce,
		"legs":                      legs,
		"type":                      "limit",
		"trigger":                   "immediate",
		"price":                     utils.RoundPrice(price),
		"quantity":                  quantity,
		"override_day_trade_checks": false,
		"override_dtbp_checks":      false,
		"ref_id":                    uuid.NewString(),
	}

	url := urls.OptionOrdersURL(nil, accountNumber, nil)
	resp, err := client.Post(ctx, url, payload, true)
	if err != nil {
		log.Printf("OrderOptionSpread: Error: %v\n", err)
		return nil, err
	}

	log.Printf("OrderOptionSpread: Order placed successfully. Order ID: %s\n", utils.GetString(resp.Data, "id"))
	return resp.Data, nil
}

func placeOptionOrder(ctx context.Context, client *robinstock_go.Client, side, positionEffect, creditOrDebit string, price, stopPrice float64, symbol string, quantity int, expirationDate, strike, optionType string, accountNumber *string, timeInForce string) (map[string]interface{}, error) {
	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	accountURL, err := getAccountURL(ctx, client, accountNumber)
	if err != nil {
		return nil, err
	}

	optionID, err := getOptionID(ctx, client, symbol, expirationDate, strike, optionType)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"account":       accountURL,
		"direction":     creditOrDebit,
		"time_in_force": timeInForce,
		"legs": []map[string]interface{}{
			{
				"position_effect": positionEffect,
				"side":            side,
				"ratio_quantity":  1,
				"option":          fmt.Sprintf("https://api.robinhood.com/options/instruments/%s/", optionID),
			},
		},
		"type":                      "limit",
		"trigger":                   "immediate",
		"price":                     utils.RoundPrice(price),
		"quantity":                  quantity,
		"override_day_trade_checks": false,
		"override_dtbp_checks":      false,
		"ref_id":                    uuid.NewString(),
	}

	if stopPrice > 0 {
		payload["trigger"] = "stop"
		payload["stop_price"] = utils.RoundPrice(stopPrice)
	}

	url := urls.OptionOrdersURL(nil, accountNumber, nil)
	resp, err := client.Post(ctx, url, payload, true)
	if err != nil {
		log.Printf("placeOptionOrder: Error: %v\n", err)
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}

	if resp.StatusCode >= 400 {
		errMsg := utils.GetString(resp.Data, "detail")
		if errMsg == "" {
			errMsg = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
		log.Printf("placeOptionOrder: Error: %s\n", errMsg)
		return nil, fmt.Errorf("%s", errMsg)
	}

	if resp.Data == nil {
		return nil, fmt.Errorf("response data is nil")
	}

	log.Printf("placeOptionOrder: Order placed successfully. Order ID: %s\n", utils.GetString(resp.Data, "id"))
	return resp.Data, nil
}

func getOptionID(ctx context.Context, client *robinstock_go.Client, symbol, expirationDate, strike, optionType string) (string, error) {
	url := "https://api.robinhood.com/options/instruments/"
	params := map[string]string{
		"chain_symbol":     symbol,
		"expiration_dates": expirationDate,
		"strike_price":     strike,
		"type":             optionType,
		"state":            "active",
	}

	resp, err := client.Get(ctx, url, params, true)
	if err != nil {
		return "", err
	}

	if results, ok := resp.Data["results"].([]interface{}); ok && len(results) > 0 {
		if firstResult, ok := results[0].(map[string]interface{}); ok {
			return utils.GetString(firstResult, "id"), nil
		}
	}

	return "", fmt.Errorf("option not found for %s %s %s %s", symbol, expirationDate, strike, optionType)
}

func buildOptionPositionsURL(accountNumber *string) string {
	url := "https://api.robinhood.com/options/positions/"
	if accountNumber != nil {
		url += "?account_numbers=" + *accountNumber
	}
	return url
}

