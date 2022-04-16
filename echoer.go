package echoer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func NewHTTPServer(addr string, version string) *http.Server {
	handler := newHTTPHandler(version)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return srv
}

func newHTTPHandler(version string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", handleIndex())
	mux.Handle("/500", handleFiveHundred())
	handler := addVersionHeader(mux, version)
	return handler
}

func addVersionHeader(next http.Handler, version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Version", version)
		next.ServeHTTP(w, r)
	})
}

func handleIndex() http.HandlerFunc {
	type response struct {
		Hostname string      `json:"hostname"`
		Header   http.Header `json:"header"`
		Path     string      `json:"path"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		hostname, err := os.Hostname()
		if err != nil {
			httpError(w, err, http.StatusInternalServerError)
		}
		resp := response{
			Hostname: hostname,
			Header:   r.Header,
			Path:     r.URL.Path,
		}
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			httpError(w, err, http.StatusInternalServerError)
		}
	}
}

func handleFiveHundred() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpError(w, errors.New("error on purpose"), http.StatusInternalServerError)
	}
}

func httpError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	fmt.Fprint(os.Stderr, err)
}
