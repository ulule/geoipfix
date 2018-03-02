package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

// Config is ipfix config
type Config struct {
	Port           int      `json:"port"`
	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	DatabasePath   string   `json:"database_path"`
}

// Load return a jsonq instance from a config path
func Load(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.Wrapf(err, "Config file %s cannot be loaded", path)
	}

	return LoadFromContent(string(content))
}

// LoadFromContent returns a jsonq instance from a config content
func LoadFromContent(content string) (*Config, error) {
	cfg := &Config{}
	dec := json.NewDecoder(strings.NewReader(content))
	err := dec.Decode(cfg)

	if err != nil {
		return nil, errors.Wrapf(err, "Config file %s cannot be parsed", content)
	}

	if cfg.DatabasePath == "" {
		cfg.DatabasePath = DatabaseURL
	}

	if cfg.Port == 0 {
		cfg.Port = DefaultPort
	}

	return cfg, nil
}
