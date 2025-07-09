package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AppConfig represents the top-level structure of the config file.
type AppConfig struct {
	FlagSets map[string]FlagSet `json:"flagSets"`
	Browsers []BrowserConfig    `json:"browsers"`
}

type FlagSet struct {
	Name  string   `json:"name"`
	Flags []string `json:"flags"`
}

type BrowserConfig struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	DesktopFiles    []string `json:"desktopFiles"`
	PWAPatterns     []string `json:"pwaPatterns"`
	ExcludePatterns []string `json:"excludePatterns"`
	FlagSets        []string `json:"flagSets"`
}

// getDefaultConfig returns the hardcoded default configuration.
func getDefaultConfig() AppConfig {
	return AppConfig{
		FlagSets: map[string]FlagSet{
			"wayland_basic": {
				Name:  "Basic Wayland Support",
				Flags: []string{"--ozone-platform=wayland"},
			},
			"wayland_full": {
				Name: "Full Wayland Support",
				Flags: []string{
					"--ozone-platform=wayland",
					"--enable-features=UseOzonePlatform,WaylandWindowDecorations",
					"--gtk-version=4",
					"--ozone-platform-hint=auto",
					"--disable-features=GlobalShortcutsPortal",
				},
			},
			"touchpad_gestures": {
				Name:  "Touchpad Gestures",
				Flags: []string{"--enable-features=TouchpadOverscrollHistoryNavigation"},
			},
			"edge_wayland": {
				Name:  "Edge Wayland Fix",
				Flags: []string{"--ozone-platform=wayland"},
			},
		},
		Browsers: []BrowserConfig{
			{
				Name:            "Brave Browser",
				DesktopFiles:    []string{"brave-browser.desktop", "brave-browser-dev.desktop", "brave-browser-beta.desktop"},
				PWAPatterns:     []string{"brave-*.desktop"},
				ExcludePatterns: []string{"brave-browser*.desktop"},
				FlagSets:        []string{"touchpad_gestures", "wayland_full"},
				Description:     "Brave Browser (all variants)",
			},
			{
				Name:            "Microsoft Edge",
				DesktopFiles:    []string{"microsoft-edge.desktop", "microsoft-edge-dev.desktop", "microsoft-edge-beta.desktop"},
				PWAPatterns:     []string{"msedge-*.desktop"},
				ExcludePatterns: []string{"microsoft-edge*.desktop"},
				FlagSets:        []string{"touchpad_gestures", "edge_wayland"},
				Description:     "Microsoft Edge (all variants)",
			},
			{
				Name:         "Visual Studio Code",
				DesktopFiles: []string{"code.desktop", "code-insiders.desktop"},
				FlagSets:     []string{"wayland_basic"},
				Description:  "Visual Studio Code",
			},
		},
	}
}

// loadConfiguration reads the config from a file, creating a default one if it doesn't exist.
func loadConfiguration() (AppConfig, error) {
	var config AppConfig

	configDir, err := os.UserConfigDir()
	if err != nil {
		return config, fmt.Errorf("could not get user config directory: %w", err)
	}

	appConfigDir := filepath.Join(configDir, "lbdm")
	configPath := filepath.Join(appConfigDir, "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("💡 No config file found. Creating a default one at: %s\n", configPath)
		defaultConfig := getDefaultConfig()

		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			return config, fmt.Errorf("could not create config directory: %w", err)
		}

		file, err := os.Create(configPath)
		if err != nil {
			return config, fmt.Errorf("could not create config file: %w", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(defaultConfig); err != nil {
			return config, fmt.Errorf("could not write default config: %w", err)
		}
		return defaultConfig, nil
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("could not read config file: %w", err)
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return config, fmt.Errorf("could not parse config file: %w", err)
	}

	return config, nil
}
