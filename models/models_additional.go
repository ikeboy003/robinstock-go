package models

type SPMover struct {
	Description                string  `json:"description"`
	InstrumentURL              string  `json:"instrument_url"`
	MarketHoursLastMovementPct float64 `json:"market_hours_last_movement_pct"`
	Symbol                     string  `json:"symbol"`
}

type Market struct {
	URL                     string `json:"url"`
	TOTP                    string `json:"totp"`
	Acronym                 string `json:"acronym"`
	Name                    string `json:"name"`
	City                    string `json:"city"`
	Country                 string `json:"country"`
	Timezone                string `json:"timezone"`
	Website                 string `json:"website"`
	OperatingMIC            string `json:"operating_mic"`
}

type MarketHours struct {
	IsOpen     bool   `json:"is_open"`
	OpensAt    string `json:"opens_at"`
	ClosesAt   string `json:"closes_at"`
	Date       string `json:"date"`
	ExtendedOpensAt  string `json:"extended_opens_at"`
	ExtendedClosesAt string `json:"extended_closes_at"`
}

type BasicProfile struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
	CreatedAt   string `json:"created_at"`
}

type InvestmentProfile struct {
	AnnualIncome           string `json:"annual_income"`
	InvestmentExperience   string `json:"investment_experience"`
	InvestmentObjective    string `json:"investment_objective"`
	LiquidityNeeds         string `json:"liquidity_needs"`
	LiquidNetWorth         string `json:"liquid_net_worth"`
	RiskTolerance          string `json:"risk_tolerance"`
	SourceOfFunds          string `json:"source_of_funds"`
	TaxBracket             string `json:"tax_bracket"`
	TimeHorizon            string `json:"time_horizon"`
	TotalNetWorth          string `json:"total_net_worth"`
}

type SecurityProfile struct {
	ObjectToDisclosure bool   `json:"object_to_disclosure"`
	SweepConsent       bool   `json:"sweep_consent"`
	ControlPerson      bool   `json:"control_person"`
	ControlPersonSecuritySymbol string `json:"control_person_security_symbol"`
	SecurityAffiliatedEmployee  bool   `json:"security_affiliated_employee"`
	SecurityAffiliatedFirmName  string `json:"security_affiliated_firm_name"`
}

type UserProfile struct {
	URL        string `json:"url"`
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	CreatedAt  string `json:"created_at"`
}

type Notification struct {
	ID         string `json:"id"`
	URL        string `json:"url"`
	Type       string `json:"type"`
	ARN        string `json:"arn"`
	Identifier string `json:"identifier"`
	Token      string `json:"token"`
}

type Holding struct {
	Symbol                string  `json:"symbol"`
	Price                 float64 `json:"price"`
	Quantity              float64 `json:"quantity"`
	AverageBuyPrice       float64 `json:"average_buy_price"`
	Equity                float64 `json:"equity"`
	PercentChange         float64 `json:"percent_change"`
	IntradayPercentChange float64 `json:"intraday_percent_change"`
	EquityChange          float64 `json:"equity_change"`
	Type                  string  `json:"type"`
	Name                  string  `json:"name"`
	ID                    string  `json:"id"`
	PERatio               float64 `json:"pe_ratio"`
	Percentage            float64 `json:"percentage"`
	PortfolioPercentage   float64 `json:"portfolio_percentage"`
}
