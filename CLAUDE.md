# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go test ./...           # run all tests
go test -run TestName   # run a single test
go test -v ./...        # verbose test output
go build ./...          # build all packages
```

To run an example:
```bash
go run ./examples/sudoku/
go run ./examples/queen/
```

## Architecture

This library implements **Donald Knuth's Algorithm 7.2.2.1M** (Dancing Links / Algorithm X) for solving **exact cover problems with multiplicities and colors (MCC)**. It also handles simpler XC and XCC variants as special cases.

### Core files

- **mcc.go** — `MCC` struct definition and constructor (`NewDancer`, `WithContext`). Holds the node/item arrays and configuration.
- **input.go** — Parses the constraint matrix from an `io.Reader` into the internal linked-list data structure.
- **dance.go** — The solver. `Dance()` runs the algorithm in a goroutine and streams solutions over a channel. Contains `cover`, `uncover`, `hide`, `unhide`, `purify`, `unpurify`, `tweak`, `untweak`.

### Data structures

Two flat arrays act as the Dancing Links structure:

- `nd []node` — each node has `up/down` circular links within its item column, an `itm` pointer back to its item, and optional `color`/`colorName` for XCC constraints.
- `cl []item` — each item has `prev/next` circular links for the active item list, `name`, `bound` (residual capacity), and `slack` (for multiplicity).

The boundary between primary and secondary items is tracked by `MCC.second`.

### Input format

Items are whitespace-separated; `|` separates primary items from secondary (optional) items. Options follow, one per line.

```
| comment line
A B C | D E         ← A B C primary; D E secondary
A C                 ← option 1
B D:red E:red       ← option 2 with color constraints
```

Multiplicity bounds use `lo:hi|itemname` notation (e.g., `2:3|X` means item X must appear 2–3 times).

### Usage pattern

```go
d := dlx.NewDancer(input)           // parse input
for sol := range d.Dance() {        // iterate solutions (channel)
    // sol is [][]string — each inner slice is one chosen option
}
```

`Dance()` respects a context (set via `dlx.WithContext(ctx)`) for cancellation and supports a `PulseInterval` for periodic heartbeat notifications during long solves.
