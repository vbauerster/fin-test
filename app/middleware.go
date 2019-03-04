package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/vbauerster/fin-test/app/payload"
	"github.com/vbauerster/fin-test/model"
)

// accountCtx middleware is used to load an Account object from
// the URL parameters passed through as the request. In case
// the Account could not be found, we stop here and return a 404.
func (s *server) accountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var account model.Account
		var err error

		if accountID := chi.URLParam(r, "accountID"); accountID != "" {
			account, err = s.db.GetAccount(accountID)
		} else {
			render.Render(w, r, payload.ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, payload.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), payload.AccountCtxKey, &account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
