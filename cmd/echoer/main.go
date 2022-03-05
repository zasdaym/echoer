package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/zasdaym/echoer"
)

var version string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var (
		showVersion = flag.Bool("version", false, "Show version")
		listenAddr  = flag.String("listen-addr", ":8080", "HTTP server listen address")
	)
	flag.Parse()

	if *showVersion {
		fmt.Print(version)
		return nil
	}

	srv := echoer.NewHTTPServer(*listenAddr, version)
	go func(ctx context.Context, srv *http.Server) {
		<-ctx.Done()
		shutdownCtx, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		srv.Shutdown(shutdownCtx)
	}(ctx, srv)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed && err != nil {
		return err
	}
	return nil
}
