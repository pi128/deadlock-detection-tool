# Deadlock Detection Tool

This Go project simulates an operating-system resource allocation state, builds a wait-for graph, detects deadlock cycles, and reports practical recovery or prevention actions.

The project also supports the theoretical discussion for the report: a current finite resource-allocation state can be checked for deadlock, but no general algorithm can perfectly predict whether every arbitrary program will eventually deadlock before it runs. That future-prediction problem is related to the same undecidability ideas behind the Halting Problem.

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
- Recovery and prevention recommendations

## Suggested Report Thesis

Operating systems can detect deadlocks in a known resource-allocation state by modeling process waits as a graph and searching for cycles. However, predicting all future deadlocks for arbitrary programs is not completely solvable in general, because doing so would require reasoning about every possible future execution path of a program.
