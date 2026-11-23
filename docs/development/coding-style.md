# Coding Style and Development Guidelines

This document defines coding conventions and development practices for **ProxyNav V2bX**.

The goal is to keep the codebase readable, maintainable, and consistent while allowing the project to support multiple proxy cores, protocols, panel integrations, and optional features.

These guidelines are based on the existing project structure and dependencies. Where the existing implementation has legacy conventions, contributors should avoid introducing unnecessary inconsistencies and should prefer incremental improvements.

## 1. General Principles

When contributing code, follow these principles:

1. **Consistency:** Follow the conventions already established in the surrounding package.
2. **Simplicity:** Prefer straightforward implementations over unnecessary abstractions.
3. **Compatibility:** Consider existing configurations, panel integrations, and deployed nodes before changing behavior.
4. **Modularity:** Keep core-specific and protocol-specific logic within the appropriate components.
5. **Reliability:** Handle errors explicitly and manage resources carefully.
6. **Testability:** Make important logic independently testable whenever practical.
7. **Security:** Never expose credentials, private keys, or sensitive user information through logs or error messages.

Before introducing a new abstraction, inspect the existing implementation to determine whether a suitable interface or helper already exists.

## 2. Go Formatting and Language Conventions

### 2.1 Formatting

All Go source files must be formatted with `gofmt`.

Before submitting changes, run:

```bash
gofmt -w .
```

To check formatting without modifying files:

```bash
gofmt -l .
```

Do not submit code that requires formatting corrections.

### 2.2 Naming

Follow standard Go naming conventions:

- Use `MixedCaps` and `mixedCaps` rather than underscores in Go identifiers.
- Use `PascalCase` for exported identifiers.
- Use `camelCase` for unexported identifiers.
- Keep package names short, lowercase, and descriptive.
- Avoid unnecessary abbreviations.
- Use initialisms consistently, such as `ID`, `IP`, `HTTP`, `TLS`, and `URL`.

Examples:

```go
type NodeConfig struct {
    NodeID   int
    ServerIP string
}

func LoadConfig(path string) error {
    // ...
    return nil
}

func validateNodeConfig(cfg NodeConfig) error {
    // ...
    return nil
}
```

Avoid generic names such as `data`, `info`, or `manager` when a more specific name improves clarity.

### 2.3 Comments and Documentation

Exported types, functions, constants, and variables should have documentation comments.

Comments should explain:

- Why a particular implementation exists.
- Important assumptions or constraints.
- Non-obvious behavior.
- Compatibility requirements.
- Concurrency and lifecycle guarantees.

Avoid comments that merely repeat the code.

Use English for new source-code comments and developer-facing documentation to maintain consistency across the project.

## 3. Package Organization

V2bX uses a modular structure. Keep new code in the package that owns the corresponding responsibility.

| Package or directory | Responsibility                                           |
| -------------------- | -------------------------------------------------------- |
| `cmd/`               | Application commands, startup, and command-line behavior |
| `conf/`              | Configuration structures, loading, and validation        |
| `api/panel/`         | Panel communication and API integration                  |
| `core/`              | Proxy core integration and core-related functionality    |
| `node/`              | Node configuration and node lifecycle                    |
| `limiter/`           | Connection and traffic limiting functionality            |
| `common/`            | Shared utilities and common functionality                |
| `example/`           | Example configuration and usage material                 |
| `test_data/`         | Test fixtures and supporting test data                   |

### 3.1 Keep Responsibilities Separate

Avoid placing unrelated functionality into a single package.

For example:

- Panel HTTP request construction belongs in the panel integration layer.
- Configuration parsing and configuration structures belong in `conf/`.
- Core-specific behavior belongs in `core/` or the relevant core implementation.
- Node-specific configuration and lifecycle behavior belongs in `node/`.
- Shared, genuinely reusable utilities may belong in `common/`.

Do not move code into `common/` merely because multiple packages use it. Shared code should have a clear and stable responsibility.

### 3.2 Avoid Circular Dependencies

Go packages must not import each other cyclically.

When two components need to communicate, consider:

- Defining a small interface at the appropriate abstraction boundary.
- Passing dependencies explicitly.
- Moving genuinely shared types into a suitable lower-level package.

Avoid introducing large interfaces that expose implementation details unnecessarily.

