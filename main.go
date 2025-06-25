package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Immich struct {
		URL    string `yaml:"url"`
		APIKey string `yaml:"api_key"`
	} `yaml:"immich"`
	Listen string `yaml:"listen"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
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
	log.Printf("proxy started，Listen %s，forwared to %s", cfg.Listen, cfg.Immich.URL)

	immichService := &ImmichService{
		client: NewIMMICHClient(cfg.Immich.URL, cfg.Immich.APIKey),
	}

	r := mux.NewRouter()

	r.HandleFunc(`/api/albums/{id:[^/]+}`, immichService.AlbumHandler).Methods("GET")
	r.HandleFunc(`/api/shared-links/me`, immichService.SharedLinksHandler).Methods("GET")
	r.HandleFunc(`/api/assets/{id:[^/]+}`, immichService.AssetHandler).Methods("GET")
	r.HandleFunc(`/api/assets/{id:[^/]+}/thumbnail`, immichService.AssetThumbnailHandler).Methods("GET")

	r.PathPrefix("/").HandlerFunc(ProxyHandler)

	log.Println("[INFO] Immich Proxy Server started on", cfg.Listen)
	if err := http.ListenAndServe(cfg.Listen, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
