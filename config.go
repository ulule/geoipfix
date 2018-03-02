package ipfix

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type corsConfig struct {
	AllowedOrigins   []string `json:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
	ExposedHeaders   []string `json:"exposed_headers"`
	MaxAge           int      `json:"max_age"`
}

type serverHTTPConfig struct {
	Port int        `json:"port"`
	Cors corsConfig `json:"cors"`
}

type serverConfig struct {
	HTTP serverHTTPConfig `json:"http"`
}

// Config is ipfix config
type Config struct {
	Debug        bool         `json:"debug"`
	DatabasePath string       `json:"database_path"`
	Server       serverConfig `json:"server"`
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

	if cfg.Server.HTTP.Port == 0 {
		cfg.Server.HTTP.Port = DefaultPort
	}

	return cfg, nil
}