### 3.3 Keep Public APIs Small

Export only identifiers that need to be accessed by other packages.

Prefer unexported helpers for internal implementation details. A smaller public API makes future refactoring easier.

## 4. Configuration Design

Configuration is central to V2bX because it controls panel integration, proxy cores, nodes, and optional features.

### 4.1 Configuration Structures

When adding configuration fields:

- Use clear and descriptive field names.
- Follow the existing serialization and naming conventions.
- Add appropriate tags when required by the configuration format.
- Use types that represent the expected values accurately.
- Consider zero values and backward compatibility.
- Document fields whose behavior is not self-explanatory.

Example:

```go
type ExampleConfig struct {
    Enabled bool   `json:"enabled"`
    Address string `json:"address"`
    Timeout int    `json:"timeout"`
}
```

This is an illustrative example. Follow the actual configuration structures and serialization conventions in `conf/` when implementing changes.

### 4.2 Validate Configuration Early

Validate required fields and incompatible combinations as close as practical to configuration loading or component initialization.

Prefer explicit errors:

```go
if cfg.Address == "" {
    return errors.New("address is required")
}
```

Avoid silently replacing invalid configuration with unexpected defaults.

For optional settings, document the default behavior and ensure it is consistent across relevant components.

### 4.3 Preserve Existing Configurations

Before changing a configuration field:

1. Search for all references to the field.
2. Check example configurations.
3. Check how the field is serialized and loaded.
4. Consider existing deployed configurations.
5. Provide a migration or compatibility strategy if necessary.

Do not rename or remove configuration fields without considering existing users.

### 4.4 JSON5 Compatibility

The project uses JSON5-compatible configuration processing.

When modifying configuration handling:

- Preserve supported JSON5 syntax.
- Avoid assuming that a standard JSON parser accepts every existing configuration file.
- Test comments, strings, nested objects, and arrays where relevant.
- Verify both parsing and semantic validation.

## 5. Logging and Error Handling

The project includes both Logrus and Zap dependencies. When modifying logging code, follow the logging API already used by the surrounding component rather than introducing another logging abstraction without a clear reason.

### 5.1 Logging Guidelines

Use log levels consistently:

| Level | Intended use                                  |
| ----- | --------------------------------------------- |
| Debug | Detailed diagnostic information               |
| Info  | Normal application lifecycle events           |
| Warn  | Recoverable problems or unexpected conditions |
| Error | Operations that failed and require attention  |

Examples of useful log context include:

- Node identifiers.
- Core names.
- Protocol types.
- Operation names.
- Relevant endpoint information.

Avoid logging:

- Passwords.
- Panel authentication tokens.
- Private keys.
- Complete user credentials.
- Sensitive request headers.
- Unredacted configuration containing secrets.

Use structured fields when the existing logging API supports them.

### 5.2 Error Handling

Handle errors explicitly.

Do not ignore errors unless the reason is clear and documented.

When returning an error from a lower-level operation, add context where it helps identify the failing operation:

```go
if err != nil {
    return fmt.Errorf("initialize node: %w", err)
}
```

Use `%w` when callers may need to inspect or unwrap the original error.

Avoid logging an error at every layer while also returning it unchanged. This can produce duplicate log entries. Prefer logging at the layer that has enough context to handle or report the failure.

### 5.3 Avoid Panic for Normal Failures

Do not use `panic` for expected runtime failures such as:

- Invalid configuration.
- Panel connection failures.
- Missing certificate files.
- Core initialization errors.
- Port conflicts.

Return errors and let the appropriate application layer decide how to handle them.

Reserve panic for genuinely unrecoverable programming errors or conditions explicitly justified by the implementation.

## 6. Panel API Development

Panel integration is an important boundary between V2bX and the management panel.

When modifying `api/panel/`:

1. Keep request construction and response handling consistent with existing methods.
2. Handle network errors, HTTP status codes, and response decoding explicitly.
3. Apply appropriate timeouts.
4. Validate response data before passing it to other components.
5. Avoid coupling panel response structures directly to unrelated internal implementations.
6. Do not log authentication credentials or sensitive response fields.

When changing an API method, consider:

- Whether existing panel versions remain compatible.
- Whether response fields can be absent or null.
- Whether the API may return unexpected values.
- Whether retries are safe for the operation.
- Whether an API failure should affect one node or the entire application.

