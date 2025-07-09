package main

func getCombinedFlags(flagSetNames []string, flagSets map[string]FlagSet) []string {
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
