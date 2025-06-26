package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// ProxyHandler handles all requests that do not match specific Immich endpoints.
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the proxy logic to forward requests to Immich
	log.Infof("Proxying request: %s", r.URL.Path)
	w.WriteHeader(http.StatusNotImplemented)
	if _, err := w.Write([]byte("Proxy not implemented yet")); err != nil {
		log.Warnf("failed to write proxy response: %v", err)
	}
}
