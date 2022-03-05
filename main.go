package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
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

	mux := http.NewServeMux()
	mux.Handle("/", indexHandler())
	mux.Handle("/500", fiveHundredHandler())
	srv := &http.Server{
		Addr:    *listenAddr,
		Handler: mux,
	}

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

func indexHandler() http.HandlerFunc {
	type response struct {
		Hostname string      `json:"hostname"`
		Header   http.Header `json:"header"`
		Path     string      `json:"path"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Version", version)
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Fprintf(os.Stderr, "get hostname: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := response{
			Hostname: hostname,
			Header:   r.Header,
			Path:     r.URL.Path,
		}
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			fmt.Fprintf(os.Stderr, "encode response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func fiveHundredHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Version", version)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
