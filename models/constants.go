package models

var (
	BaseURL = "https://api.robinhood.com"
)

const (
	ClientID   = "c82SH0WZOsabOXGP2sxqcj34FxkvfnWRZBKlBjFS"
	ApiVersion = "1.431.4"
)

type OrderSide string
type OrderType string
type OrderTrigger string
type TimeInForce string

const (
	SideBuy  OrderSide = "buy"
	SideSell OrderSide = "sell"

	TypeMarket OrderType = "market"
	TypeLimit  OrderType = "limit"

	TriggerImmediate OrderTrigger = "immediate"
	TriggerStop      OrderTrigger = "stop"

	TIFGTC TimeInForce = "gtc"
	TIFGFD TimeInForce = "gfd"
	TIFIOC TimeInForce = "ioc"
	TIFOpg TimeInForce = "opg"
)

