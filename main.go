package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FlagSet struct {
	Name  string
	Flags []string
}

type BrowserConfig struct {
	Name            string
	DesktopFiles    []string // Main desktop files (can be multiple for dev versions)
	PWAPatterns     []string // Patterns to match PWA files
	ExcludePatterns []string // Patterns to exclude (like Firefox PWAs)
	FlagSets        []string // Reference to flag sets to apply
	Description     string
}

// Global flag sets that can be mixed and matched
var flagSets = map[string]FlagSet{
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
}

// Browser configurations
var browsers = []BrowserConfig{
	{
		Name:            "Brave Browser",
		DesktopFiles:    []string{"brave-browser.desktop", "brave-browser-dev.desktop", "brave-browser-beta.desktop"},
		PWAPatterns:     []string{"brave-*.desktop"},
		ExcludePatterns: []string{"brave-browser*.desktop"}, // Exclude main browser files from PWA processing
		FlagSets:        []string{"touchpad_gestures", "wayland_full"},
		Description:     "Brave Browser (all variants)",
	},
	{
		Name:            "Microsoft Edge",
		DesktopFiles:    []string{"microsoft-edge.desktop", "microsoft-edge-dev.desktop", "microsoft-edge-beta.desktop"},
		PWAPatterns:     []string{"msedge-*.desktop"},
		ExcludePatterns: []string{"microsoft-edge*.desktop"}, // Exclude main browser files from PWA processing
		FlagSets:        []string{"touchpad_gestures", "edge_wayland"},
		Description:     "Microsoft Edge (all variants)",
	},
	{
		Name:            "Visual Studio Code",
		DesktopFiles:    []string{"code.desktop", "code-insiders.desktop"},
		PWAPatterns:     []string{}, // VSCode doesn't have PWAs
		ExcludePatterns: []string{},
		FlagSets:        []string{"wayland_basic"},
		Description:     "Visual Studio Code",
	},
	{
		Name:            "Opera",
		DesktopFiles:    []string{"opera.desktop", "opera-developer.desktop"},
		PWAPatterns:     []string{"opera-*.desktop"},
		ExcludePatterns: []string{"opera.desktop", "opera-developer.desktop"},
		FlagSets:        []string{"touchpad_gestures", "wayland_basic"},
		Description:     "Opera Browser",
	},
	{
		Name:            "Vivaldi",
		DesktopFiles:    []string{"vivaldi-stable.desktop", "vivaldi-beta.desktop"},
		PWAPatterns:     []string{"vivaldi-*.desktop"},
		ExcludePatterns: []string{"vivaldi-*.desktop"},
		FlagSets:        []string{"touchpad_gestures", "wayland_basic"},
		Description:     "Vivaldi Browser",
	},
}

func main() {
	fmt.Println("🚀 Advanced Desktop Entry Manager for Linux")
	fmt.Println("==========================================")

	userAppsDir := filepath.Join(os.Getenv("HOME"), ".local/share/applications")
	systemAppsDir := "/usr/share/applications"

	// Ensure user applications directory exists
	if err := os.MkdirAll(userAppsDir, 0755); err != nil {
		fmt.Printf("❌ Error creating user applications directory: %v\n", err)
		os.Exit(1)
	}

	// Show available flag sets
	fmt.Println("\n📋 Available Flag Sets:")
	for name, flagSet := range flagSets {
		fmt.Printf("  • %s: %s\n", name, flagSet.Name)
		for _, flag := range flagSet.Flags {
			fmt.Printf("    %s\n", flag)
		}
	}

	for _, browser := range browsers {
		fmt.Printf("\n🔍 Processing %s...\n", browser.Name)

		// Get combined flags for this browser
		combinedFlags := getCombinedFlags(browser.FlagSets)
		if len(combinedFlags) == 0 {
			fmt.Printf("⚠️  No flags configured for %s, skipping...\n", browser.Name)
			continue
		}

		fmt.Printf("🏃 Applying flags: %s\n", strings.Join(combinedFlags, " "))

		// Process main browser desktop files
		for _, desktopFile := range browser.DesktopFiles {
			if err := copyAndModifyDesktopFile(systemAppsDir, userAppsDir, desktopFile, combinedFlags); err != nil {
				fmt.Printf("⚠️  Warning: Could not process %s: %v\n", desktopFile, err)
			} else {
				fmt.Printf("✅ Updated: %s\n", desktopFile)
			}
		}

		// Process PWA files if patterns are defined
		if len(browser.PWAPatterns) > 0 {
			pwaFiles, err := findPWADesktopFiles(systemAppsDir, userAppsDir, browser.PWAPatterns, browser.ExcludePatterns)
			if err != nil {
				fmt.Printf("⚠️  Warning: Error finding %s PWA files: %v\n", browser.Name, err)
				continue
			}

			if len(pwaFiles) == 0 {
				fmt.Printf("ℹ️  No PWA files found for %s\n", browser.Name)
				continue
			}

			fmt.Printf("🔗 Found %d PWA files for %s\n", len(pwaFiles), browser.Name)

			for _, pwaFile := range pwaFiles {
				if err := copyAndModifyDesktopFile(systemAppsDir, userAppsDir, pwaFile, combinedFlags); err != nil {
					fmt.Printf("⚠️  Warning: Could not process PWA file %s: %v\n", pwaFile, err)
				} else {
					fmt.Printf("✅ Updated PWA: %s\n", pwaFile)
				}
			}
		}
	}

	fmt.Println("\n🎉 Desktop entry management completed!")
	fmt.Println("💡 Tip: Run this script after browser updates to keep entries synchronized")
	fmt.Println("🔧 Edit the browser configurations in the source code to customize for your needs")
}

