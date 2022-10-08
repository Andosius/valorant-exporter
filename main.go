package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Andosius/valorant-exporter/cfg"
	"github.com/Andosius/valorant-exporter/helpers"
	"github.com/Andosius/valorant-exporter/models"
)

func main() {
	// Create data directories first - prevent io errors
	createDirectories()

	// Expect Valorant to be running an try to locate files and fetch data
	s := helpers.Selector{}

	s.SetInformation(
		[]string{
			"Welcome! This is YOUR Valorant settings manager!",
			"In order to continue, you have to choose between these options:",
		},
	)

	s.AddOption("Save your current settings (remember current time, it will be the filename)")
	s.AddOption("Push saved settings to the server")

	idx := s.RequestSelection()

	client := models.NewClient()

	switch idx {
	case 1:

		client.GetCurrentAccountSettings().WriteConfigToFile()
		fmt.Println("Your configuration has been saved! Check", cfg.DATA_DIR, "folder.")

	case 2:

		client.ConfigManager.LoadAllConfigurationFiles()

		s.Reset()
		s.SetInformation(
			[]string{
				"Please choose the config file you wish to upload!",
			},
		)

		for _, cfg := range client.ConfigManager.Configs {
			s.AddOption(cfg.Filename)
		}

		i := s.RequestSelection()
		client.PushConfigToServer((i - 1))

	default:

		fmt.Println("Invalid option - aborting...")
		time.Sleep(time.Second * 5)
		os.Exit(0)

	}
}

func createDirectories() {
	// Create data directory if it does not exist yet
	if !Exists(cfg.DATA_DIR) {
		err := os.Mkdir(cfg.DATA_DIR, cfg.PERMS)

		helpers.Fatal("main.createDirectories:1", err)
	}
}

// https://stackoverflow.com/a/22467409
func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
