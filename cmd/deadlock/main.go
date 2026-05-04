package main

import (
	"flag"
	"fmt"
	"os"

	"deadlock-detection-tool/internal/deadlock"
)

func main() {
	file := flag.String("file", "examples/deadlock.json", "path to a JSON scenario file")
	flag.Parse()

	f, err := os.Open(*file)
	if err != nil {
		exitf("open scenario: %v", err)
	}
	defer f.Close()

	system, err := deadlock.Load(f)
	if err != nil {
		exitf("load scenario: %v", err)
	}

	result := system.Detect()
	printResult(result)
}

func printResult(r deadlock.DetectionResult) {
	fmt.Println("Deadlock Detection Tool")
	fmt.Println("=======================")
	fmt.Println()
	fmt.Println("Available resources:")
	for _, resource := range sortedIntKeys(r.Available) {
		fmt.Printf("  %s: %d\n", resource, r.Available[resource])
	}

	fmt.Println()
	fmt.Println("Wait-for graph:")
	if len(r.Edges) == 0 {
		fmt.Println("  no wait-for edges")
	} else {
		for _, edge := range r.Edges {
			fmt.Printf("  %s -> %s\n", edge.From, edge.To)
		}
	}

	fmt.Println()
	if r.Deadlocked {
		fmt.Println("Status: DEADLOCK DETECTED")
		for i, cycle := range r.Cycles {
			fmt.Printf("  Cycle %d: %s\n", i+1, formatCycle(cycle))
		}
	} else {
		fmt.Println("Status: no current deadlock detected")
	}

	if r.Safe {
		fmt.Println("Safe-state check: safe")
		fmt.Printf("Safe sequence: %s\n", join(r.SafeSequence))
	} else {
		fmt.Println("Safe-state check: unsafe")
		fmt.Printf("Unable to safely finish: %s\n", join(r.Blocked))
	}

	fmt.Println()
	fmt.Println("Recommended actions:")
	for _, rec := range r.Recommendations {
		fmt.Printf("  - %s\n", rec)
	}
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
