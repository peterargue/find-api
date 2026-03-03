package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// debug is true when LIVE_TEST_DEBUG=1 is set. Suites use dumpJSON to print
// the raw response when a field they need turns out to be empty.
var debug = os.Getenv("LIVE_TEST_DEBUG") == "1"

// dumpJSON prints v as indented JSON to stdout when debug mode is on.
func dumpJSON(label string, v any) {
	if !debug {
		return
	}
	b, _ := json.MarshalIndent(v, "    ", "  ")
	fmt.Printf("  [DEBUG] %s: %s\n", label, string(b))
}

// Test is a single named check with a run function.
// Run returns a short summary string on success, or an error on failure.
type Test struct {
	Name string
	Run  func(ctx context.Context) (summary string, err error)
}

// Suite is a named group of related tests sharing chained state.
type Suite struct {
	Name  string
	Tests []Test
}

// Runner executes suites sequentially and prints verbose output.
type Runner struct {
	suites []Suite
	pass   int
	fail   int
}

func NewRunner(suites []Suite) *Runner {
	return &Runner{suites: suites}
}

func (r *Runner) Run(ctx context.Context) {
	for _, s := range r.suites {
		fmt.Printf("\n=== %s ===\n", s.Name)
		for _, t := range s.Tests {
			summary, err := t.Run(ctx)
			if err != nil {
				fmt.Printf("  ✗ %-40s %v\n", t.Name, err)
				r.fail++
			} else {
				fmt.Printf("  ✓ %-40s (%s)\n", t.Name, summary)
				r.pass++
			}
		}
	}
}

func (r *Runner) PrintSummary() {
	total := r.pass + r.fail
	fmt.Printf("\nPASS %d/%d  FAIL %d/%d\n", r.pass, total, r.fail, total)
}

func (r *Runner) ExitCode() int {
	if r.fail > 0 {
		return 1
	}
	return 0
}

// require returns an error if v is empty; dependent tests use this to fail
// clearly when a prerequisite test did not produce a usable value.
func require(name, v string) error {
	if v == "" {
		return fmt.Errorf("prerequisite missing: no %s from previous test", name)
	}
	return nil
}

// requireUint64 returns an error if v is zero.
func requireUint64(name string, v uint64) error {
	if v == 0 {
		return fmt.Errorf("prerequisite missing: no %s from previous test", name)
	}
	return nil
}
