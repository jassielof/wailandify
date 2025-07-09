package main

import (
	"bufio"
	"bytes"
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

// findPWADesktopFiles scans only the user's application directory for PWA desktop files.
func findPWADesktopFiles(userDir string, patterns, excludePatterns []string) ([]string, error) {
	var pwaFiles []string

	entries, err := os.ReadDir(userDir)
	if err != nil {
		// If the directory doesn't exist, it's not an error, just no files found.
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not read user applications directory: %v", err)
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

		pwaFiles = append(pwaFiles, name)
	}

	return pwaFiles, nil
}

// processDesktopFile handles the logic for a single desktop file:
// 1. Ensures a copy exists in the user directory.
// 2. Modifies the user copy to apply flags.
// 3. Skips writing if the file is already up-to-date.
func processDesktopFile(systemAppsDir, userAppsDir, filename string, flags []string) {
	srcPath := filepath.Join(systemAppsDir, filename)
	dstPath := filepath.Join(userAppsDir, filename)

	// If the destination file doesn't exist, copy it from the system directory.
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			// Source doesn't exist either, so we can't do anything.
			// This is a common case (e.g., beta/dev browser not installed), so don't show a warning.
			return
		}

		input, err := os.ReadFile(srcPath)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to read system file %s: %v\n", filename, err)
			return
		}

		if err := os.WriteFile(dstPath, input, 0644); err != nil {
			fmt.Printf("⚠️  Warning: Failed to copy %s to user directory: %v\n", filename, err)
			return
		}
		fmt.Printf("📋 Copied %s to user directory\n", filename)
	}

	// Read the destination file (user's copy).
	content, err := os.ReadFile(dstPath)
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to read user file %s: %v\n", filename, err)
		return
	}

	// Process the content and generate the new version.
	modifiedContent, modifiedCount, err := modifyDesktopContent(content, flags)
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not process %s: %v\n", filename, err)
		return
	}

	// If content is unchanged, do nothing.
	if bytes.Equal(content, modifiedContent) {
		fmt.Printf("✅ Up-to-date: %s\n", filename)
		return
	}

	// Write the modified content back to the destination file.
	if err := os.WriteFile(dstPath, modifiedContent, 0644); err != nil {
		fmt.Printf("⚠️  Warning: Failed to write updated file %s: %v\n", filename, err)
		return
	}

	fmt.Printf("✅ Updated %s (%d Exec lines modified)\n", filename, modifiedCount)
}

// modifyDesktopContent takes the content of a desktop file and returns the modified version.
func modifyDesktopContent(content []byte, flags []string) ([]byte, int, error) {
	var out bytes.Buffer
	scanner := bufio.NewScanner(bytes.NewReader(content))
	execPattern := regexp.MustCompile(`^(Exec=)(.*)$`)
	modifiedCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		if matches := execPattern.FindStringSubmatch(line); len(matches) > 2 {
			modifiedCount++
			execCmd := matches[2]
			modifiedExecCmd := updateExecCommand(execCmd, flags)
			line = "Exec=" + modifiedExecCmd
		}
		if _, err := out.WriteString(line + "\n"); err != nil {
			return nil, 0, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return out.Bytes(), modifiedCount, nil
}

// parseCommandLine splits a command string into arguments, correctly handling quoted sections.
func parseCommandLine(command string) ([]string, error) {
	var args []string
	var currentArg strings.Builder
	inQuotes := false
	var quoteChar rune

	for i, r := range command {
		if inQuotes {
			if r == quoteChar {
				inQuotes = false
			} else {
				currentArg.WriteRune(r)
			}
		} else {
			switch r {
			case '"', '\'':
				inQuotes = true
				quoteChar = r
			case ' ':
				if currentArg.Len() > 0 {
					args = append(args, currentArg.String())
					currentArg.Reset()
				}
			default:
				currentArg.WriteRune(r)
			}
		}

		// Special case for the very last character in the command
		if i == len(command)-1 && currentArg.Len() > 0 {
			args = append(args, currentArg.String())
		}
	}

	if inQuotes {
		return nil, fmt.Errorf("unclosed quote in command: %s", command)
	}

	return args, nil
}

// updateExecCommand intelligently adds missing flags to an Exec command string.
// It handles quoted arguments and avoids duplicating existing flags.
func updateExecCommand(execCmd string, flagsToAdd []string) string {
	if len(flagsToAdd) == 0 {
		return execCmd
	}

	trimmedCmd := strings.TrimSpace(execCmd)
	parts, err := parseCommandLine(trimmedCmd)
	if err != nil {
		// Fallback to old behavior if parsing fails
		fmt.Printf("⚠️  Warning: Could not parse command line '%s': %v. Using simple split.\n", execCmd, err)
		parts = strings.Fields(trimmedCmd)
	}

	if len(parts) == 0 {
		return ""
	}

	executable := parts[0]
	args := parts[1:]

	// 2. Separate existing arguments into flags and non-flags (like %U, URLs).
	var existingFlags []string
	var otherArgs []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			existingFlags = append(existingFlags, arg)
		} else {
			otherArgs = append(otherArgs, arg)
		}
	}

	// 3. Determine which flags to add, avoiding duplicates.
	finalFlags := make([]string, len(existingFlags))
	copy(finalFlags, existingFlags)
	existingSet := make(map[string]bool)
	for _, f := range existingFlags {
		existingSet[f] = true
	}

	for _, flagToAdd := range flagsToAdd {
		if !existingSet[flagToAdd] {
			finalFlags = append(finalFlags, flagToAdd)
			existingSet[flagToAdd] = true // In case flagsToAdd has duplicates
		}
	}

	// 4. Reconstruct the command string.
	var newParts []string
	newParts = append(newParts, executable)
	newParts = append(newParts, finalFlags...)
	newParts = append(newParts, otherArgs...)

	// Re-quote any arguments that contain spaces
	for i, part := range newParts {
		if strings.Contains(part, " ") && !strings.HasPrefix(part, `"`) {
			newParts[i] = `"` + part + `"`
		}
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
