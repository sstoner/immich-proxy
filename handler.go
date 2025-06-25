package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ImmichService struct {
	client *IMMICHClient
}

// AlbumHandler processes requests to /api/albums/id?key=
func (s *ImmichService) AlbumHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Handling album request: %s", r.URL.String())

	albumID := GetAlbumID(r)
	if albumID == "" {
		log.Printf("[ERROR] Invalid album ID in request: %s", r.URL.String())
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] Handling album request: %s", r.URL.String())

	albumInfo, err := s.client.GetAlbumInfo(albumID)
	if err != nil {
		log.Printf("[ERROR] Failed to get album info: %v", err)
		http.Error(w, "Failed to get album info", http.StatusInternalServerError)
		return
	}

	// return json response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(albumInfo); err != nil {
		log.Printf("[ERROR] Failed to encode album info: %v", err)
		http.Error(w, "Failed to encode album info", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Successfully handled album request: %s", r.URL.String())
	w.WriteHeader(http.StatusOK)
}

// SharedLinksHandler processes requests to /api/shared-links/me?key=
func (s *ImmichService) SharedLinksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Handling shared-links request: %s", r.URL.String())
	shareKey := GetShareKey(r)
	if shareKey == "" {
		log.Printf("[ERROR] Missing share key in request: %s", r.URL.String())
		http.Error(w, "Missing share key", http.StatusBadRequest)
		return
	}
	sharedLinksInfo, err := s.client.GetSharedLinksInfo(shareKey)
	if err != nil {
		log.Printf("[ERROR] Failed to get shared links info: %v", err)
		http.Error(w, "Failed to get shared links info", http.StatusInternalServerError)
		return
	}
	// return json response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sharedLinksInfo); err != nil {
		log.Printf("[ERROR] Failed to encode shared links info: %v", err)
		http.Error(w, "Failed to encode shared links info", http.StatusInternalServerError)
		return
	}
	log.Printf("[INFO] Successfully handled shared-links request: %s", r.URL.String())
	w.WriteHeader(http.StatusOK)

}

// AssetHandler processes requests to /api/assets/id?key=
func (s *ImmichService) AssetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Handling asset request: %s", r.URL.String())
	shareKey := GetShareKey(r)
	if shareKey == "" {
		log.Printf("[ERROR] Missing share key in request: %s", r.URL.String())
		http.Error(w, "Missing share key", http.StatusBadRequest)
		return
	}
	assetID := GetAssetID(r)
	if assetID == "" {
		log.Printf("[ERROR] Invalid asset ID in request: %s", r.URL.String())
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}
	log.Printf("[INFO] Handling asset request: %s", r.URL.String())
	assetInfo, err := s.client.GetAssetInfo(assetID, shareKey)
	if err != nil {
		log.Printf("[ERROR] Failed to get asset info: %v", err)
		http.Error(w, "Failed to get asset info", http.StatusInternalServerError)
		return
	}
	// return json response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(assetInfo); err != nil {
		log.Printf("[ERROR] Failed to encode asset info: %v", err)
		http.Error(w, "Failed to encode asset info", http.StatusInternalServerError)
		return
	}
	log.Printf("[INFO] Successfully handled asset request: %s", r.URL.String())
	w.WriteHeader(http.StatusOK)
}

// AssetThumbnailHandler processes requests to /api/assets/id/thumbnail?size=preview|thumbnail&key=
func (s *ImmichService) AssetThumbnailHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Handling asset thumbnail request: %s", r.URL.String())
	shareKey := GetShareKey(r)
	if shareKey == "" {
		log.Printf("[ERROR] Missing share key in request: %s", r.URL.String())
		http.Error(w, "Missing share key", http.StatusBadRequest)
		return
	}
	assetID := GetAssetID(r)
	if assetID == "" {
		log.Printf("[ERROR] Invalid asset ID in request: %s", r.URL.String())
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}
	size := GetAssetSize(r)
	if size == "" {
		log.Printf("[ERROR] Invalid asset size in request: %s", r.URL.String())
		http.Error(w, "Invalid asset size", http.StatusBadRequest)
		return
	}
	log.Printf("[INFO] Handling asset thumbnail request: %s", r.URL.String())
	thumbnail, err := s.client.GetAssetThumbnail(assetID, size)
	if err != nil {
		log.Printf("[ERROR] Failed to get asset thumbnail: %v", err)
		http.Error(w, "Failed to get asset thumbnail", http.StatusInternalServerError)
		return
	}
	// return image response
	w.Header().Set("Content-Type", "image/jpeg")
	if _, err := w.Write(thumbnail); err != nil {
		log.Printf("[ERROR] Failed to write asset thumbnail: %v", err)
		http.Error(w, "Failed to write asset thumbnail", http.StatusInternalServerError)
		return
	}
	log.Printf("[INFO] Successfully handled asset thumbnail request: %s", r.URL.String())
	w.WriteHeader(http.StatusOK)
}
