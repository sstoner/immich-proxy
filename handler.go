package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ImmichService struct {
	client *IMMICHClient
}

// AlbumHandler processes requests to /api/albums/id?key=
func (s *ImmichService) AlbumHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Handling album request: %s", r.URL.String())

	albumID := GetAlbumID(r)
	if albumID == "" {
		log.Errorf("Invalid album ID in request: %s", r.URL.String())
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	withoutAssets := GetAlbumWithoutAssets(r)
	albumInfo, err := s.client.GetAlbumInfo(albumID, withoutAssets)
	if err != nil {
		log.Errorf("Failed to get album info: %v", err)
		http.Error(w, "Failed to get album info", http.StatusInternalServerError)
		return
	}

	// return json response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(albumInfo); err != nil {
		log.Errorf("Failed to encode album info: %v", err)
		http.Error(w, "Failed to encode album info", http.StatusInternalServerError)
		return
	}

	log.Debugf("Successfully handled album request: %s", r.URL.String())
}

// SharedLinksHandler processes requests to /api/shared-links/me?key=
func (s *ImmichService) SharedLinksHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Handling shared-links request: %s", r.URL.String())
	shareKey := GetShareKey(r)
	if shareKey == "" {
		log.Errorf("Missing share key in request: %s", r.URL.String())
		http.Error(w, "Missing share key", http.StatusBadRequest)
		return
	}
	sharedLinksInfo, err := s.client.GetSharedLinksInfo(shareKey)
	if err != nil {
		log.Errorf("Failed to get shared links info: %v", err)
		http.Error(w, "Failed to get shared links info", http.StatusInternalServerError)
		return
	}
	// return json response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sharedLinksInfo); err != nil {
		log.Errorf("Failed to encode shared links info: %v", err)
		http.Error(w, "Failed to encode shared links info", http.StatusInternalServerError)
		return
	}
	log.Debugf("Successfully handled shared-links request: %s", r.URL.String())
}

// AssetHandler processes requests to /api/assets/id?key=
func (s *ImmichService) AssetHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Handling asset request: %s", r.URL.String())
	shareKey := GetShareKey(r)
	if shareKey == "" {
		log.Debugf("Missing share key in request: %s", r.URL.String())
		http.Error(w, "Missing share key", http.StatusBadRequest)
		return
	}
	assetID := GetAssetID(r)
	if assetID == "" {
		log.Debugf("Invalid asset ID in request: %s", r.URL.String())
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	assetInfo, err := s.client.GetAssetInfo(assetID, shareKey)
	if err != nil {
		log.Errorf("Failed to get asset info: %v", err)
		http.Error(w, "Failed to get asset info", http.StatusInternalServerError)
		return
	}
	// return json response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(assetInfo); err != nil {
		log.Errorf("Failed to encode asset info: %v", err)
		http.Error(w, "Failed to encode asset info", http.StatusInternalServerError)
		return
	}
	log.Debugf("Successfully handled asset request: %s", r.URL.String())
}

// AssetThumbnailHandler processes requests to /api/assets/id/thumbnail?size=preview|thumbnail&key=
func (s *ImmichService) AssetThumbnailHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Handling asset thumbnail request: %s", r.URL.String())
	shareKey := GetShareKey(r)
	if shareKey == "" {
		log.Errorf("Missing share key in request: %s", r.URL.String())
		http.Error(w, "Missing share key", http.StatusBadRequest)
		return
	}
	assetID := GetAssetID(r)
	if assetID == "" {
		log.Debugf("Invalid asset ID in request: %s", r.URL.String())
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}
	size := GetAssetSize(r)
	if size == "" {
		log.Debugf("Invalid asset size in request: %s", r.URL.String())
		http.Error(w, "Invalid asset size", http.StatusBadRequest)
		return
	}

	thumbnail, err := s.client.GetAssetThumbnail(assetID, size, shareKey)
	if err != nil {
		log.Errorf("Failed to get asset thumbnail: %v", err)
		http.Error(w, "Failed to get asset thumbnail", http.StatusInternalServerError)
		return
	}
	// return image response
	w.Header().Set("Content-Type", "image/jpeg")
	if _, err := w.Write(thumbnail); err != nil {
		log.Errorf("Failed to write asset thumbnail: %v", err)
		http.Error(w, "Failed to write asset thumbnail", http.StatusInternalServerError)
		return
	}
	log.Debugf("Successfully handled asset thumbnail request: %s", r.URL.String())
}

func (s *ImmichService) AssetOriginalHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Handling asset original request: %s", r.URL.String())
	shareKey := GetShareKey(r)
	if shareKey == "" {
		log.Errorf("Missing share key in request: %s", r.URL.String())
		http.Error(w, "Missing share key", http.StatusBadRequest)
		return
	}
	assetID := GetAssetID(r)
	if assetID == "" {
		log.Debugf("Invalid asset ID in request: %s", r.URL.String())
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	original, err := s.client.GetAssetOriginal(assetID, shareKey)
	if err != nil {
		log.Errorf("Failed to get asset original: %v", err)
		http.Error(w, "Failed to get asset original", http.StatusInternalServerError)
		return
	}
	// return image response
	w.Header().Set("Content-Type", "image/jpeg")
	if _, err := w.Write(original); err != nil {
		log.Errorf("Failed to write asset original: %v", err)
		http.Error(w, "Failed to write asset original", http.StatusInternalServerError)
		return
	}
	log.Debugf("Successfully handled asset original request: %s", r.URL.String())
}

func (s *ImmichService) MakeAssetHandler(
	paramKeys []string,
	getData func(params map[string]string) ([]byte, error),
	contentType string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("Handling %s request: %s", contentType, r.URL.String())
		params, ok := requireParams(w, r, paramKeys...)
		if !ok {
			return
		}
		data, err := getData(params)
		if err != nil {
			log.Errorf("Failed to get %s: %v", contentType, err)
			http.Error(w, "Failed to get "+contentType, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", contentType)
		if _, err := w.Write(data); err != nil {
			log.Errorf("Failed to write %s: %v", contentType, err)
			http.Error(w, "Failed to write "+contentType, http.StatusInternalServerError)
			return
		}
		log.Debugf("Successfully handled %s request: %s", contentType, r.URL.String())
	}
}

func requireParams(w http.ResponseWriter, r *http.Request, keys ...string) (map[string]string, bool) {
	params := make(map[string]string)
	for _, key := range keys {
		var val string
		switch key {
		case "shareKey":
			val = GetShareKey(r)
		case "assetID":
			val = GetAssetID(r)
		case "size":
			val = GetAssetSize(r)
		}
		if val == "" {
			log.Warnf("Missing or invalid %s in request: %s", key, r.URL.String())
			http.Error(w, "Missing or invalid "+key, http.StatusBadRequest)
			return nil, false
		}
		params[key] = val
	}
	return params, true
}
