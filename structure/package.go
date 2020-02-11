package structure

import (
	"encoding/json"
	"fmt"
)

// Package is the definition for one package in composer.lock file
type Package struct {
	Name            string          `json:"name"`
	Version         string          `json:"version"`
	Source          json.RawMessage `json:"source,omitempty"`
	Dist            json.RawMessage `json:"dist,omitempty"`
	Require         json.RawMessage `json:"require,omitempty"`
	Provide         json.RawMessage `json:"provide,omitempty"`
	Comflict        json.RawMessage `json:"conflict,omitempty"`
	Replace         json.RawMessage `json:"replace,omitempty"`
	RequireDev      json.RawMessage `json:"require-dev,omitempty"`
	Suggest         json.RawMessage `json:"suggest,omitempty"`
	Bin             json.RawMessage `json:"bin,omitempty"`
	Type            string          `json:"type,omitempty"`
	Extra           json.RawMessage `json:"extra,omitempty"`
	Autoload        json.RawMessage `json:"autoload,omitempty"`
	NotificationUrl string          `json:"notification-url,omitempty"`
	License         json.RawMessage `json:"license,omitempty"`
	Authors         json.RawMessage `json:"authors,omitempty"`
	Description     string          `json:"description,omitempty"`
	Homepage        string          `json:"homepage,omitempty"`
	Keywords        json.RawMessage `json:"keywords,omitempty"`
	Abandoned       json.RawMessage `json:"abandoned,omitempty"`
	Time            string          `json:"time,omitempty"`
}

func (p Package) String() string {
	s, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	return string(s) + ","
}
