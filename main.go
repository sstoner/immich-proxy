package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Immich struct {
		URL                   string   `yaml:"url"`
		APIKeys               []string `yaml:"api_keys"`
		AlbumsRefreshInterval string   `yaml:"albumsRefreshInterval,omitempty"`
	} `yaml:"immich"`
	Listen   string `yaml:"listen"`
	LogLevel string `yaml:"logLevel"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Warnf("failed to close config file: %v", err)
		}
	}()
	var cfg Config
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	log.SetLevel(logLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	log.Infof("proxy started, Listen %s, forwarded to %s", cfg.Listen, cfg.Immich.URL)

	refreshInterval, err := time.ParseDuration(cfg.Immich.AlbumsRefreshInterval)
	if err != nil {
		log.Fatalf("Invalid albumsRefreshInterval: %v", err)
	}
	albumsKeys := NewAlbumsKeys(cfg.Immich.APIKeys)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	albumsKeys.StartRefreshing(ctx, refreshInterval, cfg.Immich.URL)

	immichService := &ImmichService{
		client: NewIMMICHClient(cfg.Immich.URL, albumsKeys),
	}

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

	log.Infof("[INFO] Immich Proxy Server started on %s", cfg.Listen)
	if err := http.ListenAndServe(cfg.Listen, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
