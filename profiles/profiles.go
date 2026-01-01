package profiles

import (
	"context"
	"log"

	"github.com/yourusername/robinstock_go"
	"github.com/yourusername/robinstock_go/models"
	"github.com/yourusername/robinstock_go/urls"
	"github.com/yourusername/robinstock_go/utils"
)

// GetAccountProfile returns account profile information.
func GetAccountProfile(ctx context.Context, client *robinstock_go.Client, accountNumber string) (*models.Account, error) {
	log.Printf("GetAccountProfile: Fetching profile for account %s...\n", accountNumber)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := models.BaseURL + "/accounts/" + accountNumber + "/"
	resp, err := client.Get(ctx, url, nil, true)
	if err != nil {
		log.Printf("GetAccountProfile: Error: %v\n", err)
		return nil, err
	}

	account := &models.Account{
		URL:                     utils.GetString(resp.Data, "url"),
		AccountNumber:           utils.GetString(resp.Data, "account_number"),
		Type:                    utils.GetString(resp.Data, "type"),
		CreatedAt:               utils.GetString(resp.Data, "created_at"),
		UpdatedAt:               utils.GetString(resp.Data, "updated_at"),
		Deactivated:             utils.GetBool(resp.Data, "deactivated"),
		CashBalances:            utils.GetString(resp.Data, "cash_balances"),
		PortfolioURL:            utils.GetString(resp.Data, "portfolio"),
		BuyingPower:             utils.GetString(resp.Data, "buying_power"),
		MaxAchEarlyAccessAmount: utils.GetString(resp.Data, "max_ach_early_access_amount"),
		SweepEnabled:            utils.GetBool(resp.Data, "sweep_enabled"),
		InstantEligibility:      utils.GetBool(resp.Data, "instant_eligibility"),
		CashHeldForOrders:       utils.GetString(resp.Data, "cash_held_for_orders"),
		UnsettledFunds:          utils.GetString(resp.Data, "unsettled_funds"),
		UnsettledDebit:          utils.GetString(resp.Data, "unsettled_debit"),
	}

	log.Printf("GetAccountProfile: Account %s retrieved\n", account.AccountNumber)
	return account, nil
}

