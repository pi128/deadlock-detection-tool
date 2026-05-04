package deadlock

import "sort"

func FindCycles(graph map[string][]string) [][]string {
	var cycles [][]string
	seenCycles := map[string]bool{}
	color := map[string]int{}
	stackIndex := map[string]int{}
	var stack []string

	var visit func(string)
	visit = func(node string) {
		color[node] = 1
		stackIndex[node] = len(stack)
		stack = append(stack, node)

		neighbors := append([]string(nil), graph[node]...)
		sort.Strings(neighbors)
		for _, next := range neighbors {
			switch color[next] {
			case 0:
				visit(next)
			case 1:
				start := stackIndex[next]
				cycle := append([]string(nil), stack[start:]...)
				cycle = append(cycle, next)
				key := canonicalCycle(cycle)
				if !seenCycles[key] {
					seenCycles[key] = true
					cycles = append(cycles, cycle)
				}
			}
		}

		stack = stack[:len(stack)-1]
		delete(stackIndex, node)
		color[node] = 2
	}

	for _, node := range sortedGraphKeys(graph) {
		if color[node] == 0 {
			visit(node)
		}
	}
	return cycles
}

func canonicalCycle(cycle []string) string {
	if len(cycle) <= 1 {
		return ""
	}
	nodes := append([]string(nil), cycle[:len(cycle)-1]...)
	best := ""
	for i := range nodes {
		rotated := append([]string(nil), nodes[i:]...)
		rotated = append(rotated, nodes[:i]...)
		key := join(rotated)
		if best == "" || key < best {
			best = key
		}
	}
	return best
}
