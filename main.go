package main

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

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

	r := NewRouter(immichService)

	log.Infof("[INFO] Immich Proxy Server started on %s", cfg.Listen)
	if err := http.ListenAndServe(cfg.Listen, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
