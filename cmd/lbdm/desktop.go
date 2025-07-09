package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
)

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
