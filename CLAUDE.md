# CLAUDE.md — goeval

> Project context for AI coding agents working on `goeval`. Keep this file in sync with reality — it is the **single source of truth** for project conventions.

## What this is

`goeval` — **RAGAS for Go**. A streaming, type-safe LLM-as-judge evaluation pipeline that lets Go teams measure RAG / agent / extraction quality without hand-rolling boilerplate around OpenAI-style chat APIs.

Niche: existing RAGAS implementations are Python-only (ragas, deepeval, promptfoo). Go teams shipping LLM features in production currently dump JSONL into Python notebooks. `goeval` keeps the eval pipeline in-process.

## Stack

- **Go 1.22+** (channels, generics, `slices`, `cmp`)
- **OpenAI-compatible chat completions API** (works with OpenAI / Azure OpenAI / vLLM / Ollama via `OPENAI_BASE_URL`)
- Standard library only for the core (`net/http`, `encoding/json`, `context`); zero non-stdlib runtime dependencies in `eval.go`
- `testing` + table-driven tests with subtests

## Project layout

```
goeval/
├── eval.go              # Core: Evaluator interface, streaming pipeline
├── eval_test.go         # Table-driven tests
├── README.md            # User-facing intro, install, 1-screen example
├── CHANGELOG.md         # Keep-a-changelog format
├── LICENSE              # MIT
├── go.mod               # module github.com/goncharovart/goeval
├── examples/            # Runnable examples (added incrementally)
└── docs/                # Design notes, comparison-with-ragas
```

When adding sub-packages — **prefer `internal/` for non-public helpers**, `pkg/` ONLY for explicitly exported reusable code outside the main API surface. Avoid `pkg/` cargo-culting if a flat layout suffices.

## Build & test

```bash
# Run all tests with race detector + cover
go test -race -count=1 -cover ./...

# Single test
go test -race -run TestStreamingEvaluator ./...

# Vet + lint (golangci-lint v2 config will land in next iteration)
go vet ./...

# Format
gofmt -s -w .
goimports -w .  # if installed

# Tidy modules
go mod tidy
```

## Coding conventions

### Idiomatic Go

- **Accept interfaces, return structs.** `Evaluator` interface for swappable backends; concrete types for `OpenAIEvaluator`, `LocalEvaluator`.
- **Channels for orchestration, mutexes for state.** Streaming pipeline uses unbuffered `<-chan Result`; per-evaluator metric counters live behind `sync.RWMutex`.
- **Errors wrap context.** `fmt.Errorf("evaluate %s: %w", testID, err)` — never bare `return err` without context at API boundaries.
- **Context propagation everywhere.** Every public method accepts `context.Context` as first parameter. Cancellation honored within ≤100ms.

### Modern Go (Go 1.22+ features to prefer)

| Old pattern | Prefer |
|---|---|
| `atomic.AddInt64`, `atomic.LoadInt64` | `atomic.Int64` value type |
| `for i := 0; i < N; i++` with no body deps on `i` | `for range N` |
| `if val, ok := m[key]; ok` then re-lookup | Local copy first |
| `strings.Replace(s, old, new, 1)` | `strings.Replace` with explicit `-1` for all, otherwise `Cut` |
| Manual loop to sort by field | `slices.SortFunc(s, cmp.Compare)` |

### Testing

- **Table-driven tests with subtests** named after the scenario, not the test number: `t.Run("empty input returns zero metric", ...)`.
- **`-race` clean.** All tests pass `-race -count=10` before commit.
- **Coverage target: 80%+** on `eval.go` core. Coverage measured by `go test -cover ./...`; CI fails below 75%.
- **No external network in unit tests.** Use a test HTTP server (`httptest.NewServer`) to mock OpenAI responses.

## Pre-commit hook (recommended)

```bash
#!/usr/bin/env bash
# .githooks/pre-commit
set -e
gofmt -s -w .
go vet ./...
go test -race -short ./...
golangci-lint run --new-from-rev=HEAD~1 || true  # warn-only on changed lines
```

Enable: `git config core.hooksPath .githooks`

## Roadmap to v0.1.0

- [x] Streaming `Evaluator` interface + OpenAI-compatible adapter
- [ ] Built-in metrics: faithfulness, answer-relevance, context-precision (RAGAS canonical four)
- [ ] Concurrent batch eval with bounded worker pool
- [ ] OpenTelemetry tracing for per-test spans
- [ ] CLI: `goeval run testset.jsonl --out report.json`
- [ ] Example: evaluate a RAG over local markdown corpus
- [ ] CI: GitHub Actions matrix (1.22, 1.23) + Codecov badge

See [README.md](README.md) for current state. v0.1.0 ETA: ~6 weeks evening development.

## When in doubt

- Idiomatic Go > clever Go. Prefer clarity, accept verbosity.
- Public API is unstable until v1.0.0 — breaking changes allowed in v0.x but **call them out in CHANGELOG**.
- Don't introduce a dependency just to save 20 lines. Stdlib first.
- Tests describe behavior, not implementation. If a test breaks because internals refactored without behavior change, the test was wrong.
