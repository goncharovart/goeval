# Changelog

All notable changes to goeval are tracked here. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/). The project
is **pre-v0.1.0** — the API is exploration-grade and may change
between commits without a SemVer bump.

## [Unreleased]

### Added (initial scaffold, 2026-05-22)

- Package scaffold: `eval.go` declares the core public interfaces
  (`Sample`, `Result`, `Metric`, `Judge`, `Evaluator`).
- Streaming contract pinned by `TestRun_EmptyDatasetClosesChannel`.
- MIT license, README with design principles + roadmap to v0.1.0,
  .gitignore covering Go + IDE noise + `eval_outputs/` to keep
  private evaluation data out of origin.

### Roadmap to v0.1.0

See [README.md](README.md) for the full roadmap. Six weeks to MVP
under solo evening development. First metric to land:
**faithfulness** (LLM-judge against retrieved context).
