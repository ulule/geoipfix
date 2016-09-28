package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jmoiron/jsonq"
)

// Load return a jsonq instance from a config path
func Load(path string) (*jsonq.JsonQuery, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("Config file %s cannot be loaded: %s", path, err)
	}

	return LoadFromContent(string(content))
}

// LoadFromContent returns a jsonq instance from a config content
func LoadFromContent(content string) (*jsonq.JsonQuery, error) {
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(content))
	err := dec.Decode(&data)

	if err != nil {
		return nil, fmt.Errorf("Config file %s cannot be parsed: %s", content, err)
	}

	return jsonq.NewQuery(data), nil
}
