package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Andosius/valorant-exporter/helpers"
)

const (
	DATA_DIR = "data"
	PERMS    = 0660
	SEP      = string(os.PathSeparator)
)

type Config struct {
	Type     string `json:"type"`
	Data     string `json:"data"`
	Modified int    `json:"modified"`
	Filename string `json:"-"`
}

type ConfigManager struct {
	Configs []Config
}

func (cm *ConfigManager) LoadAllConfigurationFiles() {
	// Create data directory if it does not exist yet
	if !Exists(DATA_DIR) {
		err := os.Mkdir(DATA_DIR, PERMS)

		helpers.Fatal("cm.LoadAllConfigurationFiles:1", err)
	} else {
		filepath.Walk(DATA_DIR, func(path string, info os.FileInfo, err error) error {
			helpers.Fatal("cm.LoadAllConfigurationFiles:2", err)

			// Skip directories, they don't matter to us. :)
			if !info.IsDir() {
				body, err := os.ReadFile(DATA_DIR + SEP + info.Name())
				helpers.Fatal("cm.LoadAllConfigurationFiles:3", err)

				var cfg Config
				err = json.Unmarshal(body, &cfg)

				helpers.Fatal("cm.LoadAllConfigurationFiles:4", err)

				cfg.Filename = info.Name()

				cm.Configs = append(cm.Configs, cfg)
			}
			return nil
		})
	}
}

func (cfg Config) WriteConfigToFile() {
	// Marshal struct to JSON
	body, err := json.Marshal(cfg)
	helpers.Fatal("cfg.WriteConfigToFile:1", err)

	// Since we already loaded configs, we don't have to check for data dir
	t := time.Now()
	filename := fmt.Sprintf(
		"%02d_%02d_%02d-%02d_%02d_%02d",
		t.Day(), t.Month(), t.Year(),
		t.Hour(), t.Minute(), t.Second(),
	)

	path := DATA_DIR + SEP + filename + ".json"
	err = os.WriteFile(path, body, PERMS)

	helpers.Fatal("cfg.WriteConfigToFile:2", err)
}

// https://stackoverflow.com/a/22467409
func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
