package app

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/docgen"
	"github.com/vbauerster/fin-test/store"
)

type server struct {
	db              *store.MemDB
	router          *chi.Mux
	resourceClosers []io.Closer
}

func New(gendcod bool) *server {
	s := &server{
		db:     store.New(),
		router: chi.NewRouter(),
	}

	s.resourceClosers = append(s.resourceClosers, s.db)
	s.initRoutes()

	if gendcod {
		md := docgen.MarkdownRoutesDoc(s.router, docgen.MarkdownOpts{
			ProjectPath: "github.com/vbauerster/fin-test",
			Intro:       "fin-test REST API.",
		})
		if err := ioutil.WriteFile("routes.md", []byte(md), 0644); err != nil {
			log.Println(err)
		}
	}

	return s
}

func (s *server) Serve(ctx context.Context, addr string, shutdownTimeout time.Duration) {

	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	srv.RegisterOnShutdown(func() {
		for _, c := range s.resourceClosers {
			c.Close()
		}
	})

	idleConnsClosed := make(chan struct{})
	go func() {
		<-ctx.Done()
		tctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer func() {
			cancel()
			close(idleConnsClosed)
		}()
		if err := srv.Shutdown(tctx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Can't run server: %v", err)
	}
	<-idleConnsClosed
}
