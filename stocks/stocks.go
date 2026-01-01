package stocks

import (
	"context"
	"fmt"

	"github.com/yourusername/robinstock_go"
	"github.com/yourusername/robinstock_go/models"
	"github.com/yourusername/robinstock_go/utils"
)

// GetInstrumentBySymbol returns instrument data for a symbol.
func GetInstrumentBySymbol(ctx context.Context, client *robinstock_go.Client, symbol string) (*models.Instrument, error) {
	symbol = robinstock_go.NormalizeSymbol(symbol)

	params := map[string]string{"symbol": symbol}
	resp, err := client.Get(ctx, models.BaseURL+"/instruments/", params, false)
	if err != nil {
		return nil, err
	}

	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("instrument not found for symbol: %s", symbol)
	}

	result := resp.Results[0]
	instrument := &models.Instrument{
		ID:           utils.GetString(result, "id"),
		URL:          utils.GetString(result, "url"),
		Symbol:       utils.GetString(result, "symbol"),
		Name:         utils.GetString(result, "name"),
		SimpleName:   utils.GetString(result, "simple_name"),
		ListDate:     utils.GetString(result, "list_date"),
		Country:      utils.GetString(result, "country"),
		Type:         utils.GetString(result, "type"),
		Tradeable:    utils.GetBool(result, "tradeable"),
		State:        utils.GetString(result, "state"),
		Fundamentals: utils.GetString(result, "fundamentals"),
		Quote:        utils.GetString(result, "quote"),
		Market:       utils.GetString(result, "market"),
	}

	return instrument, nil
}

// GetInstrumentsBySymbols returns instrument data for multiple symbols.
func GetInstrumentsBySymbols(ctx context.Context, client *robinstock_go.Client, symbols []string) ([]models.Instrument, error) {
	symbols = robinstock_go.NormalizeSymbols(symbols)
	symbolsParam := robinstock_go.JoinSymbols(symbols)

	params := map[string]string{"symbols": symbolsParam}
	resp, err := client.Get(ctx, models.BaseURL+"/instruments/", params, false)
	if err != nil {
		return nil, err
	}

	var instruments []models.Instrument
	for _, result := range resp.Results {
		instrument := models.Instrument{
			ID:           utils.GetString(result, "id"),
			URL:          utils.GetString(result, "url"),
			Symbol:       utils.GetString(result, "symbol"),
			Name:         utils.GetString(result, "name"),
			SimpleName:   utils.GetString(result, "simple_name"),
			ListDate:     utils.GetString(result, "list_date"),
			Country:      utils.GetString(result, "country"),
			Type:         utils.GetString(result, "type"),
			Tradeable:    utils.GetBool(result, "tradeable"),
			State:        utils.GetString(result, "state"),
			Fundamentals: utils.GetString(result, "fundamentals"),
			Quote:        utils.GetString(result, "quote"),
			Market:       utils.GetString(result, "market"),
		}
		instruments = append(instruments, instrument)
	}

	return instruments, nil
}

// GetQuote returns a stock quote for a single symbol.
func GetQuote(ctx context.Context, client *robinstock_go.Client, symbol string) (*models.Quote, error) {
	quotes, err := GetQuotes(ctx, client, symbol)
	if err != nil {
		return nil, err
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("quote not found for symbol: %s", symbol)
	}

	return &quotes[0], nil
}

