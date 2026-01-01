# Robinstock Go

Clean, production-ready Go library for Robinhood API interaction.

## Structure

```
robinstock_go/
├── account/               # Account Management
│   ├── accounts.go        # Positions, portfolio, holdings
│   └── accounts_test.go   # Tests
│
├── auth/                  # Authentication
│   ├── auth.go            # Login, logout, token management
│   ├── sheriff.go         # Sheriff verification
│   └── auth_test.go       # Tests
│
├── markets/               # Market Data
│   ├── markets.go         # Market info, hours, movers
│   └── markets_test.go    # Tests
│
├── profiles/              # User Profiles
│   ├── profiles.go        # User, investment, security profiles
│   └── profiles_test.go   # Tests
│
├── stocks/                # Stock Data
│   ├── stocks.go          # Quotes, fundamentals, historicals
│   └── stocks_test.go     # Tests
│
├── models/                # All Data Structures (DRY)
│   ├── models.go          # Core models
│   ├── models_additional.go  # Extended models
│   └── constants.go       # All constants
│
├── urls/                  # API Endpoints
│   └── urls.go            # All URL builders
│
├── utils/                 # Helper Functions
│   └── util.go            # Parsing, formatting utilities
│
├── client.go              # HTTP client
└── FUNCTIONS.md           # Implementation status

orders/  (TO BE IMPLEMENTED LAST - most complex)
```

## Design Principles

✅ **No Global State** - All state passed explicitly
✅ **No Disk Storage** - Caller manages tokens
✅ **DRY Enforced** - Structs defined once in `models/`
✅ **Container Safe** - Runs in multiple containers
✅ **Context Aware** - All calls accept `context.Context`
✅ **Thread Safe** - Concurrent use supported

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/yourusername/robinstock_go/auth"
    "github.com/yourusername/robinstock_go/stocks"
)

func main() {
    ctx := context.Background()
    client := auth.NewClient()

    // Login
    token, err := client.Login(ctx, "user", "pass", "")
    if err != nil {
        log.Fatal(err)
    }

    // Get quote
    quote, err := stocks.GetQuote(ctx, client, "AAPL")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("AAPL: $%s", quote.LastTradePrice)
}
```

## Implementation Status

See [FUNCTIONS.md](./FUNCTIONS.md) for detailed function-by-function status.

### Module Progress
- **Auth**: 100% complete ✅ (6/6 functions)
- **Account**: 100% complete ✅ (9/9 functions)
- **Profiles**: 100% complete ✅ (8/8 functions)
- **Markets**: 100% complete ✅ (7/7 functions)
- **Stocks**: 100% complete ✅ (13/13 functions)
- **Orders (Stock)**: 100% complete ✅ (18/18 functions)

## Development

```bash
# Build
go build ./...

# Test
go test ./...

# Run example
go run ./cmd/example/main.go
```

## Configuration

Environment variables:
```bash
export ROBINHOOD_USERNAME="your_username"
export ROBINHOOD_PASSWORD="your_password"
```

## Rules

1. **DRY**: Structs defined once in `models/`
2. **No global state**: Pass client explicitly
3. **No interfaces**: Unless needed for mocking
4. **No reflection**: Explicit types only
5. **Errors returned**: Never logged internally
6. **Context passed**: For cancellation/timeouts
7. **Function comments**: Required for all exported functions
8. **No inline comments**: Code should be self-documenting

## Reference

Python reference implementation available in `../python_rh_tester/` for API behavior testing.
