package main

import (
	"log"
	"net/http"
)

// ProxyHandler handles all requests that do not match specific Immich endpoints.
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the proxy logic to forward requests to Immich
	log.Printf("[INFO] Proxying request: %s", r.URL.Path)
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Proxy not implemented yet"))
}
