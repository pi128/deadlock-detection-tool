package deadlock

import (
	"fmt"
	"sort"
	"strings"
)

func splitPair(key string) (string, string, error) {
	parts := strings.Split(key, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("key %q must use process:resource format", key)
	}
	return parts[0], parts[1], nil
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedGraphKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedIntKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func copyIntMap(m map[string]int) map[string]int {
	cp := make(map[string]int, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

func join(values []string) string {
	return strings.Join(values, ", ")
}

func formatCycle(cycle []string) string {
	return strings.Join(cycle, " -> ")
}
