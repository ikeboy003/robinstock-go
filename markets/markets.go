package markets

import (
	"context"
	"log"
	"strings"

	"github.com/ikeboy003/robinstock-go"
	"github.com/ikeboy003/robinstock-go/models"
	"github.com/ikeboy003/robinstock-go/stocks"
	"github.com/ikeboy003/robinstock-go/urls"
	"github.com/ikeboy003/robinstock-go/utils"
)

// GetTopMoversSP500 returns top moving stocks in the S&P 500.
func GetTopMoversSP500(ctx context.Context, client *robinstock_go.Client, direction string) ([]models.SPMover, error) {
	log.Printf("GetTopMoversSP500: Fetching %s movers...\n", direction)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	direction = strings.ToLower(strings.TrimSpace(direction))
	if direction != "up" && direction != "down" {
		log.Println("GetTopMoversSP500: Direction must be 'up' or 'down'")
		return nil, nil
	}

	params := map[string]string{"direction": direction}
	resp, err := client.Get(ctx, urls.MoversSP500URL(), params, true)
	if err != nil {
		log.Printf("GetTopMoversSP500: Error: %v\n", err)
		return nil, err
	}

	var movers []models.SPMover
	if data, ok := resp.Data["results"].([]interface{}); ok {
		for _, item := range data {
			if moverMap, ok := item.(map[string]interface{}); ok {
				mover := models.SPMover{
					Description:                utils.GetString(moverMap, "description"),
					InstrumentURL:              utils.GetString(moverMap, "instrument_url"),
					MarketHoursLastMovementPct: getMovementPct(moverMap),
					Symbol:                     utils.GetString(moverMap, "symbol"),
				}
				movers = append(movers, mover)
			}
		}
	}

	log.Printf("GetTopMoversSP500: Received %d movers\n", len(movers))
	return movers, nil
}

