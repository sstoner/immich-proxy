package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type AlbumsKeys struct {
	ApiKeys       []string
	syncEnabled   bool              // whether to fetch albums asynchronously
	immageBaseURL string            // base URL for Immich API
	AlbumsKeys    map[string]string // map of album ID to API key
	lock          sync.Mutex        // to protect AlbumsKeys
}

func NewAlbumsKeys(keys []string, syncEnabled bool, immageBaseURL string) *AlbumsKeys {
	return &AlbumsKeys{
		ApiKeys:       keys,
		AlbumsKeys:    make(map[string]string),
		syncEnabled:   syncEnabled,
		immageBaseURL: immageBaseURL,
	}
}

func (a *AlbumsKeys) setAlbumKey(albumId, key string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	log.Debugf("Setting album key for album %s to %s", albumId, key)
	if key == "" {
		delete(a.AlbumsKeys, albumId)
	} else {
		a.AlbumsKeys[albumId] = key
	}
}

func (a *AlbumsKeys) GetAlbumKey(albumId string) string {
	key := a.getAlbumKeyFromMap(albumId)
	if key == "" && !a.syncEnabled {
		log.Debugf("Album key for %s not found in map, fetching without sync", albumId)
		key = a.GetAlbumKeyWithoutSync(albumId)
		if key == "" {
			log.Warnf("No API key found for album %s", albumId)
			return ""
		}
	} else if key == "" {
		log.Debugf("Album key for %s found in map: %s", albumId, key)
	} else {
		log.Debugf("Using cached album key for %s: %s", albumId, key)
	}

	return key
}

func (a *AlbumsKeys) GetAlbumKeyWithoutSync(albumId string) string {
	for _, key := range a.ApiKeys {
		log.Debugf("Fetching album %s with API key %s", albumId, key)
		albums, err := a.getAlbums(a.immageBaseURL, key)
		if err != nil {
			log.Errorf("Failed to fetch albums for API key %s: %v", key, err)
			continue
		}
		if slices.Contains(albums, albumId) {
			a.setAlbumKey(albumId, key)
			return key
		}
	}

	return ""
}

func (a *AlbumsKeys) getAlbumKeyFromMap(albumId string) string {
	a.lock.Lock()
	defer a.lock.Unlock()

	key, ok := a.AlbumsKeys[albumId]
	if !ok {
		return ""
	}
	return key
}

func (a *AlbumsKeys) fetchAllAlbums(baseUrl string) {
	var wg sync.WaitGroup
	log.Debugf("Fetching all albums from %s with %d API keys", baseUrl, len(a.ApiKeys))
	for _, key := range a.ApiKeys {
		wg.Add(1)
		go func(apiKey string) {
			defer wg.Done()
			albums, err := a.getAlbums(baseUrl, apiKey)
			if err != nil {
				log.Errorf("Failed to fetch albums for API key %s: %v", apiKey, err)
				return
			}
			for _, albumId := range albums {
				a.setAlbumKey(albumId, apiKey)
			}
		}(key)
	}
	wg.Wait()
}

func (a *AlbumsKeys) StartRefreshing(ctx context.Context,
	refreshInterval time.Duration, immichUrl string) {
	if !a.syncEnabled {
		log.Warn("Albums sync is disabled, not starting refresh")
		return
	}

	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()
	log.Infof("Starting albums refresh every %s", refreshInterval)
	a.fetchAllAlbums(immichUrl) // Initial fetch
	go func() {
		for {
			select {
			case <-ticker.C:
				a.fetchAllAlbums(immichUrl)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (a *AlbumsKeys) getAlbums(immichUrl, key string) ([]string, error) {
	endpoint := "/albums"
	url := fmt.Sprintf("%s/api%s", immichUrl, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("x-api-key", key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warnf("failed to close response body: %v", err)
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("immich api error on endpoint %s: %d %s: %s", url, resp.StatusCode, resp.Status, string(b))
	}

	type AlbumsResponse []struct {
		ID string `json:"id"`
	}

	var albumsResp AlbumsResponse

	err = json.NewDecoder(resp.Body).Decode(&albumsResp)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	result := make([]string, len(albumsResp))
	for i, album := range albumsResp {
		result[i] = album.ID
	}
	return result, nil
}
