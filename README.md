# Deadlock Detection Tool

This is a small Go command-line project for checking a fixed operating-system resource-allocation snapshot. It builds a wait-for graph, detects deadlock cycles, runs a simplified safe-state check, and prints practical recovery or prevention actions.

The report uses the tool to separate two ideas: current-state deadlock detection is finite graph analysis, while perfect future deadlock prediction for arbitrary programs is not generally solvable.

## Run

```sh
go run ./cmd/deadlock -file examples/deadlock.json
```

The deadlock example prints the detected cycle. A safe example is also included:

```sh
go run ./cmd/deadlock -file examples/safe.json
```

## Test

```sh
go test ./...
```

## Input Format

Scenarios are JSON files with three maps:

```json
{
  "resources": {
    "printer": 1,
    "scanner": 1
  },
  "allocations": {
    "P1:printer": 1,
    "P2:scanner": 1
  },
  "requests": {
    "P1:scanner": 1,
    "P2:printer": 1
  }
}
```

- `resources` stores the total number of each resource instance.
- `allocations` stores resources currently held by each process.
- `requests` stores resources each process is waiting to acquire.
- Allocation and request keys use `process:resource` format.

## What It Shows

- Available resource calculation
- Wait-for graph construction
- Cycle detection in the wait-for graph
- Safe-state checking similar to the idea behind Banker's Algorithm
- Input validation for impossible resource snapshots
- Recovery and prevention recommendations

## Design Notes

- The detector lives in `internal/deadlock`; `cmd/deadlock` only handles command-line input and printing.
- The JSON format uses flat `process:resource` keys so the examples stay short.
- The program sorts graph output so screenshots and tests stay stable between runs.
