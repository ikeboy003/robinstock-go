package models

import "time"

type Auth struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	DeviceToken  string    `json:"device_token"`
	ExpiresIn    int       `json:"expires_in"`
	IssuedAt     time.Time `json:"issued_at"`
}

func (a *Auth) IsExpired() bool {
	if a.ExpiresIn == 0 {
		return false
	}
	expiresAt := a.IssuedAt.Add(time.Duration(a.ExpiresIn) * time.Second)
	return time.Now().After(expiresAt)
}

type OrderRequest struct {
	Symbol        string
	Quantity      float64
	Side          OrderSide
	Type          OrderType
	Price         *float64
	StopPrice     *float64
	TimeInForce   TimeInForce
	ExtendedHours bool
	AccountNumber *string
}

type Response struct {
	StatusCode int
	Data       map[string]interface{}
	Results    []map[string]interface{}
}

type Account struct {
	URL                     string `json:"url"`
	AccountNumber           string `json:"account_number"`
	Type                    string `json:"type"`
	CreatedAt               string `json:"created_at"`
	UpdatedAt               string `json:"updated_at"`
	Deactivated             bool   `json:"deactivated"`
	CashBalances            string `json:"cash_balances"`
	PortfolioURL            string `json:"portfolio"`
	BuyingPower             string `json:"buying_power"`
	MaxAchEarlyAccessAmount string `json:"max_ach_early_access_amount"`
	SweepEnabled            bool   `json:"sweep_enabled"`
	InstantEligibility      bool   `json:"instant_eligibility"`
	CashHeldForOrders       string `json:"cash_held_for_orders"`
	UnsettledFunds          string `json:"unsettled_funds"`
	UnsettledDebit          string `json:"unsettled_debit"`
}

// Position represents a stock position.
type Position struct {
	URL                            string `json:"url"`
	Instrument                     string `json:"instrument"`
	InstrumentID                   string `json:"instrument_id"`
	Account                        string `json:"account"`
	AccountNumber                  string `json:"account_number"`
	AverageBuyPrice                string `json:"average_buy_price"`
	PendingAverageBuyPrice         string `json:"pending_average_buy_price"`
	Quantity                       string `json:"quantity"`
	IntradayAverageBuyPrice        string `json:"intraday_average_buy_price"`
	IntradayQuantity               string `json:"intraday_quantity"`
	SharesHeldForBuys              string `json:"shares_held_for_buys"`
	SharesHeldForSells             string `json:"shares_held_for_sells"`
	SharesHeldForStockGrants       string `json:"shares_held_for_stock_grants"`
	SharesHeldForOptionsCollateral string `json:"shares_held_for_options_collateral"`
	SharesHeldForOptionsEvents     string `json:"shares_held_for_options_events"`
	SharesPendingFromOptionsEvents string `json:"shares_pending_from_options_events"`
	UpdatedAt                      string `json:"updated_at"`
	CreatedAt                      string `json:"created_at"`
}

// Instrument represents a stock instrument.
type Instrument struct {
	ID                 string `json:"id"`
	URL                string `json:"url"`
	Symbol             string `json:"symbol"`
	Name               string `json:"name"`
	SimpleName         string `json:"simple_name"`
	ListDate           string `json:"list_date"`
	Country            string `json:"country"`
	Type               string `json:"type"`
	Tradeable          bool   `json:"tradeable"`
	Fundamentals       string `json:"fundamentals"`
	Quote              string `json:"quote"`
	Market             string `json:"market"`
	State              string `json:"state"`
	DayTradeRatio      string `json:"day_trade_ratio"`
	MaintenanceRatio   string `json:"maintenance_ratio"`
	MarginInitialRatio string `json:"margin_initial_ratio"`
}

// Quote represents a stock quote.
type Quote struct {
	Symbol                   string `json:"symbol"`
	InstrumentID             string `json:"instrument_id"`
	AskPrice                 string `json:"ask_price"`
	AskSize                  int    `json:"ask_size"`
	BidPrice                 string `json:"bid_price"`
	BidSize                  int    `json:"bid_size"`
	LastTradePrice           string `json:"last_trade_price"`
	LastExtendedHoursTradePrice string `json:"last_extended_hours_trade_price"`
	PreviousClose            string `json:"previous_close"`
	AdjustedPreviousClose    string `json:"adjusted_previous_close"`
	TradingHalted            bool   `json:"trading_halted"`
	HasTraded                bool   `json:"has_traded"`
	UpdatedAt                string `json:"updated_at"`
}

