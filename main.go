package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/vbauerster/fin-test/app"
)

var gendoc = flag.Bool("gendoc", false, "Generate router documentation")

func main() {
	flag.Parse()
	server := app.New(*gendoc)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		defer signal.Stop(quit)
		<-quit
		cancel()
	}()

	server.Serve(ctx, ":3333", 5*time.Second)
}
