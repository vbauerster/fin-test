package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/render"
	"github.com/vbauerster/fin-test/app/payload"
	"github.com/vbauerster/fin-test/model"
	"github.com/vbauerster/fin-test/store"
)

func (s *server) listAccounts(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, payload.NewAccountListResponse(s.db.GetAccounts())); err != nil {
		render.Render(w, r, payload.ErrRender(err))
		return
	}
}

func (s *server) listPayments(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, payload.NewPaymentListResponse(s.db.GetPayments())); err != nil {
		render.Render(w, r, payload.ErrRender(err))
		return
	}
}

// createAccount persists the posted Account and returns it
// back to the client as an acknowledgement.
// Account.Balance is ignored on purpose.
func (s *server) createAccount(w http.ResponseWriter, r *http.Request) {
	data := &payload.AccountRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payload.ErrInvalidRequest(err))
		return
	}

	// Account.Balance is ignored on purpose, client should use deposit.
	// data.Balance = data.ProtectedBalance

	data.Code = data.ProtectedCode

	switch data.Code {
	case model.CurrencyUSD, model.CurrencyEUR, model.CurrencyRUB:
		id, err := s.db.NewAccount(*data.Account)
		if err != nil {
			render.Render(w, r, payload.ErrInternal(err))
			return
		}

		account, err := s.db.GetAccount(id)
		if err != nil {
			render.Render(w, r, payload.ErrInternal(err))
			return
		}

		render.Status(r, http.StatusCreated)
		render.Render(w, r, payload.NewAccountResponse(&account))
	default:
		render.Render(w, r, payload.ErrInvalidRequest(
			fmt.Errorf("unsupported currency code: %d", data.Code),
		))
	}
}

// getAccount returns the specific Account. You'll notice it just
// fetches the Account right off the context, as its understood that
// if we made it this far, the Account must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func (s *server) getAccount(w http.ResponseWriter, r *http.Request) {
	ctxAccount := r.Context().Value(payload.AccountCtxKey).(*model.Account)

	if err := render.Render(w, r, payload.NewAccountResponse(ctxAccount)); err != nil {
		render.Render(w, r, payload.ErrRender(err))
		return
	}
}

// updateAccount updates an existing Account in our persistent store.
// It will not update Balance and Code fields on purpose.
func (s *server) updateAccount(w http.ResponseWriter, r *http.Request) {
	ctxAccount := r.Context().Value(payload.AccountCtxKey).(*model.Account)
	fmt.Printf("upd: %s\n", spew.Sdump(ctxAccount))

	data := &payload.AccountRequest{Account: ctxAccount}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payload.ErrInvalidRequest(err))
		return
	}

	fmt.Printf("upd: %s\n", spew.Sdump(data))
	account, err := s.db.UpdateAccount(data.Account.ID, *data.Account)
	if err != nil {
		render.Render(w, r, payload.ErrInternal(err))
	}

	render.Render(w, r, payload.NewAccountResponse(&account))
}

