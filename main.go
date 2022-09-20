package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Andosius/valorant-exporter/helpers"
	"github.com/Andosius/valorant-exporter/models"
)

func main() {
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
		fmt.Println("Your configuration has been saved! Check", models.DATA_DIR, "folder.")

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
