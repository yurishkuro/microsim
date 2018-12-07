package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yurishkuro/microsim/model"
)

// predefined includes hardcoded configurations.
var predefined = map[string]*model.Config{
	"hotrod": hotrod,
}

// Get loads a configuration by name.
func Get(nameOrPath string) (*model.Config, error) {
	if cfg, ok := predefined[nameOrPath]; ok {
		return cfg, nil
	}
	f, err := os.Open(nameOrPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %v", nameOrPath, err)
	}
	var cfg model.Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
