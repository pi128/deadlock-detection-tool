package main

import (
	"sort"
	"strings"
)

func sortedIntKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func join(values []string) string {
	if len(values) == 0 {
		return "(none)"
	}
	return strings.Join(values, ", ")
}

func formatCycle(cycle []string) string {
	return strings.Join(cycle, " -> ")
}
