// Package goeval is the entry point for evaluating Retrieval-Augmented
// Generation (RAG) pipelines from Go.
//
// The shape of an evaluation is intentionally simple:
//
//	dataset → Evaluator → channel of Result
//
// where the dataset is a stream of `Sample`s (each carrying a question, a
// retrieved context, and an expected answer), the Evaluator applies one
// or more Metrics to each sample, and the result channel can be consumed
// for aggregation, CI gating, or storage.
//
// This file declares the core public interfaces. Implementations live in
// the metric/ and judge/ sub-packages (lands in subsequent commits).
package goeval

import "context"

// Sample is one row of an evaluation dataset.
//
// Question is what was asked. RetrievedContext is what the RAG pipeline
// passed to the LLM (one or more chunks of text, in retrieval order).
// GeneratedAnswer is what the LLM produced. ReferenceAnswer is the gold
// answer, when available — many metrics (BLEU, ROUGE, context recall)
// only make sense with a reference.
//
// Metadata carries free-form per-sample tags (split, dataset version,
// tenant id, etc.) — surfaces back on Result so callers can group.
type Sample struct {
	ID               string
	Question         string
	RetrievedContext []string
	GeneratedAnswer  string
	ReferenceAnswer  string // optional
	Metadata         map[string]string
}

// Result is the score of one Metric on one Sample.
//
// Score is the bounded value the metric defines (most metrics use
// [0, 1]; some use binary 0/1 and a few use unbounded floats — each
// Metric documents its own range). Reasoning is the LLM-as-judge
// rationale when the metric is judge-backed; empty otherwise.
//
// Error carries a non-nil value when the metric failed to score this
// sample (judge unreachable, malformed sample, etc.). Callers should
// treat Error-bearing Results as "no signal" rather than "score 0".
type Result struct {
	SampleID string
	Metric   string  // name of the metric ("faithfulness", "context_recall", ...)
	Score    float64 // bounded per-metric; see Metric.Range() once implemented
	Reasoning string // optional LLM rationale
	Error    error
}

// Metric scores a single Sample. Implementations live under metric/.
//
// Name returns the stable identifier used in Result.Metric and in CLI
// output. It MUST be stable across versions; if a metric needs a
// semantic change, ship a new metric with a v2 suffix.
//
// Score is called once per Sample. Implementations should treat ctx
// cancellation as a fast-fail (return ctx.Err()). Implementations
// that talk to an LLM should accept a Judge via their constructor;
// goeval does not impose a default Judge in this interface, on
// purpose — that decision belongs to the implementation.
type Metric interface {
	Name() string
	Score(ctx context.Context, s Sample) (float64, string, error)
}

// Judge is the LLM-as-judge abstraction. Faithfulness, context relevance,
// answer correctness, and hallucination detection are all judge-backed
// metrics — they hand the question + answer + context to a strong LLM
// and parse the verdict.
//
// Ask returns the raw judge response; the metric is responsible for
// extracting a score from it (regex on "0.75", JSON parsing, etc.).
// We keep Ask raw rather than typed so each metric can craft its own
// prompt-and-response protocol without fighting a shared schema.
type Judge interface {
	Ask(ctx context.Context, prompt string) (string, error)
}

// Evaluator runs a stream of Samples through a fixed set of Metrics
// and emits Results on the returned channel.
//
// Run closes the result channel after the input channel closes AND
// every in-flight metric has reported. Callers should drain the
// channel; abandoning it leaks goroutines.
//
// Concurrency is bounded by Concurrency; values <= 0 default to
// runtime.GOMAXPROCS(0). Each (Sample, Metric) pair runs in its own
// goroutine up to the concurrency cap — large datasets with many
// metrics will need a backing Judge that can sustain that QPS.
type Evaluator struct {
	Metrics     []Metric
	Concurrency int
}

// Run is the streaming entry point. See Evaluator doc.
// Implementation lands in eval_runner.go in the next commit.
func (e *Evaluator) Run(ctx context.Context, in <-chan Sample) <-chan Result {
	// Placeholder; full implementation lands in eval_runner.go.
	// The signature is frozen so dependent code can be written
	// against the API while the streaming machinery materialises.
	out := make(chan Result)
	close(out)
	return out
}
