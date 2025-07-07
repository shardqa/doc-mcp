# Testing and Coverage Practices for MCP Server

## Main.go and Coverage
- The main.go file should be kept as a minimal entrypoint, delegating all logic to importable packages (such as internal/). This is the idiomatic Go approach and ensures maintainability and testability.
- Go's standard testing and coverage tools do not support direct testing or coverage of main.go's main() function, as it is not importable or callable from tests.
- All business logic should be moved out of main.go into packages. Only these packages should be tested and covered.

## Integration and End-to-End Testing
- End-to-end (CLI) tests can be written using exec.Command to spawn the binary and test its behavior, but these do not contribute to Go's coverage statistics for main.go.
- As of Go 1.20, it is possible to build coverage-instrumented binaries and collect coverage data from integration tests using the GOCOVERDIR environment variable. However, this still does not cover main.go itself, only the packages it uses.

## Best Practices for Testability and Coverage
- Refactor all logic out of main.go into importable packages (e.g., internal/server.go).
- Write in-process unit tests for these packages to maximize coverage.
- Use in-package tests (e.g., internal/server_test.go) to ensure coverage is tracked.
- Use CLI/integration tests for end-to-end safety, but do not expect these to increase coverage for main.go.

## Project-Specific Implementation
- The MCP server's main.go now only calls internal.RunServer with os.Stdin and os.Stdout.
- All server logic is in internal/server.go, which is fully testable and coverable.
- In-process unit tests for internal/server.go are tracked by Go's coverage tool.
- End-to-end tests exist for CLI safety, but do not affect coverage statistics for main.go.

## References
- Go Blog: The cover story (https://go.dev/blog/cover)
- Go Blog: Code coverage for Go integration tests (https://go.dev/blog/integration-test-coverage)

## Summary
- There is no supported or idiomatic way to directly test or cover main.go's main() function in Go.
- All logic should be in packages, and those packages should be tested and covered.
- This project follows Go best practices for testability and coverage. 