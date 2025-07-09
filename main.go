package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("🚀 Advanced Desktop Entry Manager for Linux")
	fmt.Println("==========================================")

	config, err := loadConfiguration()
	if err != nil {
		fmt.Printf("❌ Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	userAppsDir := filepath.Join(os.Getenv("HOME"), ".local/share/applications")
	systemAppsDir := "/usr/share/applications"

	// Ensure user applications directory exists
	if err := os.MkdirAll(userAppsDir, 0755); err != nil {
		fmt.Printf("❌ Error creating user applications directory: %v\n", err)
		os.Exit(1)
	}

	// Show available flag sets
	fmt.Println("\n📋 Available Flag Sets:")
	for name, flagSet := range config.FlagSets {
		fmt.Printf("  • %s: %s\n", name, flagSet.Name)
		for _, flag := range flagSet.Flags {
			fmt.Printf("    %s\n", flag)
		}
	}

	for _, browser := range config.Browsers {
		fmt.Printf("\n🔍 Processing %s...\n", browser.Name)

		// Get combined flags for this browser
		combinedFlags := getCombinedFlags(browser.FlagSets, config.FlagSets)
		if len(combinedFlags) == 0 {
			fmt.Printf("⚠️  No flags configured for %s, skipping...\n", browser.Name)
			continue
		}

		fmt.Printf("🏃 Applying flags: %s\n", strings.Join(combinedFlags, " "))

		// Process main browser desktop files
		var allFilesToProcess []string
		allFilesToProcess = append(allFilesToProcess, browser.DesktopFiles...)

		// Process PWA files if patterns are defined
		if len(browser.PWAPatterns) > 0 {
			pwaFiles, err := findPWADesktopFiles(userAppsDir, browser.PWAPatterns, browser.ExcludePatterns)
			if err != nil {
				fmt.Printf("⚠️  Warning: Error finding %s PWA files: %v\n", browser.Name, err)
			} else if len(pwaFiles) == 0 {
				fmt.Printf("ℹ️  No PWA files found for %s\n", browser.Name)
			} else {
				fmt.Printf("🔗 Found %d PWA files for %s\n", len(pwaFiles), browser.Name)
				allFilesToProcess = append(allFilesToProcess, pwaFiles...)
			}
		}

		for _, desktopFile := range allFilesToProcess {
			processDesktopFile(systemAppsDir, userAppsDir, desktopFile, combinedFlags)
		}
	}

	fmt.Println("\n🎉 Desktop entry management completed!")
	fmt.Println("💡 Tip: Run this script after browser updates to keep entries synchronized")
	fmt.Println("🔧 Edit the browser configurations in ~/.config/lbdm/config.json to customize for your needs")
}
