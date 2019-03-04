package payload

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/shopspring/decimal"
	"github.com/vbauerster/fin-test/model"
)

var AccountCtxKey = &contextKey{"Account"}

type PaymentRequest struct {
	*model.Payment

	// following protect embeded Payment's fields from overriding
	ProtectedID           string            `json:"id"`
	ProtectedSrcAccountID string            `json:"src_account_id"`
	ProtectedDstAccountID string            `json:"dst_account_id"`
	ProtectedCode         model.PaymentCode `json:"code"`
	ProtectedDate         time.Time         `json:"date"`
}

func (s *PaymentRequest) Bind(r *http.Request) error {
	if s.Payment == nil {
		return errors.New("missing required fields")
	}
	if s.Payment.SrcAccountID == s.DstAccountID {
		return errors.New("nop payment, you're rich")
	}
	if s.Amount.IsZero() || s.Amount.IsNegative() {
		return errors.New("bad amount")
	}
	return nil
}

type AccountRequest struct {
	*model.Account

	// following protect embeded Account's fields from overriding
	ProtectedID      string             `json:"id"`
	ProtectedCode    model.CurrencyCode `json:"code"`
	ProtectedBalance decimal.Decimal    `json:"balance"`
}

func (a *AccountRequest) Bind(r *http.Request) error {
	// a.Account is nil if no Account fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Account == nil {
		return errors.New("missing required fields")
	}

	switch a.Code {
	case model.CurrencyUSD, model.CurrencyEUR, model.CurrencyRUB:
		return nil
	default:
		return fmt.Errorf("unsupported currency_id: %d", a.Code)
	}
}

// ArticleResponse is the response payload for the Article data model.
// See NOTE above in ArticleRequest as well.
//
// In the ArticleResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type AccountResponse struct {
	*model.Account

	// We add an additional field to the response here.. such as this
	CurrencyCode string `json:"currency_code"`
}

func NewAccountResponse(account *model.Account) *AccountResponse {
	resp := &AccountResponse{Account: account}
	return resp
}

func (rd *AccountResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.CurrencyCode = rd.Code.String()
	return nil
}

func NewAccountListResponse(accounts []model.Account) []render.Renderer {
	var list []render.Renderer
	for i := range accounts {
		list = append(list, NewAccountResponse(&accounts[i]))
	}
	return list
}

type PaymentResponse struct {
	*model.Payment

	PaymentCode string `json:"payment_code"`
}

func NewPaymentResponse(payment *model.Payment) *PaymentResponse {
	resp := &PaymentResponse{Payment: payment}
	return resp
}

func (rd *PaymentResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.PaymentCode = rd.Code.String()
	return nil
}

func NewPaymentListResponse(payments []model.Payment) []render.Renderer {
	var list []render.Renderer
	for i := range payments {
		list = append(list, NewPaymentResponse(&payments[i]))
	}
	return list
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "payload context value " + k.name
}
