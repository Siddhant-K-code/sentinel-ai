# AGENT

## Project conventions

- **Logging**: Use structured logs with zap or logrus. Include context in error messages.
- **Errors**: Wrap errors with context using `fmt.Errorf` or `errors.Wrap`. Avoid panics in production code.
- **Security**: Never render raw user input. Always escape HTML, SQL, and other contexts.
- **Testing**: Write unit tests for all public functions. Use table-driven tests where appropriate.

## Scope

- **Targets**: Go, Rust, TypeScript/JavaScript packages under `./pkg`, `./cmd`, and `./internal`.
- **Public API**: Packages under `./pkg/api` are considered public and require careful review.
- **Private packages**: Packages under `./internal` are private to this module.

## Style

- **Go formatting**: Use `gofmt` and `goimports`. Follow standard Go conventions.
- **Commits**: Use Conventional Commits format: `type(scope): description`.
- **Code review**: All changes require at least one approval from a maintainer.

## Risk posture

- **Breaking changes**: Do not change exported symbols in `./pkg/api` without a deprecation notice.
- **Security**: Treat all user input as potentially malicious. Validate and sanitize all inputs.
- **Performance**: Avoid blocking operations in hot paths. Use timeouts for external calls.

## Dependencies

- **External tools**: Prefer tools that are available in the CI environment.
- **Version pinning**: Pin all dependency versions for reproducible builds.
- **Security updates**: Regularly update dependencies and scan for vulnerabilities.
