package main

import (
	"github.com/gorilla/mux"
)

// NewRouter creates and returns a mux.Router with all routes registered
func NewRouter(immichService *ImmichService) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(`/api/albums/{id:[^/]+}`, immichService.AlbumHandler).Methods("GET")
	r.HandleFunc(`/api/shared-links/me`, immichService.SharedLinksHandler).Methods("GET")
	r.HandleFunc(`/api/assets/{id:[^/]+}`, immichService.AssetHandler).Methods("GET")
	r.HandleFunc(`/api/assets/{id:[^/]+}/thumbnail`, immichService.MakeAssetHandler(
		[]string{"shareKey", "assetID", "size"},
		func(params map[string]string) ([]byte, error) {
			return immichService.client.GetAssetThumbnail(params["assetID"], params["size"], params["shareKey"])
		},
		"image/jpeg",
	)).Methods("GET")
	r.HandleFunc(`/api/assets/{id:[^/]+}/original`, immichService.MakeAssetHandler(
		[]string{"shareKey", "assetID"},
		func(params map[string]string) ([]byte, error) {
			return immichService.client.GetAssetOriginal(params["assetID"], params["shareKey"])
		},
		"image/jpeg",
	)).Methods("GET")

	r.PathPrefix("/").HandlerFunc(ProxyHandler)

	return r
}
