# Testing Guide

This document defines the testing strategy and validation procedures for **ProxyNav V2bX**.

The goal is to maintain compatibility across supported proxy cores, protocols, transports, panel integrations, and build configurations while reducing regressions during development and releases.

This guide distinguishes between:

- Existing project build configuration and Go testing capabilities.
- Recommended testing practices for new and modified code.
- Runtime and integration checks that require additional infrastructure.

A successful compilation is necessary, but it does not prove that a proxy node works correctly at runtime.

## 1. Testing Environment

### 1.1 Required toolchain

Use the Go version declared by the project:

```text
Go 1.25
Toolchain go1.25.0
```

The current project build configuration also uses:

```bash
export GOEXPERIMENT=jsonv2
```

Use the same toolchain and environment when reproducing build or test results.

Check the environment:

```bash
go version
go env GOTOOLCHAIN GOOS GOARCH CGO_ENABLED
```

### 1.2 Initialize the repository

Clone the repository with submodules:

```bash
git clone --recurse-submodules \
  https://github.com/proxynav/V2bX.git

cd V2bX
```

If the repository was cloned without submodules:

```bash
git submodule update --init --recursive
```

Verify the checked-out submodule revisions:

```bash
git submodule status --recursive
```

Download dependencies:

```bash
GOEXPERIMENT=jsonv2 go mod download
```

Verify dependency integrity:

```bash
go mod verify
```

### 1.3 Keep test environments isolated

Tests should not depend on:

- Production panel accounts.
- Production node credentials.
- Real customer traffic.
- Production TLS private keys.
- Uncontrolled external services.
- Machine-specific configuration files.
- Existing production databases.

Use temporary directories, test configurations, local services, and dedicated test credentials.

Never commit secrets, private keys, or production configuration into test fixtures.

## 2. Testing Levels

V2bX contains multiple layers, including CLI and configuration handling, panel integration, proxy core management, node implementations, and shared utilities.

Testing should be performed at different levels.

| Level             | Purpose                                            | Typical execution                 |
| ----------------- | -------------------------------------------------- | --------------------------------- |
| Unit tests        | Verify individual functions and components         | `go test`                         |
| Package tests     | Validate behavior within a Go package              | `go test ./path/to/package`       |
| Integration tests | Validate interactions between components           | Dedicated integration environment |
| Build tests       | Verify core and feature build configurations       | `go build`                        |
| Runtime tests     | Verify actual node and protocol behavior           | Dedicated test nodes              |
| Release tests     | Validate release artifacts and deployment behavior | Docker and deployment checks      |

Do not treat these levels as interchangeable.

For example, a successful unit test does not establish that a Trojan, Hysteria2, or AnyTLS node can establish a working client connection.

## 3. Running Go Tests

### 3.1 Run all package tests

Run the complete Go test suite:

```bash
GOEXPERIMENT=jsonv2 go test ./...
```

For verbose output:

```bash
GOEXPERIMENT=jsonv2 go test -v ./...
```

### 3.2 Run tests for a specific package

When developing a particular component, start with the affected package:

```bash
GOEXPERIMENT=jsonv2 go test ./common/...
```

Replace the package path with the actual package being modified.

This provides faster feedback than running the entire test suite after every small change.

### 3.3 Run a specific test

Run a test by name:

```bash
GOEXPERIMENT=jsonv2 go test \
  -run '^TestExample$' \
  ./path/to/package
```

Run a specific subtest:

```bash
GOEXPERIMENT=jsonv2 go test \
  -run 'TestExample/CaseName' \
  ./path/to/package
```

Replace the example names with actual test functions.

### 3.4 Disable the test cache

For troubleshooting or validation that requires a fresh test execution:

```bash
GOEXPERIMENT=jsonv2 go test -count=1 ./...
```

This is particularly useful when tests interact with files, temporary services, or other external state.

### 3.5 Test with a timeout

Set an explicit timeout for test execution:

```bash
GOEXPERIMENT=jsonv2 go test \
  -timeout 10m \
  ./...
```

Increase the timeout only when the test suite legitimately requires it. A long timeout should not be used to conceal a test that is stuck.

## 4. Test Design Guidelines

### 4.1 Test behavior, not implementation details

Tests should verify observable behavior and important invariants.

Prefer testing:

- Expected output.
- Returned errors.
- Configuration validation.
- State transitions.
- Resource cleanup.
- Protocol compatibility.
- Failure handling.

Avoid tests that merely reproduce internal implementation details and become fragile whenever code is refactored.

### 4.2 Table-driven tests

Use table-driven tests when the same behavior needs to be validated with multiple inputs.