Do not retry non-idempotent operations automatically without considering duplicate side effects.

## 7. Proxy Core and Protocol Development

V2bX supports multiple proxy cores and uses conditional compilation to include selected implementations.

Core and protocol changes should preserve clear boundaries between shared behavior and implementation-specific behavior.

### 7.1 Keep Core-Specific Logic Isolated

When adding or modifying a core integration:

- Keep core-specific configuration generation in the relevant core implementation.
- Avoid leaking core-specific types into unrelated packages.
- Reuse shared node and configuration abstractions where appropriate.
- Keep initialization and shutdown behavior explicit.
- Document core-specific limitations and dependencies.

Do not assume that two cores use identical configuration formats or lifecycle semantics.

### 7.2 Conditional Compilation and Build Tags

The project uses build tags to select proxy cores and optional features.

Examples include:

- `sing`
- `hysteria2`
- `with_quic`
- `with_grpc`
- `with_utls`
- `with_wireguard`
- `with_acme`
- `with_gvisor`

When adding code guarded by build tags:

1. Follow the existing build-tag naming and file organization.
2. Ensure the relevant implementation is excluded when its feature is disabled.
3. Check whether shared code references symbols that may not exist in other build variants.
4. Test the relevant combinations of enabled and disabled features.
5. Update build documentation when introducing a new tag.

Do not assume that a successful build with all tags enabled guarantees that reduced builds also compile.

### 7.3 Core Lifecycle

Core lifecycle code should make resource ownership and failure behavior clear.

Consider:

- Initialization order.
- Listener creation and port conflicts.
- Startup failures.
- Runtime configuration updates.
- Graceful shutdown.
- Cleanup of sockets, files, and goroutines.

Avoid starting background goroutines without a clear termination mechanism.

### 7.4 Protocol Compatibility

When implementing protocol-related changes, verify the behavior across the affected core integrations.

Consider:

- Client compatibility.
- TLS and certificate requirements.
- Transport settings.
- Authentication behavior.
- UDP support.
- Optional protocol features.
- Configuration translation between the panel and the selected core.

Do not assume that a field supported by one core has an equivalent implementation in another.

## 8. Concurrency and Resource Management

V2bX manages long-running node and core processes. Concurrency and resource cleanup must be handled carefully.

### 8.1 Goroutines

Every long-running goroutine should have a clear lifecycle.

Use context cancellation, channels, or another established synchronization mechanism to stop background work.

Avoid:

- Unbounded goroutine creation.
- Goroutines that cannot be stopped.
- Shared mutable state without synchronization.
- Blocking operations without appropriate cancellation or timeout behavior.

### 8.2 Shared State

When multiple goroutines access shared state:

- Use mutexes or other appropriate synchronization primitives.
- Keep critical sections small.
- Avoid holding locks during slow network or disk operations.
- Document non-obvious locking requirements.

Run race-enabled tests when modifying concurrent code:

```bash
go test -race ./...
```

### 8.3 Network and File Resources

Always consider cleanup for resources such as:

- HTTP response bodies.
- Network connections.
- Listening sockets.
- Files and temporary directories.
- Timers and tickers.
- Background goroutines.

Use `defer` for cleanup when its execution scope is appropriate.

For long-running services, ensure cleanup occurs when a node or core is stopped or reconfigured, not only when the entire application exits.

## 9. Dependency Management

The project has a substantial dependency graph, including proxy core implementations and networking libraries.

Before adding a dependency:

1. Check whether an existing dependency already provides the required functionality.
2. Evaluate maintenance activity and compatibility.
3. Consider transitive dependencies and binary size.
4. Verify licensing compatibility.
5. Check whether the dependency affects conditional builds.
6. Document important reasons for introducing it.

Use Go modules to manage dependencies.

```bash
go mod download
go mod verify
```

After intentionally changing dependencies, review both files:

```bash
git diff -- go.mod go.sum
```

Do not remove or upgrade indirect dependencies manually without understanding their relationship to the selected core implementations.

## 10. Testing Requirements

New functionality should include appropriate tests whenever practical.

### 10.1 Unit Tests

Use standard Go test files:

```text
example.go
example_test.go
```

Test important behavior such as:

