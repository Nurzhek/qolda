package app

import (
	"net/http"
	"time"

	"qolda/internal/config"
)

// SystemHttpServerCreate builds the system HTTP server that exposes
// service endpoints: /healthcheck and /docs/* (swagger).
func SystemHttpServerCreate() *http.Server {
	mux := http.NewServeMux()

	// healthcheck
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// docs
	docFS := http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs")))
	mux.Handle("/docs/", docFS)

	return &http.Server{
		Addr:              ":" + config.Conf.SystemHttpPort,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       time.Minute,
		MaxHeaderBytes:    300 * 1024,
	}
}
