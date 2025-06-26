package main

import (
	"net/http"
	"strings"
)

func GetShareKey(r *http.Request) string {
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys) == 0 {
		return ""
	}
	return keys[0]
}

func GetAssetSize(r *http.Request) string {
	keys, ok := r.URL.Query()["size"]
	if !ok || len(keys) == 0 {
		return ""
	}
	size := keys[0]
	if size != "preview" && size != "thumbnail" {
		return ""
	}
	return size
}

func GetAlbumWithoutAssets(r *http.Request) bool {
	keys, ok := r.URL.Query()["withoutAssets"]
	if !ok || len(keys) == 0 {
		return false
	}
	return keys[0] == "true"
}

func GetAssetID(r *http.Request) string {
	// start with /api/assets/{id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) >= 3 && parts[1] == "api" && parts[2] == "assets" {
		return parts[3]
	}
	return ""
}

func GetAlbumID(r *http.Request) string {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) >= 3 && parts[1] == "api" && parts[2] == "albums" {
		return parts[3]
	}
	return ""
}
