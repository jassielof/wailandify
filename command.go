package main

import (
    "fmt"
    "regexp"
    "strings"
)

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
