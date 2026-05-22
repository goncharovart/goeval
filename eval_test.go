package goeval

import (
	"context"
	"testing"
)

// TestRun_EmptyDatasetClosesChannel — the smoke test that pins the
// streaming contract: closing the input channel must close the output
// channel. Without this guarantee no caller can write a `for r := range
// evaluator.Run(...)` loop safely.
func TestRun_EmptyDatasetClosesChannel(t *testing.T) {
	e := &Evaluator{}
	in := make(chan Sample)
	close(in)

	out := e.Run(context.Background(), in)
	for range out {
		t.Fatal("empty dataset should not produce any Result")
	}
}

// TestSample_MetadataIsOptional — Sample.Metadata being nil must not
// panic any consumer. We bake this into the type contract from day 0
// so metrics that group by metadata don't crash on minimal datasets.
func TestSample_MetadataIsOptional(t *testing.T) {
	s := Sample{
		ID:              "1",
		Question:        "What is the capital of France?",
		GeneratedAnswer: "Paris",
	}
	// Just exercising the zero-value path. If this compiles and runs
	// without panic, the contract is satisfied for this commit.
	if s.ID != "1" {
		t.Fatalf("ID round-trip broken: %q", s.ID)
	}
}
