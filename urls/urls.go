package urls

import (
	"fmt"

	"github.com/yourusername/robinstock_go/models"
)

func LoginURL() string {
	return models.BaseURL + "/oauth2/token/"
}

func ChallengeURL(challengeID string) string {
	return fmt.Sprintf("https://api.robinhood.com/challenge/%s/respond/", challengeID)
}

func PathfinderUserMachineURL() string {
	return "https://api.robinhood.com/pathfinder/user_machine/"
}

func SheriffInquiryURL(machineID string) string {
	return fmt.Sprintf("https://api.robinhood.com/pathfinder/inquiries/%s/user_view/", machineID)
}

func SheriffChallengeStatusURL(challengeID string) string {
	return fmt.Sprintf("https://api.robinhood.com/push/%s/get_prompts_status/", challengeID)
}

func PhoenixURL() string {
	return "https://phoenix.robinhood.com/accounts/unified"
}

func AccountsURL() string {
	return "https://api.robinhood.com/accounts/"
}

func AccountURL(accountNumber string) string {
	return fmt.Sprintf("https://api.robinhood.com/accounts/%s/", accountNumber)
}

func PositionsURL() string {
	return "https://api.robinhood.com/positions/"
}

func PortfoliosURL() string {
	return "https://api.robinhood.com/portfolios/"
}

func PortfolioURL(accountNumber string) string {
	return fmt.Sprintf("https://api.robinhood.com/portfolios/%s/", accountNumber)
}

func DividendsURL() string {
	return "https://api.robinhood.com/dividends/"
}

func NotificationsURL(tracker bool) string {
	if tracker {
		return "https://api.robinhood.com/midlands/notifications/notification_tracker/"
	}
	return "https://api.robinhood.com/notifications/devices/"
}

func BankTransfersURL() string {
	return "https://api.robinhood.com/ach/transfers/"
}

func LinkedURL() string {
	return "https://api.robinhood.com/ach/relationships/"
}

func InstrumentsURL() string {
	return "https://api.robinhood.com/instruments/"
}

func QuotesURL() string {
	return "https://api.robinhood.com/quotes/"
}

func FundamentalsURL(symbol string) string {
	return fmt.Sprintf("https://api.robinhood.com/fundamentals/%s/", symbol)
}

func HistoricalsURL() string {
	return "https://api.robinhood.com/quotes/historicals/"
}

func BasicProfileURL() string {
	return "https://api.robinhood.com/user/basic_info/"
}

func InvestmentProfileURL() string {
	return "https://api.robinhood.com/user/investment_profile/"
}

func SecurityProfileURL() string {
	return "https://api.robinhood.com/user/additional_info/"
}

func UserProfileURL() string {
	return "https://api.robinhood.com/user/"
}

func MarketsURL() string {
	return "https://api.robinhood.com/markets/"
}

func MarketHoursURL(market, date string) string {
	return fmt.Sprintf("https://api.robinhood.com/markets/%s/hours/%s/", market, date)
}

func MoversSP500URL() string {
	return "https://api.robinhood.com/midlands/movers/sp500/"
}

func Top100MostPopularURL() string {
	return "https://api.robinhood.com/midlands/tags/tag/100-most-popular/"
}

func MarketCategoryURL(category string) string {
	return fmt.Sprintf("https://api.robinhood.com/midlands/tags/tag/%s/", category)
}

func OrdersURL(orderID, accountNumber, startDate *string) string {
	url := "https://api.robinhood.com/orders/"
	if orderID != nil {
		url += *orderID + "/"
	}
	params := ""
	if accountNumber != nil {
		params += "account_numbers=" + *accountNumber
	}
	if startDate != nil {
		if params != "" {
			params += "&"
		}
		params += "updated_at[gte]=" + *startDate
	}
	if params != "" {
		url += "?" + params
	}
	return url
}

func CancelURL(orderID string) string {
	return fmt.Sprintf("https://api.robinhood.com/orders/%s/cancel/", orderID)
}

func OptionOrdersURL(orderID, accountNumber, startDate *string) string {
	url := "https://api.robinhood.com/options/orders/"
	if orderID != nil {
		url += *orderID + "/"
	}
	params := ""
	if accountNumber != nil {
		params += "account_numbers=" + *accountNumber
	}
	if startDate != nil {
		if params != "" {
			params += "&"
		}
		params += "updated_at[gte]=" + *startDate
	}
	if params != "" {
		url += "?" + params
	}
	return url
}

func OptionCancelURL(id string) string {
	return fmt.Sprintf("https://api.robinhood.com/options/orders/%s/cancel/", id)
}