- Configuration parsing and validation.
- Data conversion.
- Error handling.
- Protocol or core configuration generation.
- Node-related business logic.
- Limiter behavior.

### 10.2 Test Naming

Use descriptive test names:

```go
func TestLoadConfig(t *testing.T) {
    // ...
}

func TestLoadConfig_InvalidInput(t *testing.T) {
    // ...
}
```

Use table-driven tests when several inputs exercise the same behavior.

### 10.3 Test Isolation

Tests should avoid depending on:

- Production panel services.
- Public proxy endpoints.
- Real user credentials.
- External services that are not controlled by the test.

Use test fixtures, mock implementations, or controlled local services when appropriate.

Avoid tests that modify system firewall rules, network interfaces, or production configuration.

### 10.4 Verify Build Variants

When a change touches conditional compilation, test the relevant build variants.

At minimum, verify the default supported development build and the specific core or optional feature combinations affected by the change.

## 11. Documentation Requirements

Update documentation whenever a change affects how developers or operators use V2bX.

Documentation updates may be required for:

- New configuration fields.
- Changed defaults.
- New build tags.
- New proxy core integrations.
- Protocol compatibility changes.
- Installation or upgrade behavior.
- Breaking changes.
- New environment variables or external dependencies.

Keep example configurations synchronized with the implementation.

Do not document a feature as supported until its implementation and relevant build variants have been verified.

## 12. Git and Pull Request Guidelines

### 12.1 Commit Messages

Use concise, descriptive commit messages.

Recommended prefixes:

| Prefix      | Purpose                         |
| ----------- | ------------------------------- |
| `feat:`     | New functionality               |
| `fix:`      | Bug fix                         |
| `refactor:` | Code restructuring              |
| `docs:`     | Documentation changes           |
| `test:`     | Test additions or modifications |
| `build:`    | Build and dependency changes    |
| `ci:`       | CI workflow changes             |
| `chore:`    | Maintenance tasks               |

Examples:

```text
feat: add support for a new protocol option
fix: handle missing node configuration
refactor: simplify core initialization
docs: update development guidelines
test: cover configuration validation
```

These prefixes are recommended conventions, not a claim that every historical commit follows this format.

### 12.2 Pull Request Scope

Keep each pull request focused on a specific change.

A pull request should explain:

- What changed.
- Why the change is needed.
- Which packages or components are affected.
- Whether configuration compatibility is affected.
- Which build tags or cores are relevant.
- How the change was tested.

Include logs or screenshots only when they help explain the behavior. Remove sensitive information first.

### 12.3 Pull Request Checklist

Before submitting a pull request, verify:

- [ ] Go source files are formatted with `gofmt`.
- [ ] `go vet ./...` completes successfully.
- [ ] Relevant tests pass.
- [ ] Relevant build-tag combinations compile.
- [ ] Configuration compatibility has been considered.
- [ ] Errors are handled explicitly.
- [ ] Sensitive information is not included in logs or test data.
- [ ] Documentation and examples are updated where necessary.
- [ ] Changes to `go.mod` and `go.sum` are intentional.

## 13. Repository and Module Path Considerations

The official repository is:

<https://github.com/proxynav/V2bX>

The project is based on:

<https://github.com/wyx2685/V2bX>

The current source tree may retain legacy module paths and import paths inherited from its upstream history.

Until a dedicated module-path migration is completed:

- Do not independently rewrite import paths in unrelated changes.
- Follow the module path declared in the current `go.mod`.
- Keep linker flags and build scripts consistent with the actual package paths.
- Use the ProxyNav repository URL for new documentation, issue references, and contribution links.

A module-path migration should be handled as a separate, coordinated change because it can affect imports, build scripts, dependency references, and release automation.

## 14. Final Guideline

The primary objective is not to impose a new architecture on the existing codebase. It is to improve the consistency and reliability of ProxyNav V2bX while preserving its existing behavior and multi-core design.

When in doubt:

1. Inspect the existing implementation.
2. Follow the conventions of the affected package.
3. Keep changes small and focused.
4. Preserve compatibility unless a breaking change is intentional.
5. Add tests for the behavior being changed.
6. Document important design decisions.

For questions, bug reports, and contributions, use the official repository:

<https://github.com/proxynav/V2bX>
