package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/matryer/is"
	"github.com/shopspring/decimal"
	"github.com/vbauerster/fin-test/app/payload"
	"github.com/vbauerster/fin-test/store"
)

func newTestServer() *httptest.Server {
	s := &server{
		db:     new(store.MemDB),
		router: chi.NewRouter(),
	}
	s.initRoutes()
	return httptest.NewServer(s.router)
}

func postAccount(url, payload string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPost, url+"/accounts", strings.NewReader(payload))
	return http.DefaultClient.Do(req)
}

func getAccount(url, id string) (*http.Response, error) {
	uri := fmt.Sprintf("%s%s%s", url, "/accounts/", id)
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	return http.DefaultClient.Do(req)
}

func doDeposit(url, accID string, amount decimal.Decimal) (*http.Response, error) {
	uri := fmt.Sprintf("%s%s%s%s", url, "/accounts/", accID, "/deposit")
	payload := fmt.Sprintf("{%q: %q}", "amount", amount)
	req, _ := http.NewRequest(http.MethodPost, uri, strings.NewReader(payload))
	return http.DefaultClient.Do(req)
}

func doTransfer(url, srcID, dstID string, amount decimal.Decimal) (*http.Response, error) {
	uri := fmt.Sprintf("%s%s%s%s", url, "/accounts/", srcID, "/transfer")
	payload := fmt.Sprintf("{%q: %q, %q: %q}",
		"amount",
		amount,
		"dst_account_id",
		dstID,
	)
	req, _ := http.NewRequest(http.MethodPost, uri, strings.NewReader(payload))
	return http.DefaultClient.Do(req)
}

func TestCreateAccountIgnoresBalance(t *testing.T) {
	is := is.New(t)
	ts := newTestServer()
	defer ts.Close()

	reqPayload := `{"code": 1, "name": "my acc", "balance": "100"}`
	res, err := postAccount(ts.URL, reqPayload)
	is.NoErr(err)
	is.Equal(res.StatusCode, http.StatusCreated)

	ar := new(payload.AccountResponse)
	is.NoErr(json.NewDecoder(res.Body).Decode(ar))
	is.NoErr(res.Body.Close())

	is.True(ar.ID != "")
	is.Equal(ar.CurrencyCode, "USD")
	is.True(ar.Balance.IsZero())
}

func TestDepositAccount(t *testing.T) {
	is := is.New(t)
	ts := newTestServer()
	defer ts.Close()

	reqPayload := `{"code": 1, "name": "my acc"}`
	res, err := postAccount(ts.URL, reqPayload)
	is.NoErr(err)
	is.Equal(res.StatusCode, http.StatusCreated)
	ar := new(payload.AccountResponse)
	is.NoErr(json.NewDecoder(res.Body).Decode(ar))
	is.NoErr(res.Body.Close())

	accID := ar.ID
	amount := decimal.New(10, 0)
	res, err = doDeposit(ts.URL, accID, amount)
	is.NoErr(err)
	is.Equal(res.StatusCode, http.StatusCreated)
	is.NoErr(res.Body.Close())

	res, err = getAccount(ts.URL, accID)
	is.NoErr(err)
	is.Equal(res.StatusCode, http.StatusOK)
	ar = new(payload.AccountResponse)
	is.NoErr(json.NewDecoder(res.Body).Decode(ar))
	is.NoErr(res.Body.Close())

	is.Equal(accID, ar.ID)
	is.Equal(ar.CurrencyCode, "USD")
	is.True(ar.Balance.Equal(amount))
}

func TestTransferAccount(t *testing.T) {
	is := is.New(t)
	ts := newTestServer()
	defer ts.Close()

	var accounts [2]string
	reqPayload := `{"code": 1, "name": "my acc"}`
	for i := 0; i < 2; i++ {
		res, err := postAccount(ts.URL, reqPayload)
		is.NoErr(err)
		is.Equal(res.StatusCode, http.StatusCreated)
		ar := new(payload.AccountResponse)
		is.NoErr(json.NewDecoder(res.Body).Decode(ar))
		is.NoErr(res.Body.Close())
		accounts[i] = ar.ID
	}

	amount := decimal.New(10, 0)
	res, err := doDeposit(ts.URL, accounts[0], amount)
	is.NoErr(err)
	is.Equal(res.StatusCode, http.StatusCreated)
	is.NoErr(res.Body.Close())

	res, err = doTransfer(ts.URL, accounts[0], accounts[1], amount)
	is.NoErr(err)
	is.Equal(res.StatusCode, http.StatusCreated)
	is.NoErr(res.Body.Close())

	res, err = getAccount(ts.URL, accounts[1])
	is.NoErr(err)
	is.Equal(res.StatusCode, http.StatusOK)
	ar := new(payload.AccountResponse)
	is.NoErr(json.NewDecoder(res.Body).Decode(ar))
	is.NoErr(res.Body.Close())

	is.Equal(accounts[1], ar.ID)
	is.True(ar.Balance.Equal(amount))
}
