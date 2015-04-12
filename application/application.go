package application

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"strings"
)

type Application struct {
	Jq *jsonq.JsonQuery
}

func NewFromConfigPath(path string) (*Application, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("Your config file %s cannot be loaded: %s", path, err)
	}

	return NewFromConfig(string(content))
}

func NewFromConfig(content string) (*Application, error) {
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(content))
	err := dec.Decode(&data)

	if err != nil {
		return nil, fmt.Errorf("Your config file %s cannot be parsed: %s", content, err)
	}

	jq := jsonq.NewQuery(data)

	return NewFromJsonQuery(jq)
}

func (a *Application) Port() int {
	port, err := a.Jq.Int("port")

	if err != nil {
		port = DefaultPort
	}

	return port
}

func NewFromJsonQuery(jq *jsonq.JsonQuery) (*Application, error) {
	app := &Application{}
	app.Jq = jq

	return app, nil
}