// GetTop100MostPopular returns the 100 most popular stocks with their quotes.
func GetTop100MostPopular(ctx context.Context, client *robinstock_go.Client) ([]models.Quote, error) {
	log.Println("GetTop100MostPopular: Fetching top 100...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	resp, err := client.Get(ctx, urls.Top100MostPopularURL(), nil, true)
	if err != nil {
		log.Printf("GetTop100MostPopular: Error: %v\n", err)
		return nil, err
	}

	// Extract instrument URLs
	var instrumentURLs []string
	if instruments, ok := resp.Data["instruments"].([]interface{}); ok {
		for _, inst := range instruments {
			if url, ok := inst.(string); ok {
				instrumentURLs = append(instrumentURLs, url)
			}
		}
	}

	log.Printf("GetTop100MostPopular: Found %d instruments\n", len(instrumentURLs))

	// Get symbols from instrument URLs
	symbols := getSymbolsFromInstruments(ctx, client, instrumentURLs)

	if len(symbols) == 0 {
		return nil, nil
	}

	// Get quotes for all symbols
	quotes, err := stocks.GetQuotes(ctx, client, symbols...)
	if err != nil {
		return nil, err
	}

	log.Printf("GetTop100MostPopular: Received %d quotes\n", len(quotes))
	return quotes, nil
}

// GetStocksByMarketTag returns stocks from a specific market category.
func GetStocksByMarketTag(ctx context.Context, client *robinstock_go.Client, tag string) ([]models.Quote, error) {
	log.Printf("GetStocksByMarketTag: Fetching stocks for tag '%s'...\n", tag)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	resp, err := client.Get(ctx, urls.MarketCategoryURL(tag), nil, true)
	if err != nil {
		log.Printf("GetStocksByMarketTag: Error: %v\n", err)
		return nil, err
	}

	// Extract instrument URLs
	var instrumentURLs []string
	if instruments, ok := resp.Data["instruments"].([]interface{}); ok {
		for _, inst := range instruments {
			if url, ok := inst.(string); ok {
				instrumentURLs = append(instrumentURLs, url)
			}
		}
	}

	log.Printf("GetStocksByMarketTag: Found %d instruments\n", len(instrumentURLs))

	// Get symbols from instrument URLs
	symbols := getSymbolsFromInstruments(ctx, client, instrumentURLs)

	if len(symbols) == 0 {
		return nil, nil
	}

	// Get quotes for all symbols
	quotes, err := stocks.GetQuotes(ctx, client, symbols...)
	if err != nil {
		return nil, err
	}

	log.Printf("GetStocksByMarketTag: Received %d quotes\n", len(quotes))
	return quotes, nil
}

// GetMarkets returns all available markets.
func GetMarkets(ctx context.Context, client *robinstock_go.Client) ([]models.Market, error) {
	log.Println("GetMarkets: Fetching markets...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	results, err := client.FetchAllPages(ctx, urls.MarketsURL(), true)
	if err != nil {
		log.Printf("GetMarkets: Error: %v\n", err)
		return nil, err
	}

	var markets []models.Market
	for _, result := range results {
		market := models.Market{
			URL:          utils.GetString(result, "url"),
			TOTP:         utils.GetString(result, "totp"),
			Acronym:      utils.GetString(result, "acronym"),
			Name:         utils.GetString(result, "name"),
			City:         utils.GetString(result, "city"),
			Country:      utils.GetString(result, "country"),
			Timezone:     utils.GetString(result, "timezone"),
			Website:      utils.GetString(result, "website"),
			OperatingMIC: utils.GetString(result, "operating_mic"),
		}
		markets = append(markets, market)
	}

	log.Printf("GetMarkets: Received %d markets\n", len(markets))
	return markets, nil
}

// GetMarketHours returns trading hours for a specific market and date.
func GetMarketHours(ctx context.Context, client *robinstock_go.Client, market, date string) (*models.MarketHours, error) {
	log.Printf("GetMarketHours: Fetching hours for %s on %s...\n", market, date)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	resp, err := client.Get(ctx, urls.MarketHoursURL(market, date), nil, true)
	if err != nil {
		log.Printf("GetMarketHours: Error: %v\n", err)
		return nil, err
	}

	hours := &models.MarketHours{
		IsOpen:           utils.GetBool(resp.Data, "is_open"),
		OpensAt:          utils.GetString(resp.Data, "opens_at"),
		ClosesAt:         utils.GetString(resp.Data, "closes_at"),
		Date:             utils.GetString(resp.Data, "date"),
		ExtendedOpensAt:  utils.GetString(resp.Data, "extended_opens_at"),
		ExtendedClosesAt: utils.GetString(resp.Data, "extended_closes_at"),
	}

	log.Printf("GetMarketHours: Market is open: %v\n", hours.IsOpen)
	return hours, nil
}

func getMovementPct(moverMap map[string]interface{}) float64 {
	if priceMovement, ok := moverMap["price_movement"].(map[string]interface{}); ok {
		return utils.GetFloat(priceMovement, "market_hours_last_movement_pct")
	}
	return 0.0
}

// GetEarnings returns earnings reports for all stocks.
func GetEarnings(ctx context.Context, client *robinstock_go.Client) ([]map[string]interface{}, error) {
	log.Println("GetEarnings: Fetching earnings reports...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := models.BaseURL + "/marketdata/earnings/"
	results, err := client.FetchAllPages(ctx, url, true)
	if err != nil {
		log.Printf("GetEarnings: Error: %v\n", err)
		return nil, err
	}

	log.Printf("GetEarnings: Retrieved %d earnings reports\n", len(results))
	return results, nil
}

// GetEvents returns upcoming options events.
func GetEvents(ctx context.Context, client *robinstock_go.Client) ([]map[string]interface{}, error) {
	log.Println("GetEvents: Fetching market events...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := models.BaseURL + "/options/events/"
	results, err := client.FetchAllPages(ctx, url, true)
	if err != nil {
		log.Printf("GetEvents: Error: %v\n", err)
		return nil, err
	}

	log.Printf("GetEvents: Retrieved %d events\n", len(results))
	return results, nil
}

func getSymbolsFromInstruments(ctx context.Context, client *robinstock_go.Client, instrumentURLs []string) []string {
	var symbols []string

	// For now, fetch sequentially (can optimize with goroutines later if needed)
	for _, url := range instrumentURLs {
		symbol, err := stocks.GetSymbolByURL(ctx, client, url)
		if err == nil && symbol != "" {
			symbols = append(symbols, symbol)
		}
	}

	return symbols
}
