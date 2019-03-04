package store

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/vbauerster/fin-test/model"
)

type StoreError string

func (s StoreError) Error() string {
	return string(s)
}

const (
	ErrNotFound = StoreError("record not found")
)

type MemDB struct {
	// we could use RWmutex, but just Mutex is sufficient for this small
	// struct. Refer to https://github.com/golang/go/issues/17973 before
	// using RWmutex in any prod.
	sync.Mutex
	accounts []model.Account
	payments []model.Payment
	_idCount uint64
}

var db *MemDB = &MemDB{
	accounts: []model.Account{
		{ID: "1", Code: model.CurrencyUSD, Balance: decimal.NewFromFloat(136.61), Name: "My USD primary account"},
		{ID: "2", Code: model.CurrencyUSD, Balance: decimal.NewFromFloat(136.62), Name: "My USD secondary account"},
		{ID: "3", Code: model.CurrencyUSD, Balance: decimal.NewFromFloat(136.63), Name: "My USD spare account"},
		{ID: "4", Code: model.CurrencyEUR, Balance: decimal.NewFromFloat(136.64), Name: "My EUR account"},
		{ID: "5", Code: model.CurrencyRUB, Balance: decimal.NewFromFloat(136.65), Name: "My RUB account"},
	},
}

func init() {
	db._idCount = uint64(len(db.accounts))
}

func New() *MemDB {
	return db
}

// Close implement io.Closer
func (db *MemDB) Close() error {
	return nil
}

func (db *MemDB) NewAccount(account model.Account) (id string, err error) {
	db.Lock()
	defer db.Unlock()
	return db.newAccount(account)
}

func (db *MemDB) newAccount(account model.Account) (id string, err error) {
	db._idCount++
	account.ID = strconv.FormatUint(db._idCount, 10)
	db.accounts = append(db.accounts, account)
	return account.ID, nil
}

func (db *MemDB) GetAccounts() []model.Account {
	db.Lock()
	defer db.Unlock()
	return db.getAccounts()
}

func (db *MemDB) getAccounts() []model.Account {
	accounts := make([]model.Account, len(db.accounts))
	copy(accounts, db.accounts)
	return accounts
}

func (db *MemDB) GetAccount(id string) (model.Account, error) {
	db.Lock()
	defer db.Unlock()
	return db.getAccount(id)
}

func (db *MemDB) getAccount(id string) (account model.Account, err error) {
	for _, a := range db.accounts {
		if a.ID == id {
			return a, nil
		}
	}
	return account, ErrNotFound
}

func (db *MemDB) UpdateAccount(id string, account model.Account) (model.Account, error) {
	db.Lock()
	defer db.Unlock()
	return db.updateAccount(id, account)
}

func (db *MemDB) updateAccount(id string, account model.Account) (result model.Account, err error) {
	for i, a := range db.accounts {
		if a.ID == id {
			account.ID = id
			db.accounts[i] = account
			return account, nil
		}
	}
	return result, ErrNotFound
}

func (db *MemDB) RemoveAccount(id string) (model.Account, error) {
	db.Lock()
	defer db.Unlock()
	return db.removeAccount(id)
}

func (db *MemDB) removeAccount(id string) (result model.Account, err error) {
	for i, a := range db.accounts {
		if a.ID == id {
			db.accounts = append((db.accounts)[:i], (db.accounts)[i+1:]...)
			return a, nil
		}
	}
	return result, errors.New("account not found.")
}

func (db *MemDB) PaymentTransaction(payment model.Payment) (err error) {
	db.Lock()
	defer db.Unlock()

	var aa [2]model.Account
	for i, id := range [...]string{payment.SrcAccountID, payment.DstAccountID} {
		aa[i], err = db.getAccount(id)
		if err != nil {
			return err
		}
	}
	if aa[0].Balance.IsNegative() || payment.Amount.GreaterThan(aa[0].Balance) {
		return fmt.Errorf("insufficient funds of account id %s: %s %s", payment.SrcAccountID, aa[0].Balance, aa[0].Code)
	}
	if aa[0].Code != aa[1].Code {
		return fmt.Errorf("currency conversion %s-%s not supported", aa[0].Code, aa[1].Code)
	}

	aa[0].Balance = aa[0].Balance.Sub(payment.Amount)
	aa[1].Balance = aa[1].Balance.Add(payment.Amount)
	payment.Date = time.Now()
	db.newPayment(payment)
	return err
}

func (db *MemDB) DoTransaction(cb func(*MemDB) error) error {
	db.Lock()
	defer db.Unlock()

	tx := &MemDB{
		accounts: make([]model.Account, len(db.accounts)),
		payments: make([]model.Payment, len(db.payments)),
		_idCount: db._idCount,
	}

	copy(tx.accounts, db.accounts)
	copy(tx.payments, db.payments)

	if err := cb(tx); err != nil {
		return err
	}

	db.accounts = tx.accounts
	db.payments = tx.payments
	db._idCount = tx._idCount
	return nil
}

func (db *MemDB) NewPayment(payment model.Payment) (id string, err error) {
	db.Lock()
	defer db.Unlock()
	return db.newPayment(payment)
}

func (db *MemDB) newPayment(payment model.Payment) (id string, err error) {
	db._idCount++
	payment.ID = strconv.FormatUint(db._idCount, 10)
	db.payments = append(db.payments, payment)
	return payment.ID, nil
}

func (db *MemDB) GetPayments() []model.Payment {
	db.Lock()
	defer db.Unlock()
	return db.getPayments()
}

func (db *MemDB) getPayments() []model.Payment {
	payments := make([]model.Payment, len(db.payments))
	copy(payments, db.payments)
	return payments
}

func (db *MemDB) GetPayment(id string) (model.Payment, error) {
	db.Lock()
	defer db.Unlock()
	return db.getPayment(id)
}

func (db *MemDB) getPayment(id string) (payment model.Payment, err error) {
	for _, p := range db.payments {
		if p.ID == id {
			return p, nil
		}
	}
	return payment, ErrNotFound
}
