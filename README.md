# goeval

> RAGAS for Go. Faithfulness, context recall, answer relevance, hallucination — measure your RAG pipeline without leaving the Go stack.

[![Go Reference](https://pkg.go.dev/badge/github.com/goncharovart/goeval.svg)](https://pkg.go.dev/github.com/goncharovart/goeval)
[![Go Report Card](https://goreportcard.com/badge/github.com/goncharovart/goeval)](https://goreportcard.com/report/github.com/goncharovart/goeval)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

> ⚠️ **Pre-MVP.** API is exploration-grade until `v0.1.0`. Star/watch to follow development. First metric (faithfulness) lands in the next commit batch.

## Why goeval

Python owns the RAG-eval ecosystem: [RAGAS](https://github.com/explodinggradients/ragas) (8k+ stars), [DeepEval](https://github.com/confident-ai/deepeval) (3k+ stars), [TruLens](https://github.com/truera/trulens), Phoenix from Arize. All Python.

In Go — nothing comparable. Teams building RAG pipelines in Go (on `langchaingo`, `cloudwego/eino`, or hand-rolled OpenAI/Anthropic clients) currently have two options:

1. **Export to Python** — run RAGAS in a separate process, plumb the data back. Slow, breaks the single-binary deployment story.
2. **Roll your own metrics** — every team reinvents faithfulness, context-recall, hallucination detection.

Both are wasteful. **goeval** is the Go-native alternative.

## Design principles

- **Streaming-first.** Eval runs are pipelines: dataset → evaluator → metric. Channels are first-class. Eval a 10k-sample dataset without blocking your CI on memory.
- **LLM-as-judge done right.** Faithfulness, context relevance, answer correctness rely on a strong LLM. goeval supports a `Judge` abstraction so you can swap GPT-4 → Claude → local Llama transparently.
- **Deterministic metrics in addition.** Context recall (overlap-based), BLEU/ROUGE-style — no LLM dependency, fast and reproducible.
- **No framework lock-in.** Adapters for `langchaingo`, `eino`, raw OpenAI/Anthropic Go SDKs, but the core is dependency-light.
- **CI-friendly.** Exit code on regression, JSON/Markdown reports, GitHub Action template.

## Roadmap to v0.1.0

- [ ] `Evaluator`, `Metric`, `Dataset`, `Judge` interfaces
- [ ] Streaming pipeline (dataset → evaluator → channel of `Result`)
- [ ] Metric: **faithfulness** (LLM-judge against retrieved context)
- [ ] Metric: **context relevance** (LLM-judge)
- [ ] Metric: **answer relevance** (cosine on embeddings)
- [ ] Metric: **context recall** (deterministic overlap)
- [ ] Metric: **hallucination detection** (LLM-judge + groundedness check)
- [ ] Judge adapters: OpenAI, Anthropic, local Ollama
- [ ] Dataset adapter: JSON/JSONL ingestion
- [ ] Adapter: `langchaingo`
- [ ] Adapter: `cloudwego/eino`
- [ ] CLI: `goeval run dataset.jsonl --metric faithfulness,context_recall`
- [ ] GitHub Action template for CI regression gates
- [ ] CHANGELOG + CONTRIBUTING + CoC
- [ ] v0.1.0 tag + release

Estimated runway to v0.1.0: **6 weeks** (one solo maintainer, evenings).

## Install (when v0.1.0 ships)

```bash
go get github.com/goncharovart/goeval@v0.1.0
```

## Inspiration

- [explodinggradients/ragas](https://github.com/explodinggradients/ragas) — Python progenitor of RAG-eval
- [confident-ai/deepeval](https://github.com/confident-ai/deepeval) — 50+ metrics, LLM-as-judge
- [truera/trulens](https://github.com/truera/trulens) — observability + eval

## Status & maintenance

Pre-MVP solo development. Issues + PRs welcome; expect slow review (evenings only) until v0.1.0.

Build openly — every architectural decision goes into [docs/design.md](docs/design.md) once it stabilises.

## License

MIT — see [LICENSE](LICENSE).