func getCombinedFlags(flagSetNames []string) []string {
	var combinedFlags []string
	seen := make(map[string]bool)

	for _, flagSetName := range flagSetNames {
		if flagSet, exists := flagSets[flagSetName]; exists {
			for _, flag := range flagSet.Flags {
				if !seen[flag] {
					combinedFlags = append(combinedFlags, flag)
					seen[flag] = true
				}
			}
		}
	}

	return combinedFlags
}

func findPWADesktopFiles(systemDir, userDir string, patterns, excludePatterns []string) ([]string, error) {
	var pwaFiles []string

	// Check both system and user directories
	dirs := []string{systemDir, userDir}

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue // Skip if directory doesn't exist or can't be read
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()

			// Check if file matches any pattern
			matched := false
			for _, pattern := range patterns {
				if matched, _ = filepath.Match(pattern, name); matched {
					break
				}
			}

			if !matched {
				continue
			}

			// Check if file should be excluded
			excluded := false
			for _, excludePattern := range excludePatterns {
				if matched, _ = filepath.Match(excludePattern, name); matched {
					excluded = true
					break
				}
			}

			if excluded {
				continue
			}

			// Check if we already have this file
			alreadyAdded := false
			for _, existing := range pwaFiles {
				if existing == name {
					alreadyAdded = true
					break
				}
			}

			if !alreadyAdded {
				pwaFiles = append(pwaFiles, name)
			}
		}
	}

	return pwaFiles, nil
}

func copyAndModifyDesktopFile(srcDir, dstDir, filename string, flags []string) error {
	srcPath := filepath.Join(srcDir, filename)
	dstPath := filepath.Join(dstDir, filename)
	userHome, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home directory: %v", err)
	}
	userSrcPath := filepath.Join(userHome, ".local/share/applications", filename)

	// Check if source file exists, first in system dir, then in user dir
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		// If not in system dir, check user dir
		if _, userErr := os.Stat(userSrcPath); userErr == nil {
			srcPath = userSrcPath // It exists in user dir, use that as source
		} else {
			return fmt.Errorf("source file %s does not exist in %s or %s", filename, srcDir, filepath.Dir(userSrcPath))
		}
	}

	// Read the entire source file into memory first to avoid issues
	// when srcPath and dstPath are the same file.
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("error reading source file: %v", err)
	}

	// Create destination file (this will truncate the file if it exists)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer dstFile.Close()

	// Process file line by line from memory
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	writer := bufio.NewWriter(dstFile)
	defer writer.Flush()

	execPattern := regexp.MustCompile(`^Exec=(.*)$`)
	execLineCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is an Exec line (can appear in multiple sections)
		if matches := execPattern.FindStringSubmatch(line); matches != nil {
			execLineCount++
			execCmd := matches[1]

			// Remove any existing flags that we're about to add (to avoid duplicates)
			for _, flag := range flags {
				// Handle both with and without equals sign
				flagBase := strings.Split(flag, "=")[0]
				execCmd = removeFlagFromCommand(execCmd, flagBase)
			}

			// Find the executable part and add flags after it
			modifiedExecCmd := addFlagsToExecCommand(execCmd, flags)
			line = "Exec=" + modifiedExecCmd
		}

		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing to destination file: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading source file: %v", err)
	}

	if execLineCount > 0 {
		fmt.Printf("    Modified %d Exec lines\n", execLineCount)
	}

	return nil
}

func addFlagsToExecCommand(execCmd string, flags []string) string {
	if len(flags) == 0 {
		return execCmd
	}

	// Parse the command to find the executable and its arguments
	parts := strings.Fields(execCmd)
	if len(parts) == 0 {
		return execCmd
	}

	executable := parts[0]
	flagsStr := strings.Join(flags, " ")

	// Find where to insert flags - after executable but before other arguments
	// Look for the first argument that doesn't start with -- (likely a file/URL argument)
	// (Currently, flags are always inserted after executable)

	// Build the new command
	var newParts []string
	newParts = append(newParts, executable)
	newParts = append(newParts, strings.Fields(flagsStr)...)
	if len(parts) > 1 {
		newParts = append(newParts, parts[1:]...)
	}

	return strings.Join(newParts, " ")
}

func removeFlagFromCommand(command, flagName string) string {
	// Remove flag and its value (if any) from the command
	// This handles both --flag=value and --flag value formats

	// Pattern for --flag=value
	equalPattern := regexp.MustCompile(fmt.Sprintf(`\s*%s=[^\s]*`, regexp.QuoteMeta(flagName)))
	command = equalPattern.ReplaceAllString(command, "")

	// Pattern for --flag value (where value doesn't start with -)
	spacePattern := regexp.MustCompile(fmt.Sprintf(`\s*%s\s+[^\s-][^\s]*`, regexp.QuoteMeta(flagName)))
	command = spacePattern.ReplaceAllString(command, "")

	// Pattern for standalone --flag
	standalonePattern := regexp.MustCompile(fmt.Sprintf(`\s*%s(?:\s|$)`, regexp.QuoteMeta(flagName)))
	command = standalonePattern.ReplaceAllString(command, " ")

	return strings.TrimSpace(command)
}
