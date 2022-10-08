package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Andosius/valorant-exporter/cfg"
	"github.com/Andosius/valorant-exporter/helpers"
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
	filepath.Walk(cfg.DATA_DIR, func(path string, info os.FileInfo, err error) error {
		helpers.Fatal("cm.LoadAllConfigurationFiles:2", err)

		// Skip directories, they don't matter to us. :)
		if !info.IsDir() {
			body, err := os.ReadFile(cfg.DATA_DIR + cfg.SEP + info.Name())
			helpers.Fatal("cm.LoadAllConfigurationFiles:3", err)

			var cfg_file Config
			err = json.Unmarshal(body, &cfg_file)

			helpers.Fatal("cm.LoadAllConfigurationFiles:4", err)

			cfg_file.Filename = info.Name()

			cm.Configs = append(cm.Configs, cfg_file)
		}
		return nil
	})
}

func (cfg_file Config) WriteConfigToFile() {
	// Marshal struct to JSON
	body, err := json.Marshal(cfg_file)
	helpers.Fatal("cfg.WriteConfigToFile:1", err)

	// Since we already loaded configs, we don't have to check for data dir
	t := time.Now()
	filename := fmt.Sprintf(
		"%02d_%02d_%02d-%02d_%02d_%02d",
		t.Day(), t.Month(), t.Year(),
		t.Hour(), t.Minute(), t.Second(),
	)

	path := cfg.DATA_DIR + cfg.SEP + filename + ".json"
	err = os.WriteFile(path, body, cfg.PERMS)

	helpers.Fatal("cfg.WriteConfigToFile:2", err)
}
