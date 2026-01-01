package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func GetString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func GetInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

func GetFloat(data map[string]interface{}, key string) float64 {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return 0.0
}

func GetBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func NormalizeSymbol(symbol string) string {
	return strings.ToUpper(strings.TrimSpace(symbol))
}

func NormalizeSymbols(symbols []string) []string {
	normalized := make([]string, len(symbols))
	for i, sym := range symbols {
		normalized[i] = NormalizeSymbol(sym)
	}
	return normalized
}

func JoinSymbols(symbols []string) string {
	return strings.Join(symbols, ",")
}

// ParseFloat converts a string or interface to float64.
func ParseFloat(v interface{}) float64 {
	switch val := v.(type) {
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	}
	return 0.0
}

func BuildURL(base string, params map[string]string) string {
	if len(params) == 0 {
		return base
	}

	var parts []string
	for key, val := range params {
		if val != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", key, val))
		}
	}

	if len(parts) == 0 {
		return base
	}

	return base + "?" + strings.Join(parts, "&")
}

// RoundPrice rounds a price to 2 decimal places.
func RoundPrice(price float64) float64 {
	return math.Round(price*100) / 100
}

// Address returns a pointer to a string.
func Address(s string) *string {
	return &s
}