Example:

```go
func TestValidatePort(t *testing.T) {
    tests := []struct {
        name    string
        port    int
        wantErr bool
    }{
        {
            name:    "valid port",
            port:    443,
            wantErr: false,
        },
        {
            name:    "zero port",
            port:    0,
            wantErr: true,
        },
        {
            name:    "negative port",
            port:    -1,
            wantErr: true,
        },
        {
            name:    "port out of range",
            port:    65536,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidatePort(tt.port)

            if (err != nil) != tt.wantErr {
                t.Fatalf(
                    "ValidatePort(%d) error = %v, wantErr %v",
                    tt.port,
                    err,
                    tt.wantErr,
                )
            }
        })
    }
}
```

The function above is an illustrative example. Adapt the test to the actual function and package under development.

### 4.3 Test error paths

Tests should cover both successful and unsuccessful execution.

For example:

- Missing required configuration.
- Invalid configuration values.
- Panel API request failures.
- Authentication failures.
- Unsupported protocol settings.
- Certificate loading failures.
- Core initialization errors.
- Resource cleanup after partial initialization.

Error paths are particularly important in node management code because a partially initialized core or failed panel request can leave background resources running.

### 4.4 Keep tests deterministic

Tests should avoid relying on:

- Current wall-clock time without a controlled clock.
- Random values without fixed seeds.
- External DNS responses.
- Public API availability.
- Uncontrolled network conditions.
- Host-specific filesystem paths.

When external behavior is unavoidable, isolate it behind an interface or another testable boundary.

## 5. Recommended Package Testing Coverage

The following areas are recommended priorities for building and improving the project's test suite.

They describe desired coverage, not a claim that every listed test currently exists.

### 5.1 CLI and command handling

Relevant area:

```text
cmd/
```

Recommended test cases:

- Command registration.
- Argument and flag parsing.
- Default configuration paths.
- Invalid command-line arguments.
- Version output.
- Server startup error handling.
- Configuration loading failures.

The CLI should fail clearly when required configuration is missing or invalid.

### 5.2 Configuration management

Relevant area:

```text
conf/
```

Recommended test cases:

- Valid configuration loading.
- Invalid JSON5 syntax.
- Missing required fields.
- Default values.
- Invalid core names.
- Invalid node settings.
- Invalid TLS configuration.
- Configuration serialization, where applicable.

Configuration tests should cover the difference between omitted fields, explicit zero values, and invalid values.

### 5.3 Panel API integration

Relevant area:

```text
api/panel/
```

Recommended test cases:

- Authentication and request headers.
- Successful API responses.
- HTTP error responses.
- Invalid response bodies.
- Timeouts and connection failures.
- User and node data parsing.
- Panel communication errors.
- Retry behavior, where implemented.

Where possible, use a local HTTP test server instead of a real panel.

Tests must not depend on a production panel or production API credentials.

### 5.4 Proxy core management

Relevant area:

```text
core/
```

Recommended test cases:

- Core registration and selection.
- Core initialization.
- Invalid core configuration.
- Startup and shutdown lifecycle.
- Resource cleanup.
- Core-specific error propagation.
- Compatibility with enabled build tags.

Special attention should be given to differences between Xray, sing-box, and Hysteria2 integrations.

### 5.5 Node implementations

Relevant area:

```text
node/
```

Recommended test cases:

- Node configuration parsing.
- Protocol-specific initialization.
- User and credential synchronization.
- Node startup and shutdown.
- Invalid protocol parameters.
- Traffic and connection handling.
- TLS certificate loading.
- Failed initialization cleanup.

Node tests should verify that invalid configuration does not leave partially initialized resources behind.

### 5.6 Shared utilities and traffic limiting

Relevant areas:

```text
common/
limiter/
```

Recommended test cases:

- Input validation.
- Utility function edge cases.
- Traffic accounting.
- User and node identification.
- Concurrent access.
- State cleanup.
- Error handling.

For concurrent code, include race detection and tests that exercise simultaneous reads and writes.

## 6. Build Configuration Testing

V2bX supports multiple proxy cores and optional features through Go build tags.

The current documented full build uses:

```text
Core tags:
- sing
- xray
- hysteria2

Optional feature tags:
- with_quic
- with_grpc
- with_utls
- with_wireguard
- with_acme
- with_gvisor
```

These tags affect which source files and dependencies are included in a build.

A successful build with one tag combination does not prove that another combination works.

### 6.1 Full build validation

Use the project's documented release build configuration:

```bash
version="test"
mkdir -p build_assets

GOEXPERIMENT=jsonv2 go build \
  -v \
  -o build_assets/V2bX \
  -tags "sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor" \
  -trimpath \
  -ldflags "-X 'github.com/proxynav/V2bX/cmd.version=$version' -s -w -buildid="
```

### 6.2 Core-specific build validation

When modifying a core integration, test the corresponding build configuration.

For example:

```bash
GOEXPERIMENT=jsonv2 go build \
  -tags "xray" \
  ./...
```

```bash
GOEXPERIMENT=jsonv2 go build \
  -tags "sing" \
  ./...
```

```bash
GOEXPERIMENT=jsonv2 go build \
  -tags "hysteria2" \
  ./...
```

These are candidate validation commands. They must be checked against the actual source-level build constraints before being adopted as a supported build matrix.

Some combinations may require additional feature tags or may not be supported independently.

### 6.3 Optional feature validation

When changing optional functionality, validate the relevant feature tag together with its required core tags.

For example:

```bash
GOEXPERIMENT=jsonv2 go build \
  -tags "sing with_quic" \
  ./...
```

```bash
GOEXPERIMENT=jsonv2 go build \
  -tags "sing with_acme" \
  ./...
```

Do not assume that every optional feature can be enabled independently of the core that uses it.

### 6.4 Build matrix

The following matrix is a recommended starting point:

| Build target          | Validation                              |
| --------------------- | --------------------------------------- |
| Full build            | All documented core and optional tags   |
| Xray integration      | Xray-related source and build tags      |
| sing-box integration  | sing-related source and build tags      |
| Hysteria2 integration | Hysteria2-related source and build tags |
| QUIC                  | Relevant core plus `with_quic`          |
| gRPC                  | Relevant core plus `with_grpc`          |
| uTLS                  | Relevant core plus `with_utls`          |
| WireGuard             | Relevant core plus `with_wireguard`     |
| ACME                  | Relevant core plus `with_acme`          |
| gVisor                | Relevant core plus `with_gvisor`        |

The exact supported combinations should be confirmed by inspecting source build constraints and successful compilation. Add only verified combinations to the project's official CI matrix.

## 7. Concurrency and Race Detection

Concurrency is important in a long-running proxy node process.

Potential areas include:

- User synchronization.
- Traffic accounting.
- Connection management.
- Core lifecycle.
- Periodic tasks.
- Configuration updates.
- Shared state.

Run the Go race detector where supported:

```bash
GOEXPERIMENT=jsonv2 go test -race ./...
```

The race detector requires compatible build and runtime environments. Some dependencies or target platforms may impose additional requirements.

For local troubleshooting, run a focused package:

```bash
GOEXPERIMENT=jsonv2 go test \
  -race \
  -count=1 \
  ./path/to/package
```

Do not interpret a successful race-enabled test run as proof that every production concurrency path is race-free.

## 8. Integration Testing

Integration tests verify behavior across component boundaries.

They should be separated from ordinary unit tests because they may require additional services, privileges, or network resources.

### 8.1 Panel integration

Use a dedicated test panel or a controlled HTTP test server.

Validate:

- API authentication.
- Node information retrieval.
- User synchronization.
- Online status reporting.
- Traffic reporting.
- API timeout and failure handling.

Never run destructive integration tests against production panel data.

### 8.2 Proxy core integration

For core-related changes, use a dedicated test configuration and isolated ports.

Validate:

- Core initialization.
- Inbound creation.
- User configuration.
- Outbound and routing behavior.
- Connection establishment.
- Shutdown and restart.
- Cleanup after failed startup.

Where practical, use a local client to establish real connections through the test node.

### 8.3 TLS and certificate integration

For TLS-related changes, validate:

- Certificate and private-key loading.
- Certificate hostname matching.
- Certificate expiration handling.
- TLS handshake success.
- Invalid certificate handling.
- SNI behavior.
- ECH configuration and client compatibility, where supported.

For ECH-related changes, distinguish between successful TLS connectivity and successful ECH negotiation. A normal TLS handshake alone does not establish that ECH was used.

### 8.4 Network protocol integration

For network-related changes, validate the actual transport and protocol behavior.

Examples include:

- TCP connectivity.
- UDP forwarding.
- QUIC handshakes.
- gRPC transport.
- HTTP/2 behavior.
- DNS resolution.
- IPv4 and IPv6 behavior.

Use packet captures and socket inspection when necessary.

Useful Linux commands include:

```bash
ss -lntup
```

```bash
tcpdump -ni any
```

```bash
ip route
```

These checks complement automated tests; they do not replace them.

## 9. Protocol Compatibility Testing