// GetAllAccountProfiles returns all account profiles for the user.
func GetAllAccountProfiles(ctx context.Context, client *robinstock_go.Client) ([]models.Account, error) {
	log.Println("GetAllAccountProfiles: Fetching all account profiles...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := models.BaseURL + "/accounts/?default_to_all_accounts=true"
	resp, err := client.Get(ctx, url, nil, true)
	if err != nil {
		log.Printf("GetAllAccountProfiles: Error: %v\n", err)
		return nil, err
	}

	var accounts []models.Account
	if results, ok := resp.Data["results"].([]interface{}); ok {
		for _, item := range results {
			if accountData, ok := item.(map[string]interface{}); ok {
				account := models.Account{
					URL:                     utils.GetString(accountData, "url"),
					AccountNumber:           utils.GetString(accountData, "account_number"),
					Type:                    utils.GetString(accountData, "type"),
					CreatedAt:               utils.GetString(accountData, "created_at"),
					UpdatedAt:               utils.GetString(accountData, "updated_at"),
					Deactivated:             utils.GetBool(accountData, "deactivated"),
					CashBalances:            utils.GetString(accountData, "cash_balances"),
					PortfolioURL:            utils.GetString(accountData, "portfolio"),
					BuyingPower:             utils.GetString(accountData, "buying_power"),
					MaxAchEarlyAccessAmount: utils.GetString(accountData, "max_ach_early_access_amount"),
					SweepEnabled:            utils.GetBool(accountData, "sweep_enabled"),
					InstantEligibility:      utils.GetBool(accountData, "instant_eligibility"),
					CashHeldForOrders:       utils.GetString(accountData, "cash_held_for_orders"),
					UnsettledFunds:          utils.GetString(accountData, "unsettled_funds"),
					UnsettledDebit:          utils.GetString(accountData, "unsettled_debit"),
				}
				accounts = append(accounts, account)
			}
		}
	}

	log.Printf("GetAllAccountProfiles: Retrieved %d accounts\n", len(accounts))
	return accounts, nil
}

// GetBasicProfile returns the user's basic profile information.
func GetBasicProfile(ctx context.Context, client *robinstock_go.Client) (*models.BasicProfile, error) {
	log.Println("GetBasicProfile: Fetching basic profile...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	resp, err := client.Get(ctx, urls.BasicProfileURL(), nil, true)
	if err != nil {
		log.Printf("GetBasicProfile: Error: %v\n", err)
		return nil, err
	}

	profile := &models.BasicProfile{
		FirstName:   utils.GetString(resp.Data, "first_name"),
		LastName:    utils.GetString(resp.Data, "last_name"),
		Email:       utils.GetString(resp.Data, "email"),
		PhoneNumber: utils.GetString(resp.Data, "phone_number"),
		Username:    utils.GetString(resp.Data, "username"),
		CreatedAt:   utils.GetString(resp.Data, "created_at"),
	}

	log.Printf("GetBasicProfile: Retrieved profile for %s %s\n", profile.FirstName, profile.LastName)
	return profile, nil
}

// GetInvestmentProfile returns the user's investment profile.
func GetInvestmentProfile(ctx context.Context, client *robinstock_go.Client) (*models.InvestmentProfile, error) {
	log.Println("GetInvestmentProfile: Fetching investment profile...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	resp, err := client.Get(ctx, urls.InvestmentProfileURL(), nil, true)
	if err != nil {
		log.Printf("GetInvestmentProfile: Error: %v\n", err)
		return nil, err
	}

	profile := &models.InvestmentProfile{
		AnnualIncome:         utils.GetString(resp.Data, "annual_income"),
		InvestmentExperience: utils.GetString(resp.Data, "investment_experience"),
		InvestmentObjective:  utils.GetString(resp.Data, "investment_objective"),
		LiquidityNeeds:       utils.GetString(resp.Data, "liquidity_needs"),
		LiquidNetWorth:       utils.GetString(resp.Data, "liquid_net_worth"),
		RiskTolerance:        utils.GetString(resp.Data, "risk_tolerance"),
		SourceOfFunds:        utils.GetString(resp.Data, "source_of_funds"),
		TaxBracket:           utils.GetString(resp.Data, "tax_bracket"),
		TimeHorizon:          utils.GetString(resp.Data, "time_horizon"),
		TotalNetWorth:        utils.GetString(resp.Data, "total_net_worth"),
	}

	log.Printf("GetInvestmentProfile: Investment objective: %s\n", profile.InvestmentObjective)
	return profile, nil
}

// GetPortfolioProfile returns portfolio information for a specific account.
func GetPortfolioProfile(ctx context.Context, client *robinstock_go.Client, accountNumber string) (*models.Portfolio, error) {
	log.Printf("GetPortfolioProfile: Fetching portfolio for account %s...\n", accountNumber)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := models.BaseURL + "/portfolios/" + accountNumber + "/"
	resp, err := client.Get(ctx, url, nil, true)
	if err != nil {
		log.Printf("GetPortfolioProfile: Error: %v\n", err)
		return nil, err
	}

	portfolio := &models.Portfolio{
		URL:                                    utils.GetString(resp.Data, "url"),
		Account:                                utils.GetString(resp.Data, "account"),
		StartDate:                              utils.GetString(resp.Data, "start_date"),
		MarketValue:                            utils.GetString(resp.Data, "market_value"),
		Equity:                                 utils.GetString(resp.Data, "equity"),
		EquityPreviousClose:                    utils.GetString(resp.Data, "equity_previous_close"),
		ExtendedHoursMarketValue:               utils.GetString(resp.Data, "extended_hours_market_value"),
		ExtendedHoursEquity:                    utils.GetString(resp.Data, "extended_hours_equity"),
		ExtendedHoursPreviousClose:             utils.GetString(resp.Data, "extended_hours_previous_close"),
		LastCoreMarketValue:                    utils.GetString(resp.Data, "last_core_market_value"),
		LastCoreEquity:                         utils.GetString(resp.Data, "last_core_equity"),
		ExcessMargin:                           utils.GetString(resp.Data, "excess_margin"),
		ExcessMaintenanceWithUnclearedDeposits: utils.GetString(resp.Data, "excess_maintenance_with_uncleared_deposits"),
		AdjustedEquityPreviousClose:            utils.GetString(resp.Data, "adjusted_equity_previous_close"),
		Withdrawable:                           utils.GetString(resp.Data, "withdrawable"),
		UnsettledFunds:                         utils.GetString(resp.Data, "unsettled_funds"),
	}

	log.Printf("GetPortfolioProfile: Market value: %s\n", portfolio.MarketValue)
	return portfolio, nil
}

// GetAllPortfolioProfiles returns portfolio information for all accounts.
func GetAllPortfolioProfiles(ctx context.Context, client *robinstock_go.Client) ([]models.Portfolio, error) {
	log.Println("GetAllPortfolioProfiles: Fetching all portfolio profiles...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := models.BaseURL + "/portfolios/"
	results, err := client.FetchAllPages(ctx, url, true)
	if err != nil {
		log.Printf("GetAllPortfolioProfiles: Error: %v\n", err)
		return nil, err
	}

	var portfolios []models.Portfolio
	for _, result := range results {
		portfolio := models.Portfolio{
			URL:                                    utils.GetString(result, "url"),
			Account:                                utils.GetString(result, "account"),
			StartDate:                              utils.GetString(result, "start_date"),
			MarketValue:                            utils.GetString(result, "market_value"),
			Equity:                                 utils.GetString(result, "equity"),
			EquityPreviousClose:                    utils.GetString(result, "equity_previous_close"),
			ExtendedHoursMarketValue:               utils.GetString(result, "extended_hours_market_value"),
			ExtendedHoursEquity:                    utils.GetString(result, "extended_hours_equity"),
			ExtendedHoursPreviousClose:             utils.GetString(result, "extended_hours_previous_close"),
			LastCoreMarketValue:                    utils.GetString(result, "last_core_market_value"),
			LastCoreEquity:                         utils.GetString(result, "last_core_equity"),
			ExcessMargin:                           utils.GetString(result, "excess_margin"),
			ExcessMaintenanceWithUnclearedDeposits: utils.GetString(result, "excess_maintenance_with_uncleared_deposits"),
			AdjustedEquityPreviousClose:            utils.GetString(result, "adjusted_equity_previous_close"),
			Withdrawable:                           utils.GetString(result, "withdrawable"),
			UnsettledFunds:                         utils.GetString(result, "unsettled_funds"),
		}
		portfolios = append(portfolios, portfolio)
	}

	log.Printf("GetAllPortfolioProfiles: Retrieved %d portfolios\n", len(portfolios))
	return portfolios, nil
}

// GetSecurityProfile returns the user's security profile.
func GetSecurityProfile(ctx context.Context, client *robinstock_go.Client) (*models.SecurityProfile, error) {
	log.Println("GetSecurityProfile: Fetching security profile...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	resp, err := client.Get(ctx, urls.SecurityProfileURL(), nil, true)
	if err != nil {
		log.Printf("GetSecurityProfile: Error: %v\n", err)
		return nil, err
	}

	profile := &models.SecurityProfile{
		ObjectToDisclosure:          utils.GetBool(resp.Data, "object_to_disclosure"),
		SweepConsent:                utils.GetBool(resp.Data, "sweep_consent"),
		ControlPerson:               utils.GetBool(resp.Data, "control_person"),
		ControlPersonSecuritySymbol: utils.GetString(resp.Data, "control_person_security_symbol"),
		SecurityAffiliatedEmployee:  utils.GetBool(resp.Data, "security_affiliated_employee"),
		SecurityAffiliatedFirmName:  utils.GetString(resp.Data, "security_affiliated_firm_name"),
	}

	log.Println("GetSecurityProfile: Security profile retrieved")
	return profile, nil
}

// GetUserProfile returns the full user profile.
func GetUserProfile(ctx context.Context, client *robinstock_go.Client) (*models.UserProfile, error) {
	log.Println("GetUserProfile: Fetching user profile...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	resp, err := client.Get(ctx, urls.UserProfileURL(), nil, true)
	if err != nil {
		log.Printf("GetUserProfile: Error: %v\n", err)
		return nil, err
	}

	profile := &models.UserProfile{
		URL:       utils.GetString(resp.Data, "url"),
		ID:        utils.GetString(resp.Data, "id"),
		Username:  utils.GetString(resp.Data, "username"),
		Email:     utils.GetString(resp.Data, "email"),
		FirstName: utils.GetString(resp.Data, "first_name"),
		LastName:  utils.GetString(resp.Data, "last_name"),
		CreatedAt: utils.GetString(resp.Data, "created_at"),
	}

	log.Printf("GetUserProfile: Retrieved profile for user %s\n", profile.Username)
	return profile, nil
}

// GetPortfolioHistoricals returns historical portfolio performance for an account.
func GetPortfolioHistoricals(ctx context.Context, client *robinstock_go.Client, accountNumber, interval, span string) ([]models.HistoricalData, error) {
	log.Printf("GetPortfolioHistoricals: Fetching historicals for account %s...\n", accountNumber)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	params := map[string]string{
		"interval": interval,
		"span":     span,
	}

	url := models.BaseURL + "/portfolios/historicals/" + accountNumber + "/"
	resp, err := client.Get(ctx, url, params, true)
	if err != nil {
		log.Printf("GetPortfolioHistoricals: Error: %v\n", err)
		return nil, err
	}

	var historicals []models.HistoricalData
	if equityHistoricals, ok := resp.Data["equity_historicals"].([]interface{}); ok {
		for _, item := range equityHistoricals {
			if hist, ok := item.(map[string]interface{}); ok {
				historical := models.HistoricalData{
					OpenPrice:    utils.GetString(hist, "open_equity"),
					ClosePrice:   utils.GetString(hist, "close_equity"),
					HighPrice:    utils.GetString(hist, "high_equity"),
					LowPrice:     utils.GetString(hist, "low_equity"),
					Session:      utils.GetString(hist, "session"),
					Interpolated: utils.GetBool(hist, "interpolated"),
				}
				historicals = append(historicals, historical)
			}
		}
	}

	log.Printf("GetPortfolioHistoricals: Retrieved %d historical data points\n", len(historicals))
	return historicals, nil
}
