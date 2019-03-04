package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func (s *server) initRoutes() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	s.router.Route("/accounts", func(r chi.Router) {
		r.Get("/", s.listAccounts)
		r.Post("/", s.createAccount) // POST /accounts

		r.Route("/{accountID}", func(r chi.Router) {
			r.Use(s.accountCtx)               // Load the *Account on the request context
			r.Get("/", s.getAccount)          // GET /accounts/123
			r.Put("/", s.updateAccount)       // PUT /accounts/123
			r.Delete("/", s.deleteAccount)    // DELETE /accounts/123
			r.Post("/deposit", s.doDeposit)   // POST /accounts/123/deposit
			r.Post("/withdraw", s.doWithdraw) // POST /accounts/123/withdraw
			r.Post("/transfer", s.doTransfer) // POST /accounts/123/transfer
		})
	})

	s.router.Route("/payments", func(r chi.Router) {
		r.Get("/", s.listPayments)

		r.Route("/{paymentID}", func(r chi.Router) {
			r.Use(s.paymentCtx)
			r.Get("/", s.getPayment) // GET /payments/123
		})
	})
}
