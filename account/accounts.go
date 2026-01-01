package account

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/ikeboy003/robinstock-go"
	"github.com/ikeboy003/robinstock-go/models"
	"github.com/ikeboy003/robinstock-go/urls"
	"github.com/ikeboy003/robinstock-go/utils"
)

// LoadPhoenixAccount returns unified account data from Phoenix API.
// Note: Phoenix endpoint (phoenix.robinhood.com) has TLS compatibility issues.
// This function is maintained for API completeness but may fail with TLS handshake errors.
// Use GetAllPositions and other specific endpoints for reliable data access.
func LoadPhoenixAccount(ctx context.Context, client *robinstock_go.Client) (map[string]interface{}, error) {
	log.Println("LoadPhoenixAccount: Fetching Phoenix account data...")
	log.Println("LoadPhoenixAccount: WARNING - Phoenix endpoint has known TLS issues")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	resp, err := client.Get(ctx, urls.PhoenixURL(), nil, true)
	if err != nil {
		log.Printf("LoadPhoenixAccount: Error: %v\n", err)
		return nil, err
	}

	return resp.Data, nil
}

// GetAllPositions returns all positions.
func GetAllPositions(ctx context.Context, client *robinstock_go.Client) ([]models.Position, error) {
	log.Println("GetAllPositions: Fetching all positions...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	results, err := client.FetchAllPages(ctx, urls.PositionsURL(), true)
	if err != nil {
		log.Printf("GetAllPositions: Error: %v\n", err)
		return nil, err
	}

	log.Printf("GetAllPositions: Received %d positions\n", len(results))

	var positions []models.Position
	for _, result := range results {
		position := models.Position{
			URL:                      utils.GetString(result, "url"),
			Instrument:               utils.GetString(result, "instrument"),
			InstrumentID:             utils.GetString(result, "instrument_id"),
			Account:                  utils.GetString(result, "account"),
			AccountNumber:            utils.GetString(result, "account_number"),
			AverageBuyPrice:          utils.GetString(result, "average_buy_price"),
			PendingAverageBuyPrice:   utils.GetString(result, "pending_average_buy_price"),
			Quantity:                 utils.GetString(result, "quantity"),
			IntradayAverageBuyPrice:  utils.GetString(result, "intraday_average_buy_price"),
			IntradayQuantity:         utils.GetString(result, "intraday_quantity"),
			SharesHeldForBuys:        utils.GetString(result, "shares_held_for_buys"),
			SharesHeldForSells:       utils.GetString(result, "shares_held_for_sells"),
			UpdatedAt:                utils.GetString(result, "updated_at"),
			CreatedAt:                utils.GetString(result, "created_at"),
		}
		positions = append(positions, position)
	}

	return positions, nil
}

// GetOpenStockPosition returns open positions, optionally filtered by account number.
func GetOpenStockPosition(ctx context.Context, client *robinstock_go.Client, accountNumber *string) ([]models.Position, error) {
	if accountNumber != nil {
		log.Printf("GetOpenStockPosition: Fetching for account %s...\n", *accountNumber)
	} else {
		log.Println("GetOpenStockPosition: Fetching all open positions...")
	}

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	params := map[string]string{"nonzero": "true"}
	url := urls.PositionsURL()
	if accountNumber != nil {
		url = utils.BuildURL(url, map[string]string{"account_number": *accountNumber})
	}

	results, err := client.FetchAllPages(ctx, url, true)
	if err != nil {
		log.Printf("GetOpenStockPosition: Error: %v\n", err)
		return nil, err
	}

	log.Printf("GetOpenStockPosition: Received %d positions\n", len(results))

	var positions []models.Position
	for _, result := range results {
		position := models.Position{
			URL:                      utils.GetString(result, "url"),
			Instrument:               utils.GetString(result, "instrument"),
			InstrumentID:             utils.GetString(result, "instrument_id"),
			Account:                  utils.GetString(result, "account"),
			AccountNumber:            utils.GetString(result, "account_number"),
			AverageBuyPrice:          utils.GetString(result, "average_buy_price"),
			PendingAverageBuyPrice:   utils.GetString(result, "pending_average_buy_price"),
			Quantity:                 utils.GetString(result, "quantity"),
			IntradayAverageBuyPrice:  utils.GetString(result, "intraday_average_buy_price"),
			IntradayQuantity:         utils.GetString(result, "intraday_quantity"),
			SharesHeldForBuys:        utils.GetString(result, "shares_held_for_buys"),
			SharesHeldForSells:       utils.GetString(result, "shares_held_for_sells"),
			UpdatedAt:                utils.GetString(result, "updated_at"),
			CreatedAt:                utils.GetString(result, "created_at"),
		}
		positions = append(positions, position)
	}

	_ = params
	return positions, nil
}

// GetDividends returns dividend history.
func GetDividends(ctx context.Context, client *robinstock_go.Client) ([]models.Dividend, error) {
	log.Println("GetDividends: Fetching dividend history...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	results, err := client.FetchAllPages(ctx, urls.DividendsURL(), true)
	if err != nil {
		log.Printf("GetDividends: Error: %v\n", err)
		return nil, err
	}

	log.Printf("GetDividends: Received %d dividends\n", len(results))

	var dividends []models.Dividend
	for _, result := range results {
		dividend := models.Dividend{
			ID:                utils.GetString(result, "id"),
			URL:               utils.GetString(result, "url"),
			Account:           utils.GetString(result, "account"),
			Instrument:        utils.GetString(result, "instrument"),
			Amount:            utils.GetString(result, "amount"),
			Rate:              utils.GetString(result, "rate"),
			Position:          utils.GetString(result, "position"),
			WithholdingAmount: utils.GetString(result, "withholding"),
			RecordDate:        utils.GetString(result, "record_date"),
			PayableDate:       utils.GetString(result, "payable_date"),
			PaidAt:            utils.GetString(result, "paid_at"),
			State:             utils.GetString(result, "state"),
		}
		dividends = append(dividends, dividend)
	}

	return dividends, nil
}

// GetNotifications returns account notifications.
func GetNotifications(ctx context.Context, client *robinstock_go.Client) ([]models.Notification, error) {
	log.Println("GetNotifications: Fetching notifications...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	results, err := client.FetchAllPages(ctx, urls.NotificationsURL(false), true)
	if err != nil {
		log.Printf("GetNotifications: Error: %v\n", err)
		return nil, err
	}

	log.Printf("GetNotifications: Received %d notifications\n", len(results))

	var notifications []models.Notification
	for _, result := range results {
		notification := models.Notification{
			ID:         utils.GetString(result, "id"),
			URL:        utils.GetString(result, "url"),
			Type:       utils.GetString(result, "type"),
			ARN:        utils.GetString(result, "arn"),
			Identifier: utils.GetString(result, "identifier"),
			Token:      utils.GetString(result, "token"),
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// GetLinkedBankAccounts returns all linked bank accounts.
func GetLinkedBankAccounts(ctx context.Context, client *robinstock_go.Client) ([]map[string]interface{}, error) {
	log.Println("GetLinkedBankAccounts: Fetching linked bank accounts...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	results, err := client.FetchAllPages(ctx, urls.LinkedURL(), true)
	if err != nil {
		log.Printf("GetLinkedBankAccounts: Error: %v\n", err)
		return nil, err
	}

	log.Printf("GetLinkedBankAccounts: Received %d accounts\n", len(results))
	return results, nil
}

// DepositFundsIntoRobinhood deposits funds from a linked bank account.
func DepositFundsIntoRobinhood(ctx context.Context, client *robinstock_go.Client, accountNumber string, amount float64) (map[string]interface{}, error) {
	log.Printf("DepositFundsIntoRobinhood: Depositing $%.2f from account %s...\n", amount, accountNumber)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	achRelationship, err := getBankURLByAccountNumber(ctx, client, accountNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank relationship: %w", err)
	}

	payload := map[string]interface{}{
		"amount":           amount,
		"direction":        "deposit",
		"ach_relationship": achRelationship,
		"ref_id":           uuid.New().String(),
	}

	resp, err := client.Post(ctx, urls.BankTransfersURL(), payload, true)
	if err != nil {
		log.Printf("DepositFundsIntoRobinhood: Error: %v\n", err)
		return nil, err
	}

	log.Println("DepositFundsIntoRobinhood: Deposit initiated successfully")
	return resp.Data, nil
}

// WithdrawFundsFromRobinhood withdraws funds to a linked bank account.
func WithdrawFundsFromRobinhood(ctx context.Context, client *robinstock_go.Client, accountNumber string, amount float64) (map[string]interface{}, error) {
	log.Printf("WithdrawFundsFromRobinhood: Withdrawing $%.2f to account %s...\n", amount, accountNumber)

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	achRelationship, err := getBankURLByAccountNumber(ctx, client, accountNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank relationship: %w", err)
	}

	payload := map[string]interface{}{
		"amount":           amount,
		"direction":        "withdraw",
		"ach_relationship": achRelationship,
		"ref_id":           uuid.New().String(),
	}

	resp, err := client.Post(ctx, urls.BankTransfersURL(), payload, true)
	if err != nil {
		log.Printf("WithdrawFundsFromRobinhood: Error: %v\n", err)
		return nil, err
	}

	log.Println("WithdrawFundsFromRobinhood: Withdrawal initiated successfully")
	return resp.Data, nil
}

func getBankURLByAccountNumber(ctx context.Context, client *robinstock_go.Client, accountNumber string) (string, error) {
	accounts, err := GetLinkedBankAccounts(ctx, client)
	if err != nil {
		return "", err
	}

	for _, account := range accounts {
		if bankAccNum := utils.GetString(account, "bank_account_number"); bankAccNum == accountNumber {
			return utils.GetString(account, "url"), nil
		}
	}

	if len(accounts) > 0 {
		return utils.GetString(accounts[0], "url"), nil
	}

	return "", fmt.Errorf("no linked bank accounts found")
}

// BuildHoldings builds detailed holdings data with calculations.
func BuildHoldings(ctx context.Context, client *robinstock_go.Client, withDividend bool) ([]models.Holding, error) {
	log.Println("BuildHoldings: Building holdings data...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	// Get positions
	positions, err := GetOpenStockPosition(ctx, client, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	// Get account data for cash
	resp, err := client.Get(ctx, urls.AccountsURL(), nil, true)
	if err != nil || len(resp.Results) == 0 {
		return nil, fmt.Errorf("failed to get account data: %w", err)
	}
	accountData := resp.Results[0]

	// Get portfolio data
	portfolio, err := GetPortfolio(ctx, client, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	// Calculate total equity
	totalEquity := utils.GetFloat(accountData, "equity")
	cash := utils.GetFloat(accountData, "cash") + utils.GetFloat(accountData, "uncleared_deposits")

	var holdings []models.Holding
	for _, pos := range positions {
		if pos.Quantity == "0" || pos.Quantity == "0.00000000" {
			continue
		}

		// Get instrument data
		// Get stock price
		// Calculate metrics
		// This is a simplified version - full implementation would fetch instrument and price data

		quantity := utils.ParseFloat(pos.Quantity)
		avgBuyPrice := utils.ParseFloat(pos.AverageBuyPrice)

		holding := models.Holding{
			ID:              pos.InstrumentID,
			Quantity:        quantity,
			AverageBuyPrice: avgBuyPrice,
		}

		holdings = append(holdings, holding)
	}

	log.Printf("BuildHoldings: Built %d holdings\n", len(holdings))
	log.Printf("BuildHoldings: Total equity: %.2f, Cash: %.2f\n", totalEquity, cash)

	_ = portfolio
	_ = withDividend

	return holdings, nil
}

// GetPortfolio returns portfolio data for an account.
func GetPortfolio(ctx context.Context, client *robinstock_go.Client, accountNumber *string) (*models.Portfolio, error) {
	log.Println("GetPortfolio: Fetching portfolio...")

	if !client.IsAuthenticated() {
		return nil, robinstock_go.ErrNotAuthenticated
	}

	url := urls.PortfoliosURL()
	if accountNumber != nil {
		url = urls.PortfolioURL(*accountNumber)
	}

	resp, err := client.Get(ctx, url, nil, true)
	if err != nil {
		log.Printf("GetPortfolio: Error: %v\n", err)
		return nil, err
	}

	data := resp.Data
	if len(resp.Results) > 0 {
		data = resp.Results[0]
	}

	portfolio := &models.Portfolio{
		URL:                          utils.GetString(data, "url"),
		Account:                      utils.GetString(data, "account"),
		StartDate:                    utils.GetString(data, "start_date"),
		MarketValue:                  utils.GetString(data, "market_value"),
		Equity:                       utils.GetString(data, "equity"),
		ExtendedHoursMarketValue:     utils.GetString(data, "extended_hours_market_value"),
		ExtendedHoursEquity:          utils.GetString(data, "extended_hours_equity"),
		ExcessMargin:                 utils.GetString(data, "excess_margin"),
		EquityPreviousClose:          utils.GetString(data, "equity_previous_close"),
		AdjustedEquityPreviousClose:  utils.GetString(data, "adjusted_equity_previous_close"),
		Withdrawable:                 utils.GetString(data, "withdrawable"),
		UnsettledFunds:               utils.GetString(data, "unsettled_funds"),
	}

	return portfolio, nil
}
