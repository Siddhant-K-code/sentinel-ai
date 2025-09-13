# AGENT

## Project conventions

- Logging: use structured logs (zap/logrus).
- Errors: wrap with context; avoid panics.
- Security: never render raw user input; escape HTML.

## Scope

- Targets: Go, Rust, TS packages under ./pkg and ./cmd.
- Public API: packages under ./pkg/api are public.

## Style

- Go formatting: `gofmt` + `goimports`.
- Commits: Conventional Commits.

## Risk posture

- Do not change exported symbols in ./pkg/api without a deprecation note.
