package main

import (
	"net/http"
	"strings"
)

/*
/api/albums/id?key=
/api/shared-links/me?key=
/api/assets/id&key=
/api/assets/id/thumbnail?size=preview|thumbnail&key=
*/

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
