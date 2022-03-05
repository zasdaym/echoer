package echoer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func NewHTTPServer(addr string, version string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", indexHandler())
	mux.Handle("/500", fiveHundredHandler())

	srv := &http.Server{
		Addr:    addr,
		Handler: versionMiddleware(mux, version),
	}
	return srv
}

func versionMiddleware(next http.Handler, version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Version", version)
		next.ServeHTTP(w, r)
	})
}

func indexHandler() http.HandlerFunc {
	type response struct {
		Hostname string      `json:"hostname"`
		Header   http.Header `json:"header"`
		Path     string      `json:"path"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusInternalServerError)
	}
}