// Order represents a stock order.
type Order struct {
	ID                     string      `json:"id"`
	URL                    string      `json:"url"`
	Account                string      `json:"account"`
	Instrument             string      `json:"instrument"`
	Symbol                 string      `json:"symbol"`
	Type                   string      `json:"type"`
	Side                   string      `json:"side"`
	TimeInForce            string      `json:"time_in_force"`
	Trigger                string      `json:"trigger"`
	Price                  string      `json:"price"`
	StopPrice              string      `json:"stop_price"`
	Quantity               string      `json:"quantity"`
	AveragePrice           string      `json:"average_price"`
	CumulativeQuantity     string      `json:"cumulative_quantity"`
	State                  string      `json:"state"`
	CreatedAt              string      `json:"created_at"`
	UpdatedAt              string      `json:"updated_at"`
	LastTransactionAt      string      `json:"last_transaction_at"`
	Executions             []Execution `json:"executions"`
	ExtendedHours          bool        `json:"extended_hours"`
	OverrideDayTradeChecks bool        `json:"override_day_trade_checks"`
	OverrideDtbpChecks     bool        `json:"override_dtbp_checks"`
	RefID                  string      `json:"ref_id"`
}

// Execution represents an order execution.
type Execution struct {
	ID               string `json:"id"`
	Price            string `json:"price"`
	Quantity         string `json:"quantity"`
	SettlementDate   string `json:"settlement_date"`
	Timestamp        string `json:"timestamp"`
}

// Portfolio represents portfolio data.
type Portfolio struct {
	URL                             string `json:"url"`
	Account                         string `json:"account"`
	StartDate                       string `json:"start_date"`
	MarketValue                     string `json:"market_value"`
	Equity                          string `json:"equity"`
	ExtendedHoursMarketValue        string `json:"extended_hours_market_value"`
	ExtendedHoursEquity             string `json:"extended_hours_equity"`
	ExtendedHoursPreviousClose      string `json:"extended_hours_previous_close"`
	LastCoreMarketValue             string `json:"last_core_market_value"`
	LastCoreEquity                  string `json:"last_core_equity"`
	ExcessMargin                    string `json:"excess_margin"`
	ExcessMaintenanceWithUnclearedDeposits string `json:"excess_maintenance_with_uncleared_deposits"`
	EquityPreviousClose             string `json:"equity_previous_close"`
	AdjustedEquityPreviousClose     string `json:"adjusted_equity_previous_close"`
	Withdrawable                    string `json:"withdrawable"`
	UnsettledFunds                  string `json:"unsettled_funds"`
}

// Dividend represents a dividend payment.
type Dividend struct {
	ID               string    `json:"id"`
	URL              string    `json:"url"`
	Account          string    `json:"account"`
	Instrument       string    `json:"instrument"`
	Amount           string    `json:"amount"`
	Rate             string    `json:"rate"`
	Position         string    `json:"position"`
	WithholdingAmount string   `json:"withholding"`
	RecordDate       string    `json:"record_date"`
	PayableDate      string    `json:"payable_date"`
	PaidAt           string    `json:"paid_at"`
	State            string    `json:"state"`
}

// Fundamental represents stock fundamentals.
type Fundamental struct {
	Open              string `json:"open"`
	High              string `json:"high"`
	Low               string `json:"low"`
	Volume            string `json:"volume"`
	AverageVolume     string `json:"average_volume"`
	AverageVolume2Weeks string `json:"average_volume_2_weeks"`
	High52Weeks       string `json:"high_52_weeks"`
	Low52Weeks        string `json:"low_52_weeks"`
	MarketCap         string `json:"market_cap"`
	DividendYield     string `json:"dividend_yield"`
	PERatio           string `json:"pe_ratio"`
	Description       string `json:"description"`
	Instrument        string `json:"instrument"`
}

// HistoricalData represents historical price data.
type HistoricalData struct {
	BeginsAt         time.Time `json:"begins_at"`
	OpenPrice        string    `json:"open_price"`
	ClosePrice       string    `json:"close_price"`
	HighPrice        string    `json:"high_price"`
	LowPrice         string    `json:"low_price"`
	Volume           int       `json:"volume"`
	Session          string    `json:"session"`
	Interpolated     bool      `json:"interpolated"`
}

