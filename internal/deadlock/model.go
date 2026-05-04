package deadlock

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

type System struct {
	Resources   map[string]int `json:"resources"`
	Allocations map[string]int `json:"allocations"`
	Requests    map[string]int `json:"requests"`
}

type Edge struct {
	From string
	To   string
}

type DetectionResult struct {
	Available       map[string]int
	WaitForGraph    map[string][]string
	Edges           []Edge
	Deadlocked      bool
	Cycles          [][]string
	Safe            bool
	SafeSequence    []string
	Blocked         []string
	Recommendations []string
}

func Load(r io.Reader) (*System, error) {
	var s System
	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return nil, err
	}
	if s.Resources == nil {
		s.Resources = map[string]int{}
	}
	if s.Allocations == nil {
		s.Allocations = map[string]int{}
	}
	if s.Requests == nil {
		s.Requests = map[string]int{}
	}
	return &s, s.Validate()
}

func (s *System) Validate() error {
	for id, total := range s.Resources {
		if total < 0 {
			return fmt.Errorf("resource %q has negative total", id)
		}
	}
	for key, amount := range s.Allocations {
		p, r, err := splitPair(key)
		if err != nil {
			return err
		}
		if amount < 0 {
			return fmt.Errorf("allocation %q has negative amount", key)
		}
		if _, ok := s.Resources[r]; !ok {
			return fmt.Errorf("allocation %q references unknown resource %q", key, r)
		}
		if p == "" {
			return fmt.Errorf("allocation %q has empty process", key)
		}
	}
	for key, amount := range s.Requests {
		p, r, err := splitPair(key)
		if err != nil {
			return err
		}
		if amount < 0 {
			return fmt.Errorf("request %q has negative amount", key)
		}
		if _, ok := s.Resources[r]; !ok {
			return fmt.Errorf("request %q references unknown resource %q", key, r)
		}
		if p == "" {
			return fmt.Errorf("request %q has empty process", key)
		}
	}
	available := s.Available()
	for r, amount := range available {
		if amount < 0 {
			return fmt.Errorf("resource %q is over-allocated by %d instance(s)", r, -amount)
		}
	}
	return nil
}

func (s *System) Detect() DetectionResult {
	available := s.Available()
	graph, edges := s.WaitForGraph(available)
	cycles := FindCycles(graph)
	safe, sequence, blocked := s.SafeSequence(available)
	return DetectionResult{
		Available:       available,
		WaitForGraph:    graph,
		Edges:           edges,
		Deadlocked:      len(cycles) > 0,
		Cycles:          cycles,
		Safe:            safe,
		SafeSequence:    sequence,
		Blocked:         blocked,
		Recommendations: Recommendations(cycles, blocked, graph),
	}
}

func (s *System) Processes() []string {
	seen := map[string]bool{}
	for key := range s.Allocations {
		p, _, _ := splitPair(key)
		seen[p] = true
	}
	for key := range s.Requests {
		p, _, _ := splitPair(key)
		seen[p] = true
	}
	return sortedKeys(seen)
}

func (s *System) Available() map[string]int {
	available := copyIntMap(s.Resources)
	for key, amount := range s.Allocations {
		_, r, _ := splitPair(key)
		available[r] -= amount
	}
	return available
}

func (s *System) WaitForGraph(available map[string]int) (map[string][]string, []Edge) {
	graph := map[string][]string{}
	edgeSet := map[string]bool{}
	var edges []Edge
	for _, p := range s.Processes() {
		graph[p] = nil
	}
	for reqKey, requested := range s.Requests {
		if requested <= 0 {
			continue
		}
		waitingProcess, resource, _ := splitPair(reqKey)
		if available[resource] >= requested {
			continue
		}
		for allocKey, held := range s.Allocations {
			holder, heldResource, _ := splitPair(allocKey)
			if holder == waitingProcess || heldResource != resource || held <= 0 {
				continue
			}
			key := waitingProcess + "->" + holder
			if edgeSet[key] {
				continue
			}
			edgeSet[key] = true
			graph[waitingProcess] = append(graph[waitingProcess], holder)
			edges = append(edges, Edge{From: waitingProcess, To: holder})
		}
	}
	for p := range graph {
		sort.Strings(graph[p])
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].From == edges[j].From {
			return edges[i].To < edges[j].To
		}
		return edges[i].From < edges[j].From
	})
	return graph, edges
}

func (s *System) SafeSequence(available map[string]int) (bool, []string, []string) {
	work := copyIntMap(available)
	finished := map[string]bool{}
	processes := s.Processes()
	var sequence []string

	progress := true
	for progress {
		progress = false
		for _, p := range processes {
			if finished[p] || !s.requestsCanBeSatisfied(p, work) {
				continue
			}
			finished[p] = true
			sequence = append(sequence, p)
			for key, amount := range s.Allocations {
				holder, resource, _ := splitPair(key)
				if holder == p {
					work[resource] += amount
				}
			}
			progress = true
		}
	}

	var blocked []string
	for _, p := range processes {
		if !finished[p] {
			blocked = append(blocked, p)
		}
	}
	return len(blocked) == 0, sequence, blocked
}

func (s *System) requestsCanBeSatisfied(process string, available map[string]int) bool {
	for key, amount := range s.Requests {
		p, resource, _ := splitPair(key)
		if p == process && amount > available[resource] {
			return false
		}
	}
	return true
}

func Recommendations(cycles [][]string, blocked []string, graph map[string][]string) []string {
	if len(cycles) == 0 && len(blocked) == 0 {
		return []string{"No deadlock was found. Continue monitoring resource requests before granting new allocations."}
	}
	recs := []string{
		"Abort or roll back one process in a detected cycle to release its resources.",
		"Require all processes to request shared resources in the same global order.",
		"Use timeouts for long waits and retry after releasing partial allocations.",
		"Before granting a request, run a safe-state check similar to Banker's Algorithm.",
	}
	if len(cycles) > 0 {
		recs = append(recs, "Detected cycle: "+formatCycle(cycles[0])+". Breaking any edge in this cycle removes the immediate deadlock.")
	}
	if len(blocked) > 0 {
		recs = append(recs, "Blocked processes: "+join(blocked)+". Inspect their outgoing wait-for edges first.")
		_ = graph
	}
	return recs
}
