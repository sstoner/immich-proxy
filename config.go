package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Immich struct {
		URL                   string   `yaml:"url"`
		APIKeys               []string `yaml:"api_keys"`
		AlbumsSyncEnabled     bool     `yaml:"albumsSyncEnabled,omitempty"`
		AlbumsRefreshInterval string   `yaml:"albumsRefreshInterval,omitempty"`
	} `yaml:"immich"`
	Listen   string     `yaml:"listen"`
	LogLevel string     `yaml:"logLevel"`
	Cors     CORSConfig `yaml:"cors,omitempty"`
}

type CORSConfig struct {
	AllowOrigin      string `yaml:"allowOrigin"`
	AllowMethods     string `yaml:"allowMethods"`
	AllowHeaders     string `yaml:"allowHeaders"`
	AllowCredentials bool   `yaml:"allowCredentials"`
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

func (c *Config) GetCORSConfig() *CORSConfig {
	return &c.Cors
}