// deleteAccount removes an existing Account from our persistent store.
func (s *server) deleteAccount(w http.ResponseWriter, r *http.Request) {
	ctxAccount := r.Context().Value(payload.AccountCtxKey).(*model.Account)

	account, err := s.db.RemoveAccount(ctxAccount.ID)
	if err != nil {
		render.Render(w, r, payload.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, payload.NewAccountResponse(&account))
}

func (s *server) doDeposit(w http.ResponseWriter, r *http.Request) {
	ctxAccount := r.Context().Value(payload.AccountCtxKey).(*model.Account)

	data := &payload.PaymentRequest{
		Payment: &model.Payment{
			Code:         model.PaymentDEPOSIT,
			SrcAccountID: ctxAccount.ID,
			Date:         time.Now(),
		},
	}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payload.ErrInvalidRequest(err))
		return
	}
	spew.Dump(data)

	var id string
	err := s.db.DoTransaction(func(tx *store.MemDB) (err error) {
		payment := data.Payment
		account, err := tx.GetAccount(payment.SrcAccountID)
		if err != nil {
			return err
		}

		account.Balance = account.Balance.Add(payment.Amount)
		_, err = tx.UpdateAccount(account.ID, account)
		if err != nil {
			return err
		}
		id, err = tx.NewPayment(*payment)
		return err
	})

	if err != nil {
		if err == store.ErrNotFound {
			render.Render(w, r, payload.ErrInvalidRequest(err))
			return
		}
		render.Render(w, r, payload.ErrInternal(err))
		return
	}

	payment, err := s.db.GetPayment(id)
	if err != nil {
		render.Render(w, r, payload.ErrInternal(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, payload.NewPaymentResponse(&payment))
}

func (s *server) doWithdraw(w http.ResponseWriter, r *http.Request) {
	ctxAccount := r.Context().Value(payload.AccountCtxKey).(*model.Account)

	data := &payload.PaymentRequest{
		Payment: &model.Payment{
			Code:         model.PaymentWITHDRAW,
			SrcAccountID: ctxAccount.ID,
			Date:         time.Now(),
		},
	}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payload.ErrInvalidRequest(err))
		return
	}
	spew.Dump(data)

	var id string
	var rd render.Renderer
	err := s.db.DoTransaction(func(tx *store.MemDB) (err error) {
		payment := data.Payment
		account, err := tx.GetAccount(payment.SrcAccountID)
		if err != nil {
			return err
		}

		if payment.Amount.GreaterThan(account.Balance) {
			err = fmt.Errorf("insufficient funds of account id %q: %v %s", payment.SrcAccountID, account.Balance, account.Code)
			rd = payload.ErrInvalidRequest(err)
			return err
		}

		account.Balance = account.Balance.Sub(payment.Amount)
		_, err = tx.UpdateAccount(account.ID, account)
		if err != nil {
			return err
		}
		id, err = tx.NewPayment(*payment)
		return err
	})

	if err != nil {
		if rd != nil {
			render.Render(w, r, rd)
			return
		}
		render.Render(w, r, payload.ErrInternal(err))
		return
	}

	payment, err := s.db.GetPayment(id)
	if err != nil {
		render.Render(w, r, payload.ErrInternal(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, payload.NewPaymentResponse(&payment))
}

func (s *server) doTransfer(w http.ResponseWriter, r *http.Request) {
	account := r.Context().Value(payload.AccountCtxKey).(*model.Account)

	data := &payload.PaymentRequest{
		Payment: &model.Payment{
			Code:         model.PaymentTRANSFER,
			SrcAccountID: account.ID,
			Date:         time.Now(),
		},
	}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payload.ErrInvalidRequest(err))
		return
	}

	data.DstAccountID = data.ProtectedDstAccountID

	var rd render.Renderer
	var id string
	err := s.db.DoTransaction(func(tx *store.MemDB) (err error) {
		payment := data.Payment
		var aa [2]model.Account
		for i, id := range [...]string{payment.SrcAccountID, payment.DstAccountID} {
			aa[i], err = tx.GetAccount(id)
			if err != nil {
				if err == store.ErrNotFound {
					rd = payload.ErrInvalidRequest(err)
				}
				return err
			}
		}
		if aa[0].Balance.IsNegative() || payment.Amount.GreaterThan(aa[0].Balance) {
			err = fmt.Errorf("insufficient funds of account id %q: %v %s", payment.SrcAccountID, aa[0].Balance, aa[0].Code)
			rd = payload.ErrInvalidRequest(err)
			return err
		}
		if aa[0].Code != aa[1].Code {
			err = fmt.Errorf("conversion %s-%s not supported", aa[0].Code, aa[1].Code)
			rd = payload.ErrInvalidRequest(err)
			return err
		}

		aa[0].Balance = aa[0].Balance.Sub(payment.Amount)
		aa[1].Balance = aa[1].Balance.Add(payment.Amount)
		for i := range aa {
			_, err = tx.UpdateAccount(aa[i].ID, aa[i])
			if err != nil {
				return err
			}
		}
		id, err = tx.NewPayment(*payment)
		return err
	})
	if err != nil {
		if rd != nil {
			render.Render(w, r, rd)
			return
		}
		render.Render(w, r, payload.ErrInternal(err))
		return
	}

	payment, err := s.db.GetPayment(id)
	if err != nil {
		render.Render(w, r, payload.ErrInternal(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, payload.NewPaymentResponse(&payment))
}