// GetQuotes returns stock quotes for multiple symbols.
func GetQuotes(ctx context.Context, client *robinstock_go.Client, symbols ...string) ([]models.Quote, error) {
	symbols = robinstock_go.NormalizeSymbols(symbols)
	symbolsParam := robinstock_go.JoinSymbols(symbols)

	params := map[string]string{"symbols": symbolsParam}
	resp, err := client.Get(ctx, models.BaseURL+"/quotes/", params, false)
	if err != nil {
		return nil, err
	}

	var quotes []models.Quote
	for _, result := range resp.Results {
		quote := models.Quote{
			Symbol:                      utils.GetString(result, "symbol"),
			InstrumentID:                utils.GetString(result, "instrument_id"),
			AskPrice:                    utils.GetString(result, "ask_price"),
			AskSize:                     utils.GetInt(result, "ask_size"),
			BidPrice:                    utils.GetString(result, "bid_price"),
			BidSize:                     utils.GetInt(result, "bid_size"),
			LastTradePrice:              utils.GetString(result, "last_trade_price"),
			LastExtendedHoursTradePrice: utils.GetString(result, "last_extended_hours_trade_price"),
			PreviousClose:               utils.GetString(result, "previous_close"),
			AdjustedPreviousClose:       utils.GetString(result, "adjusted_previous_close"),
			TradingHalted:               utils.GetBool(result, "trading_halted"),
			HasTraded:                   utils.GetBool(result, "has_traded"),
			UpdatedAt:                   utils.GetString(result, "updated_at"),
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

// GetFundamentals returns fundamental data for a symbol.
func GetFundamentals(ctx context.Context, client *robinstock_go.Client, symbol string) (*models.Fundamental, error) {
	symbol = robinstock_go.NormalizeSymbol(symbol)

	url := fmt.Sprintf("%s/fundamentals/%s/", models.BaseURL, symbol)
	resp, err := client.Get(ctx, url, nil, false)
	if err != nil {
		return nil, err
	}

	fundamental := &models.Fundamental{
		Open:                utils.GetString(resp.Data, "open"),
		High:                utils.GetString(resp.Data, "high"),
		Low:                 utils.GetString(resp.Data, "low"),
		Volume:              utils.GetString(resp.Data, "volume"),
		AverageVolume:       utils.GetString(resp.Data, "average_volume"),
		AverageVolume2Weeks: utils.GetString(resp.Data, "average_volume_2_weeks"),
		High52Weeks:         utils.GetString(resp.Data, "high_52_weeks"),
		Low52Weeks:          utils.GetString(resp.Data, "low_52_weeks"),
		MarketCap:           utils.GetString(resp.Data, "market_cap"),
		DividendYield:       utils.GetString(resp.Data, "dividend_yield"),
		PERatio:             utils.GetString(resp.Data, "pe_ratio"),
		Description:         utils.GetString(resp.Data, "description"),
		Instrument:          utils.GetString(resp.Data, "instrument"),
	}

	return fundamental, nil
}

// GetSymbolByURL returns the stock symbol for a given instrument URL.
func GetSymbolByURL(ctx context.Context, client *robinstock_go.Client, url string) (string, error) {
	resp, err := client.Get(ctx, url, nil, false)
	if err != nil {
		return "", err
	}
	return utils.GetString(resp.Data, "symbol"), nil
}

// GetHistoricals returns historical price data for a symbol.
func GetHistoricals(ctx context.Context, client *robinstock_go.Client, symbol string, interval string, span string) ([]models.HistoricalData, error) {
	symbol = robinstock_go.NormalizeSymbol(symbol)

	params := map[string]string{
		"symbols":  symbol,
		"interval": interval,
		"span":     span,
	}

	url := models.BaseURL + "/quotes/historicals/"
	resp, err := client.Get(ctx, url, params, false)
	if err != nil {
		return nil, err
	}

	resultsArray, ok := resp.Data["results"].([]interface{})
	if !ok || len(resultsArray) == 0 {
		return nil, fmt.Errorf("no historical data found")
	}

	firstResult, ok := resultsArray[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid historical data format")
	}

	historicalsArray, ok := firstResult["historicals"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no historicals array found")
	}

	var historicals []models.HistoricalData
	for _, item := range historicalsArray {
		if hist, ok := item.(map[string]interface{}); ok {
			historical := models.HistoricalData{
				OpenPrice:    utils.GetString(hist, "open_price"),
				ClosePrice:   utils.GetString(hist, "close_price"),
				HighPrice:    utils.GetString(hist, "high_price"),
				LowPrice:     utils.GetString(hist, "low_price"),
				Volume:       utils.GetInt(hist, "volume"),
				Session:      utils.GetString(hist, "session"),
				Interpolated: utils.GetBool(hist, "interpolated"),
			}
			historicals = append(historicals, historical)
		}
	}

	return historicals, nil
}

// GetRatings returns analyst ratings for a stock.
func GetRatings(ctx context.Context, client *robinstock_go.Client, symbol string) (map[string]interface{}, error) {
	symbol = robinstock_go.NormalizeSymbol(symbol)

	instrument, err := GetInstrumentBySymbol(ctx, client, symbol)
	if err != nil {
		return nil, err
	}

	url := models.BaseURL + "/midlands/ratings/" + instrument.ID + "/"
	resp, err := client.Get(ctx, url, nil, false)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// GetNews returns news articles for a stock symbol.
func GetNews(ctx context.Context, client *robinstock_go.Client, symbol string) ([]map[string]interface{}, error) {
	symbol = robinstock_go.NormalizeSymbol(symbol)

	url := models.BaseURL + "/midlands/news/" + symbol + "/"
	resp, err := client.Get(ctx, url, nil, false)
	if err != nil {
		return nil, err
	}

	var news []map[string]interface{}
	if results, ok := resp.Data["results"].([]interface{}); ok {
		for _, item := range results {
			if newsItem, ok := item.(map[string]interface{}); ok {
				news = append(news, newsItem)
			}
		}
	}

	return news, nil
}

// GetPopularity returns popularity data for a stock.
func GetPopularity(ctx context.Context, client *robinstock_go.Client, symbol string) (map[string]interface{}, error) {
	symbol = robinstock_go.NormalizeSymbol(symbol)

	instrument, err := GetInstrumentBySymbol(ctx, client, symbol)
	if err != nil {
		return nil, err
	}

	url := models.BaseURL + "/instruments/" + instrument.ID + "/popularity/"
	resp, err := client.Get(ctx, url, nil, false)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// GetSplits returns stock split history for a symbol.
func GetSplits(ctx context.Context, client *robinstock_go.Client, symbol string) ([]map[string]interface{}, error) {
	symbol = robinstock_go.NormalizeSymbol(symbol)

	instrument, err := GetInstrumentBySymbol(ctx, client, symbol)
	if err != nil {
		return nil, err
	}

	url := models.BaseURL + "/instruments/" + instrument.ID + "/splits/"
	results, err := client.FetchAllPages(ctx, url, false)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetLatestPrice returns the latest price for symbols.
func GetLatestPrice(ctx context.Context, client *robinstock_go.Client, priceType *string, includeExtendedHours bool, symbols ...string) ([]interface{}, error) {
	quotes, err := GetQuotes(ctx, client, symbols...)
	if err != nil {
		return nil, err
	}

	var prices []interface{}
	for _, quote := range quotes {
		if priceType != nil {
			switch *priceType {
			case "ask_price":
				prices = append(prices, quote.AskPrice)
			case "bid_price":
				prices = append(prices, quote.BidPrice)
			default:
				prices = append(prices, nil)
			}
		} else {
			if includeExtendedHours && quote.LastExtendedHoursTradePrice != "" {
				prices = append(prices, quote.LastExtendedHoursTradePrice)
			} else {
				prices = append(prices, quote.LastTradePrice)
			}
		}
	}

	return prices, nil
}