Protocol changes should include a protocol-specific validation plan.

The following table provides a recommended baseline.

| Protocol or feature | Suggested validation                                   |
| ------------------- | ------------------------------------------------------ |
| VLESS               | Authentication, transport, TLS and connection          |
| VMess               | Authentication, transport and connection               |
| Trojan              | TLS certificate, SNI and authentication                |
| Shadowsocks         | Cipher configuration and authentication                |
| Hysteria2           | QUIC, authentication and UDP behavior                  |
| AnyTLS              | TLS handshake, authentication and client compatibility |
| ECH                 | Client ECH negotiation and fallback behavior           |
| Reality             | Handshake and client compatibility                     |
| gRPC                | HTTP/2 transport and service connectivity              |
| WireGuard           | Tunnel establishment and packet forwarding             |

The table is a test-planning reference. Actual protocol support and optional features must be confirmed against the current implementation and enabled build configuration.

For every protocol change, test both:

1. The expected successful connection.
2. Relevant failure cases, such as invalid credentials, invalid TLS configuration, or unsupported transport parameters.

When compatibility with third-party clients is important, include at least one representative client implementation in the integration test plan.

## 10. Docker Build and Runtime Testing

The current Dockerfile builds the application using Go 1.25.0 and enables:

```text
sing xray hysteria2
with_quic with_grpc with_utls
with_wireguard with_acme with_gvisor
```

It also sets:

```dockerfile
ENV CGO_ENABLED=0
```

The runtime image is based on Alpine and installs timezone and CA certificate packages.

Validate the Docker build:

```bash
docker build -t proxynav/v2bx:test .
```

Check the resulting image:

```bash
docker image inspect proxynav/v2bx:test
```

Run the container with a dedicated test configuration:

```bash
docker run --rm \
  -v "$(pwd)/test-config:/etc/V2bX:ro" \
  proxynav/v2bx:test
```

The example assumes that `test-config` contains a valid configuration at the path expected by the current Docker entrypoint.

Validate:

- Container startup.
- Configuration loading.
- Core initialization.
- Logging.
- Network port availability.
- Graceful shutdown.
- Restart behavior.

Do not assume that a successful Docker build proves that the container can start with a particular production configuration.

## 11. Performance and Load Testing

Performance tests should be introduced when modifying:

- Connection handling.
- Traffic accounting.
- User synchronization.
- Core lifecycle.
- Network forwarding.
- Memory-intensive data structures.
- Concurrent shared state.

### 11.1 Go benchmarks

Use Go benchmarks for isolated functions:

```go
func BenchmarkExample(b *testing.B) {
    for b.Loop() {
        ExampleFunction()
    }
}
```

Run benchmarks:

```bash
GOEXPERIMENT=jsonv2 go test \
  -bench=. \
  -benchmem \
  ./path/to/package
```

Use representative input sizes and compare results against a known baseline.

### 11.2 Runtime load tests

Runtime load tests should use dedicated test nodes and controlled traffic generators.

Monitor:

- CPU utilization.
- Memory consumption.
- Goroutine count.
- Connection count.
- Traffic throughput.
- Connection establishment latency.
- Error rates.
- File descriptor usage.

Avoid running uncontrolled load tests against production nodes or third-party services.

## 12. Debugging Failed Tests

### 12.1 Increase test verbosity

```bash
GOEXPERIMENT=jsonv2 go test -v ./path/to/package
```

### 12.2 Disable caching

```bash
GOEXPERIMENT=jsonv2 go test -count=1 ./...
```

### 12.3 Inspect package dependencies

```bash
go list -deps ./path/to/package
```

### 12.4 Inspect build constraints

When a package behaves differently under different build tags, inspect source constraints:

```bash
grep -R -nE '^//go:build|^// \\+build' \
  --include='*.go' .
```

Use the output to determine which files are included in each build configuration.

### 12.5 Inspect environment and module state

```bash
go version
go env
go mod verify
git submodule status --recursive
```

### 12.6 Capture test output

```bash
GOEXPERIMENT=jsonv2 go test \
  -json \
  ./... > test-results.json
```

This is useful when analyzing failures across multiple packages or preserving test output in CI.

## 13. Continuous Integration Recommendations

The repository's actual CI configuration should be treated as authoritative. The following is a recommended CI validation plan, not a claim that every job is currently configured.

### 13.1 Pull request checks

For every pull request, consider running:

1. Go module verification.
2. Module consistency checks.
3. Unit tests.
4. Static analysis.
5. Full documented build.
6. Relevant core-specific builds.

Example:

