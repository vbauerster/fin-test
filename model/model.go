package model

import (
	"time"

	"github.com/shopspring/decimal"
)

//go:generate stringer -type=CurrencyCode -trimprefix=Currency
type CurrencyCode int

const (
	_ CurrencyCode = iota
	CurrencyUSD
	CurrencyEUR
	CurrencyRUB
)

//go:generate stringer -type=PaymentCode -trimprefix=Payment
type PaymentCode int

const (
	_ PaymentCode = iota
	PaymentDEPOSIT
	PaymentWITHDRAW
	PaymentTRANSFER
)

type Account struct {
	ID      string          `json:"id"`
	Code    CurrencyCode    `json:"code"`
	Name    string          `json:"name"`
	Balance decimal.Decimal `json:"balance"`
}

type Payment struct {
	ID           string          `json:"id"`
	Code         PaymentCode     `json:"code"`
	SrcAccountID string          `json:"src_account_id"`
	DstAccountID string          `json:"dst_account_id,omitempty"`
	Amount       decimal.Decimal `json:"amount"`
	Date         time.Time       `json:"date"`
}
