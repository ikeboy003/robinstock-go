package screener

import (
	"context"
	"fmt"

	robinstock "github.com/ikeboy003/robinstock-go"
)

// ScreenerRequest represents a request to the Robinhood screener API.
type ScreenerRequest struct {
	Columns       []string              `json:"columns"`
	Indicators    []Indicator           `json:"indicators"`
	IsPollable    bool                  `json:"is_pollable"`
	SortBy        string                `json:"sort_by"`
	SortDirection string                `json:"sort_direction"`
}

// Indicator represents a filter indicator in the screener.
type Indicator struct {
	Key    string `json:"key"`
	Filter Filter `json:"filter"`
}

// Filter represents a filter for an indicator.
type Filter struct {
	Type       string      `json:"type"`
	Selection  *Selection  `json:"selection,omitempty"`
	Selections []Selection `json:"selections,omitempty"`
}

// Selection represents a filter selection.
type Selection struct {
	ID              string           `json:"id"`
	SecondaryFilter *SecondaryFilter `json:"secondary_filter,omitempty"`
}

// SecondaryFilter represents a secondary filter with min/max values.
type SecondaryFilter struct {
	Type string   `json:"type"`
	Min  *float64 `json:"min,omitempty"`
	Max  *float64 `json:"max,omitempty"`
}

// ScreenerResult represents a result from the screener API.
type ScreenerResult struct {
	InstrumentID string                 `json:"instrument_id"`
	Symbol       string                 `json:"symbol"`
	Data         map[string]interface{} `json:"data"`
}

// ScreenerResponse represents the response from the screener API.
type ScreenerResponse struct {
	Results []ScreenerResult `json:"results"`
	Next    *string          `json:"next"`
}

// Scan performs a screener scan with the given request.
func Scan(ctx context.Context, client *robinstock.Client, request ScreenerRequest) ([]ScreenerResult, error) {
	if !client.IsAuthenticated() {
		return nil, robinstock.ErrNotAuthenticated
	}

	url := "https://bonfire.robinhood.com/screeners/scan/"
	resp, err := client.Post(ctx, url, request, true)
	if err != nil {
		return nil, fmt.Errorf("screener scan failed: %w", err)
	}

	results := parseScreenerResponse(resp.Data)
	return results, nil
}

func parseScreenerResponse(data map[string]interface{}) []ScreenerResult {
	var results []ScreenerResult

	rowsData, ok := data["rows"].([]interface{})
	if !ok || len(rowsData) == 0 {
		return results
	}

	for i, item := range rowsData {
		if i == 0 {
			continue
		}

		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		instrumentID := getString(itemMap, "instrument_id")
		symbol := getString(itemMap, "instrument_symbol")

		parsedData := make(map[string]interface{})
		parsedData["instrument_id"] = instrumentID
		parsedData["symbol"] = symbol

		if items, ok := itemMap["items"].([]interface{}); ok {
			for _, item := range items {
				if itemComponentMap, ok := item.(map[string]interface{}); ok {
					if component, ok := itemComponentMap["component"].(map[string]interface{}); ok {
						componentType, _ := component["sdui_component_type"].(string)

						switch componentType {
						case "TABLE_INSTRUMENT_NAME":
							if name, ok := component["name"].(string); ok {
								parsedData["name"] = name
							}
						case "TABLE_1D_CHANGE_ITEM":
							if defaultValue, ok := component["default_value"].(map[string]interface{}); ok {
								if value, ok := defaultValue["value"].(string); ok {
									parsedData["1d_price_change"] = value
								}
								if direction, ok := defaultValue["direction"].(string); ok {
									parsedData["1d_direction"] = direction
								}
							}
						case "TABLE_SHARE_PRICE_ITEM":
							if defaultValue, ok := component["default_value"].(map[string]interface{}); ok {
								if price, ok := defaultValue["price"].(string); ok {
									parsedData["price"] = price
								}
							}
						case "TEXT":
							if value, ok := component["text"].(string); ok {
								if title, ok := component["title"].(string); ok {
									parsedData[title] = value
								} else {
									if len(parsedData) == 5 {
										parsedData["volume"] = value
									} else if len(parsedData) == 6 {
										parsedData["market_cap"] = value
									} else if len(parsedData) == 7 {
										parsedData["avg_volume"] = value
									}
								}
							}
						}
					}
				}
			}
		}

		result := ScreenerResult{
			InstrumentID: instrumentID,
			Symbol:       symbol,
			Data:         parsedData,
		}

		results = append(results, result)
	}

	return results
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// NewLosersScreener creates a screener request for stocks that fell by at least minPercentChange.
func NewLosersScreener(minPercentChange float64) ScreenerRequest {
	return ScreenerRequest{
		Columns: []string{
			"sparkline",
			"1d_price_change",
			"price",
			"todays_volume",
			"market_cap",
			"average_volume_30_days",
		},
		Indicators: []Indicator{
			{
				Key: "1d_price_change",
				Filter: Filter{
					Type: "SINGLE_SELECT",
					Selection: &Selection{
						ID: "custom",
						SecondaryFilter: &SecondaryFilter{
							Type: "PERCENT_RANGE",
							Max:  &minPercentChange,
						},
					},
				},
			},
		},
		IsPollable:    true,
		SortBy:        "1d_price_change",
		SortDirection: "ASC",
	}
}

// NewGainersScreener creates a screener request for stocks that gained by at least minPercentChange.
func NewGainersScreener(minPercentChange float64) ScreenerRequest {
	return ScreenerRequest{
		Columns: []string{
			"sparkline",
			"1d_price_change",
			"price",
			"todays_volume",
			"market_cap",
			"average_volume_30_days",
		},
		Indicators: []Indicator{
			{
				Key: "1d_price_change",
				Filter: Filter{
					Type: "SINGLE_SELECT",
					Selection: &Selection{
						ID: "custom",
						SecondaryFilter: &SecondaryFilter{
							Type: "PERCENT_RANGE",
							Min:  &minPercentChange,
						},
					},
				},
			},
		},
		IsPollable:    true,
		SortBy:        "1d_price_change",
		SortDirection: "DESC",
	}
}

// NewVolumeScreener creates a screener request for stocks with high volume.
func NewVolumeScreener(minVolumeChange float64) ScreenerRequest {
	return ScreenerRequest{
		Columns: []string{
			"sparkline",
			"1d_price_change",
			"price",
			"todays_volume",
			"market_cap",
			"average_volume_30_days",
			"todays_volume_ratio",
		},
		Indicators: []Indicator{
			{
				Key: "todays_volume_ratio",
				Filter: Filter{
					Type: "SINGLE_SELECT",
					Selection: &Selection{
						ID: "custom",
						SecondaryFilter: &SecondaryFilter{
							Type: "PERCENT_RANGE",
							Min:  &minVolumeChange,
						},
					},
				},
			},
		},
		IsPollable:    true,
		SortBy:        "todays_volume_ratio",
		SortDirection: "DESC",
	}
}

// Market cap constants for convenience.
const (
	MegaCap   = "mkt_cap_mega_cap"
	LargeCap  = "mkt_cap_large_cap"
	MidCap    = "mkt_cap_mid_cap"
	SmallCap  = "mkt_cap_small_cap"
	MicroCap  = "mkt_cap_micro_cap"
	AllCaps   = "mkt_cap_all"
)