```bash
set -euo pipefail

export GOEXPERIMENT=jsonv2

go mod verify
go mod tidy
git diff --exit-code -- go.mod go.sum

go test ./...
go vet ./...

mkdir -p build_assets

version="ci"
go build \
  -v \
  -o build_assets/V2bX \
  -tags "sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor" \
  -trimpath \
  -ldflags "-X 'github.com/proxynav/V2bX/cmd.version=$version' -s -w -buildid="
```

The `go mod tidy` step should be run in a clean working tree. It is a validation step: unexpected modifications should fail the check and be reviewed.

### 13.2 Extended CI checks

For scheduled workflows or important releases, consider adding:

- Race detector tests.
- Additional supported build-tag combinations.
- Docker image builds.
- Dependency vulnerability scanning.
- Integration tests with dedicated services.
- Runtime protocol compatibility tests.

### 13.3 Test matrix management

Avoid creating an unnecessarily large CI matrix.

Start with the configurations that represent supported release builds and the core integrations affected by changes.

Expand the matrix when a new core, protocol, or build-tag combination becomes an officially supported target.

### 13.4 CI failure handling

A failed check must not be ignored merely because another build configuration succeeds.

When a failure is platform-specific or caused by unavailable integration infrastructure:

- Record the failing job and environment.
- Identify whether it is a code regression or an infrastructure issue.
- Preserve the failure output.
- Document any temporary exception.
- Ensure the affected validation is performed before release.

## 14. Test Data and Fixtures

Store reusable test data in an appropriate package-local `testdata/` directory.

Example:

```text
package/
├── config.go
├── config_test.go
└── testdata/
    ├── valid-config.json
    ├── invalid-config.json
    └── sample-response.json
```

Go ignores directories named `testdata` when enumerating packages, making them suitable for test fixtures.

Test fixture guidelines:

- Use minimal representative configuration.
- Keep valid and invalid examples separate.
- Do not include real credentials.
- Do not include production private keys.
- Document fixtures that require special interpretation.
- Avoid duplicating large configuration files unnecessarily.

For certificate tests, generate dedicated test certificates or use non-production fixtures with clearly documented limitations.

## 15. Pull Request Testing Requirements

Every code change should include an appropriate validation plan.

### General changes

- [ ] Affected package tests pass.
- [ ] Relevant error paths are tested.
- [ ] No unrelated test or fixture changes are included.
- [ ] Full Go tests pass when applicable.

### Core or protocol changes

- [ ] Relevant core build succeeds.
- [ ] Full documented build succeeds.
- [ ] Protocol-specific tests are performed.
- [ ] TLS and transport behavior are checked when affected.
- [ ] Runtime behavior is validated in a test environment.

### Dependency changes

- [ ] Module integrity is verified.
- [ ] Dependency changes are reviewed.
- [ ] Relevant build configurations succeed.
- [ ] Runtime compatibility is checked for critical dependencies.

### Release changes

- [ ] Release build succeeds.
- [ ] Docker build succeeds.
- [ ] Container startup is validated.
- [ ] Critical protocols are smoke-tested.
- [ ] Known limitations are documented.

## 16. Test Coverage Roadmap

The following roadmap is recommended for gradually improving the project's test coverage.

| Priority | Area                 | Goal                                               |
| -------- | -------------------- | -------------------------------------------------- |
| High     | Configuration        | Validate parsing, defaults, and invalid inputs     |
| High     | Panel API            | Cover request handling and error responses         |
| High     | Core lifecycle       | Verify initialization and cleanup                  |
| High     | Node lifecycle       | Verify startup, shutdown, and failure handling     |
| High     | Build matrix         | Validate supported core and feature combinations   |
| Medium   | Traffic limiter      | Verify accounting and concurrent access            |
| Medium   | TLS and certificates | Cover loading and handshake behavior               |
| Medium   | Protocol integration | Verify representative client connections           |
| Medium   | Docker runtime       | Verify startup and configuration loading           |
| Later    | Performance          | Establish repeatable benchmarks and load baselines |

Prioritize tests that protect critical runtime behavior and prevent regressions in node startup, panel synchronization, and protocol connectivity.

## 17. References

- [ProxyNav V2bX Repository](https://github.com/proxynav/V2bX)
- [Go Testing Documentation](https://go.dev/doc/tutorial/add-a-test)
- [Go Test Command Reference](https://pkg.go.dev/cmd/go#hdr-Test_packages)
- [Go Test Package](https://pkg.go.dev/testing)
- [Go Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [GitHub Actions: Building and Testing Go](https://docs.github.com/en/actions/tutorials/build-and-test-code/go)
