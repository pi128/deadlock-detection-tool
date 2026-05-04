# Deadlock Detection Tool

## Name and KSU-ID

Name:  
KSU-ID:

## Abstract

This project implements a Go-based deadlock detection simulator for operating systems. The program reads a resource-allocation scenario, computes available resources, builds a wait-for graph between processes, detects cycles that represent deadlock, and recommends recovery or prevention strategies. The project demonstrates that deadlock detection is practical for a known finite system state, while future deadlock prediction for arbitrary programs is theoretically limited and related to undecidability concepts such as the Halting Problem.

## Introduction

Modern operating systems run many processes concurrently. These processes often compete for limited resources such as files, printers, locks, memory regions, and devices. A deadlock occurs when a group of processes wait forever because each process is holding a resource needed by another process in the same group.

## Background Information

Deadlock is usually described through four necessary conditions: mutual exclusion, hold and wait, no preemption, and circular wait. If all four conditions hold at the same time, deadlock can occur. Operating systems use detection, prevention, avoidance, or recovery strategies to handle this problem.

Deadlock detection for a current state is solvable because the system has a finite number of processes, resources, allocations, and requests. The state can be represented as a graph and analyzed for cycles. However, predicting whether any arbitrary program will eventually deadlock in the future is not generally decidable. A complete predictor would need to reason about all possible program paths, inputs, schedules, and resource requests.

## Problem Definition

The problem is to determine whether a set of processes is currently deadlocked based on known resource ownership and pending requests. The tool must identify which processes are involved and suggest practical methods for breaking or preventing the deadlock.

## Proposed Solution

The solution models the system as a wait-for graph. Each process is a node. If process `P1` is waiting for a resource held by process `P2`, the graph contains an edge from `P1` to `P2`. A cycle in this graph means the processes in the cycle are waiting on each other and cannot proceed without outside intervention.

## Methods

The tool uses these steps:

1. Read a JSON scenario containing total resources, current allocations, and current requests.
2. Calculate currently available resources.
3. Build the wait-for graph.
4. Search the graph for cycles using depth-first search.
5. Run a safe-state check to find whether all processes can eventually finish.
6. Print detected cycles, safe sequence information, and recommendations.

## Implementation and Experimental Results

The project is implemented in Go. The main command is located in `cmd/deadlock`, and the detection logic is located in `internal/deadlock`.

Example deadlock scenario:

```json
{
  "resources": {"printer": 1, "scanner": 1, "database": 1},
  "allocations": {"P1:printer": 1, "P2:scanner": 1, "P3:database": 1},
  "requests": {"P1:scanner": 1, "P2:database": 1, "P3:printer": 1}
}
```

In this scenario, `P1` waits for `P2`, `P2` waits for `P3`, and `P3` waits for `P1`. The tool reports the cycle `P1 -> P2 -> P3 -> P1`, which identifies the deadlock.

## Discussion

The graph-based method is useful because it gives a clear explanation of why the system is stuck. Instead of only saying that a deadlock exists, the tool reports the exact cycle of process dependencies. This makes it easier for an administrator or operating system policy to choose which process to abort, roll back, or preempt.

## Limitations

The simulator analyzes a static snapshot. It does not execute real operating-system processes or intercept real kernel locks. It also does not solve the general future prediction problem. A program that is safe in one state may still deadlock later depending on future inputs, scheduling decisions, or resource requests.

## Conclusion

The project shows that deadlock detection can be solved for a known finite resource-allocation state using a wait-for graph and cycle detection. It also shows why the broader problem of predicting all possible future deadlocks is much harder and cannot be completely solved for arbitrary programs.

## References

1. Abraham Silberschatz, Peter Baer Galvin, and Greg Gagne, *Operating System Concepts*.
2. Andrew S. Tanenbaum and Herbert Bos, *Modern Operating Systems*.
3. Edsger W. Dijkstra, "Cooperating Sequential Processes."
4. Alan M. Turing, "On Computable Numbers, with an Application to the Entscheidungsproblem."
