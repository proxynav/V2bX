# Debugging Guide

This document describes practical debugging techniques for **ProxyNav V2bX**, a multi-core proxy node management application.

It is intended for developers investigating build failures, configuration errors, panel communication problems, proxy core startup failures, protocol issues, and runtime problems.

## 1. Project Repository and Code Structure

The official project repository is:

- **ProxyNav V2bX:** https://github.com/proxynav/V2bX
- **Original fork source:** https://github.com/wyx2685/V2bX

Use the official ProxyNav repository when reporting issues, submitting pull requests, or referencing project documentation.

The codebase is organized into several areas. When debugging, start with the module closest to the observed failure.

| Directory    | Debugging focus                                                   |
| ------------ | ----------------------------------------------------------------- |
| `cmd/`       | Application entry point, command initialization, and startup flow |
| `conf/`      | Configuration loading, parsing, and configuration structures      |
| `api/panel/` | Panel API communication and response handling                     |
| `core/`      | Proxy core integration and lifecycle management                   |
| `node/`      | Node configuration and node-related processing                    |
| `limiter/`   | Connection, traffic, and rate-limiting behavior                   |
| `common/`    | Shared utilities and common functionality                         |
| `example/`   | Example configurations and usage references                       |
| `test_data/` | Test data and debugging fixtures                                  |

The application entry point is `main.go`, which invokes the command package.

When investigating a problem, trace the execution path from the entry point to the affected module rather than debugging unrelated components.

## 2. Establish a Reproducible Environment

Before debugging, record the application version, build configuration, operating system, and runtime environment.

### 2.1 Record the Go Environment

```bash
go version
go env
go env GOOS GOARCH CGO_ENABLED
```

The current project `go.mod` declares Go 1.25 and toolchain 1.25.0. Use a compatible Go toolchain when building the project.

### 2.2 Record the Project Revision

```bash
git remote -v
git branch --show-current
git rev-parse HEAD
git status --short
git log -1 --oneline
```

Make sure the working tree and branch are known before comparing behavior or reproducing a bug.

### 2.3 Record Build Tags and Build Options

V2bX supports conditional compilation for different proxy cores and optional features.

The selected build tags can affect which implementations are included in the binary. A problem that occurs in one build may not occur in another.

Inspect the build command and record:

- Enabled proxy core tags, such as `xray`, `sing`, and `hysteria2`.
- Optional feature tags, such as `with_acme`, `with_quic`, and `with_grpc`.
- The `GOEXPERIMENT` setting.
- The Go toolchain version.
- Any linker flags used to inject version information.

Do not assume that a binary contains every core or optional feature.

## 3. Build and Static Analysis

Start with compilation and static analysis before investigating runtime behavior.

### 3.1 Check Dependencies

```bash
go mod download
go mod verify
```

If dependency resolution fails, inspect:

```bash
go env GOPROXY GOSUMDB GOPRIVATE
```

Do not run `go mod tidy` as an initial debugging step in an established working tree. It may modify `go.mod` and `go.sum`. Review and preserve intentional dependency changes.

### 3.2 Run Static Analysis

```bash
gofmt -w .
go vet ./...
```

Use `gofmt -w .` only when formatting changes are intended. During investigation, you can check formatting without modifying files:

```bash
gofmt -l .
```

### 3.3 Run Tests

```bash
go test ./...
go test -v ./...
```

Run a specific package when narrowing down a failure:

```bash
go test -v ./core/...
```

Replace the package path with the actual package being investigated.

For race detection:

```bash
go test -race ./...
```

Race detection may require additional compiler and platform support. If it is unavailable in the target environment, run it in a supported development environment.

### 3.4 Build with the Required Features

The project README documents a build using conditional tags and `GOEXPERIMENT=jsonv2`.

A representative build command is:

```bash
GOEXPERIMENT=jsonv2 go build -v \
  -o build_assets/V2bX \
  -tags "sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor" \
  -trimpath \
  -ldflags "-s -w"
```

Adjust the tags to match the intended release configuration.

For debugging, avoid stripping symbols with `-s -w` because they make stack traces and debugger inspection less convenient.

Use the same build tags as the affected production binary whenever possible. Otherwise, you may fail to reproduce a problem caused by conditional compilation.

### 3.5 Inspect the Resulting Binary

```bash
file ./build_assets/V2bX
go version -m ./build_assets/V2bX
```

These commands help identify the binary architecture, Go build information, and embedded module metadata.

## 4. Configuration Debugging

V2bX loads its application configuration through the `conf` package.

The configuration structure includes:

- `Log`: logging configuration.
- `Cores`: proxy core configuration.
- `Nodes`: node configuration.

The configuration loader reads the file, processes it through the project's JSON5 reader, and then unmarshals the result into the application configuration structure.

### 4.1 Check the Configuration File

First, verify that the application is using the intended file.

```bash
ls -l /path/to/config.json
stat /path/to/config.json
```

Check:

- Whether the file exists.
- Whether the service account can read it.
- Whether the file is the expected version.
- Whether the service is using the correct configuration path.
- Whether the configuration contains the expected `Log`, `Cores`, and `Nodes` sections.

Do not assume that a configuration file in the current working directory is automatically used by the service.

### 4.2 Check JSON5 Parsing Errors

Configuration loading can fail because of malformed JSON5 input or incompatible field values.

Typical problems include:

- Missing or incorrectly placed delimiters.
- Invalid strings or escape sequences.
- Incorrect field names.
- Values with incompatible types.
- Invalid nested objects or arrays.

When a parsing error occurs, inspect the complete error message and the configuration around the reported location.

For a quick syntax check, use a JSON5-compatible validator. A regular JSON validator may reject valid JSON5 syntax, such as comments or unquoted property names.

### 4.3 Verify Configuration Semantics

Successful parsing does not guarantee a valid runtime configuration.

Check the configuration against the corresponding structures and initialization code in:

- `conf/`
- `core/`
- `node/`

In particular, verify that each configured core and node references valid settings and that all required fields are present.

When a configuration loads successfully but startup fails, continue tracing the initialization path instead of treating the parser as the source of the problem.

### 4.4 Protect Sensitive Configuration

Configuration files may contain panel credentials, API tokens, passwords, and TLS-related information.

Before sharing a configuration in an issue or debugging session:

- Remove passwords and tokens.
- Replace private addresses and domains when necessary.
- Do not publish private keys or certificate material.
- Preserve field names and structural relationships so the issue remains reproducible.

## 5. Application Logging

Logs are usually the first source of information for startup and runtime failures.

### 5.1 Check the Configured Log Level

Inspect the `Log` section of the configuration and verify the configured log level and output destination.

During development, use a sufficiently detailed log level if supported by the current logging implementation.

Avoid enabling verbose logging indefinitely in production because it may generate large amounts of data or expose sensitive information.

### 5.2 Run the Application in the Foreground

When possible, reproduce the issue by running the binary directly with the same configuration and environment as the production service.

For example:

```bash
/path/to/V2bX --help
```

Use the actual command-line options shown by the application to specify its configuration file or other runtime settings.

Do not assume a particular configuration flag name without checking the command implementation.

Running in the foreground is useful because it makes startup errors and process termination easier to observe.

### 5.3 Inspect systemd Logs

If V2bX is managed by systemd, first identify the actual unit name:

```bash
systemctl list-units --type=service | grep -i v2
```

Then inspect its status and logs:

```bash
systemctl status <service-name>
journalctl -u <service-name> -n 200 --no-pager
journalctl -u <service-name> -f
```

Check the service definition when the process behaves differently under systemd:

```bash
systemctl cat <service-name>
```

Pay attention to:

- `ExecStart`
- `WorkingDirectory`
- `User` and `Group`
- Environment variables
- File and directory permissions
- Restart policies

A service may use a different configuration path or environment from an interactive shell.

## 6. Panel API Debugging

The `api/panel/` package is the starting point for investigating panel communication problems.

Typical symptoms include:

- Failed panel authentication.
- HTTP request errors.
- Unexpected API responses.
- Node information not updating.
- User or traffic information not synchronizing.
- Unexpected behavior after a panel-side configuration change.

### 6.1 Verify Network Connectivity

Check DNS resolution:

```bash
dig +short panel.example.com
```

Check TCP connectivity:

```bash
nc -vz panel.example.com 443
```

Check HTTPS connectivity:

```bash
curl -v https://panel.example.com/
```

Replace the hostname, port, and path with the actual panel endpoint.

A successful connection to the panel homepage does not prove that the API endpoint or authentication flow is working.

### 6.2 Inspect HTTP Failures

When an API request fails, record:

- Request URL and HTTP method.
- HTTP status code.
- Request timeout or connection error.
- Response body, with sensitive fields removed.
- Relevant request and response headers.
- Whether the failure is consistent or intermittent.

Common HTTP status codes:

| Status                | Possible cause                                 |
| --------------------- | ---------------------------------------------- |
| `400`                 | Invalid request or parameters                  |
| `401`                 | Authentication failure                         |
| `403`                 | Permission or access restriction               |
| `404`                 | Incorrect API endpoint                         |
| `429`                 | Rate limiting                                  |
| `500`                 | Panel-side internal error                      |
| `502` / `503` / `504` | Reverse proxy or upstream availability problem |

Inspect the actual API response and the corresponding code path before concluding that a status code has a specific meaning for V2bX.

### 6.3 Trace the Request Lifecycle

When debugging panel integration:

1. Find the relevant API method in `api/panel/`.
2. Identify how the request URL, headers, and parameters are constructed.
3. Inspect response status handling and response decoding.
4. Follow the returned data into `node/` or other consuming modules.
5. Check whether the resulting state is applied to the relevant core.

Avoid logging authentication tokens, passwords, or complete user credentials.

## 7. Proxy Core Debugging

V2bX integrates multiple proxy cores. Core-related failures should be investigated separately from panel communication and node configuration.

The main areas to inspect are:

- Core configuration and initialization in `core/`.
- Node configuration and runtime integration in `node/`.
- Build tags and conditional compilation.
- Core-specific configuration generation and startup behavior.

### 7.1 Confirm the Core Was Included in the Binary

Check the build command and tags used to produce the binary.

For example, a build that omits the `sing` tag may not contain the corresponding sing-box implementation.

Rebuild with the required tags before investigating a missing-core problem.

### 7.2 Check Core Initialization

When a core fails to start, inspect:

- Whether the core implementation was compiled in.
- Whether the configuration is valid.
- Whether required files and directories exist.
- Whether the process has permission to access those files.
- Whether the required ports are available.
- Whether the core reports an initialization or configuration error.

Follow the initialization code in `core/` and the corresponding node startup path rather than assuming every core uses the same lifecycle.

### 7.3 Validate Generated Core Configuration

If V2bX generates configuration for an underlying core, inspect the generated configuration when possible.

Use the configuration validation command supported by the exact core version included in the build.

Do not assume that Xray and sing-box accept the same configuration format or validation flags.

When a core configuration fails:

1. Capture the generated configuration.
2. Remove credentials and other sensitive information.
3. Validate it with the matching core version.
4. Compare the generated configuration with the expected configuration.
5. Trace the generation code to identify incorrect field mapping or missing values.

### 7.4 Check Core Ports and Processes

Inspect listening sockets:

```bash
ss -lntup
```

Search for a specific port:

```bash
ss -lntup | grep ':443'
```

Check running processes:

```bash
ps aux | grep -E 'V2bX|xray|sing-box'
```

A port conflict, incorrect bind address, or insufficient permission can prevent a core from starting.

## 8. Protocol and TLS Troubleshooting

Protocol failures may originate from the client, network, TLS configuration, generated core configuration, or the underlying core.

Start by identifying the exact protocol, node, core, and transport involved.

### 8.1 Check DNS Resolution

```bash
dig A node.example.com
dig AAAA node.example.com
```

Verify that the returned addresses point to the intended server.

If IPv4 and IPv6 are both configured, test them independently. A working IPv4 path does not guarantee that the IPv6 path is correctly configured.

### 8.2 Check TCP and UDP Connectivity

For TCP:

```bash
nc -vz node.example.com 443
```

For UDP, a successful local socket check does not prove that packets reach the server.

Inspect listening sockets:

```bash
ss -lntup
```

Capture traffic when necessary:

```bash
tcpdump -ni any port 443
```

Replace the port and protocol filter with the actual node configuration.

### 8.3 Inspect TLS Certificates

Use OpenSSL to inspect a TLS endpoint:

```bash
openssl s_client \
  -connect node.example.com:443 \
  -servername node.example.com \
  -showcerts \
  -brief </dev/null
```

Check:

- Whether the TCP connection succeeds.
- Whether a TLS handshake completes.
- Whether the certificate matches the expected hostname.
- Whether the certificate chain verifies.
- Whether the negotiated TLS version and cipher suite are acceptable.

A successful TLS handshake does not prove that the proxy protocol or application-level authentication is working.

### 8.4 Investigate ECH and TLS Fingerprint Issues

When troubleshooting ECH-related or TLS fingerprint-related problems, distinguish between:

- DNS resolution and HTTPS/SVCB record availability.
- Client ECH support.
- Server-side TLS configuration.
- The selected core and its ECH implementation.
- SNI and certificate configuration.
- Client fingerprint and ALPN settings.

Compare the actual client configuration with the configuration generated by V2bX.

Do not infer that ECH is active merely because a TLS handshake succeeds. Confirm the relevant client and server behavior using tools and logs that support ECH inspection.

## 9. Network and Firewall Diagnostics

When a node is reachable from the server but not from clients, investigate the complete packet path.

### 9.1 Check Listening Ports

```bash
ss -lntup
```

Verify that the expected service is listening on the correct address and port.

A service bound to `127.0.0.1` is not directly reachable through the server's public interface.

### 9.2 Inspect Firewall Rules

For iptables:

```bash
iptables -L -n -v
iptables -t nat -L -n -v
```

For nftables:

```bash
nft list ruleset
```

Check whether the relevant traffic is accepted, rejected, dropped, or redirected.

### 9.3 Capture Packets

```bash
tcpdump -ni any host <client-ip>
```

For a specific port:

```bash
tcpdump -ni any 'port 443'
```

Packet captures help determine whether:

- Packets reach the server.
- The expected interface receives traffic.
- The server sends a response.
- Traffic is dropped or redirected before reaching the application.

### 9.4 Check Routing

```bash
ip addr
ip route
ip rule
```

For more detailed route selection:

```bash
ip route get <destination-ip>
```

Consider policy routing, NAT, reverse-path filtering, and upstream firewall rules when traffic reaches the host but does not follow the expected path.

## 10. Runtime and Performance Debugging

Use Go runtime tools when the application starts successfully but experiences crashes, high CPU usage, memory growth, or unexpected concurrency behavior.

### 10.1 Collect CPU and Memory Information

```bash
ps -o pid,ppid,%cpu,%mem,rss,vsz,cmd -C V2bX
```

Use the actual process name if it differs.

Inspect process resource limits:

```bash
cat /proc/<pid>/limits
```

Check system memory:

```bash
free -h
```

### 10.2 Inspect Goroutines

If the application exits unexpectedly or appears blocked, capture the available logs and stack traces.

For a Go process, sending `SIGQUIT` typically causes the runtime to print goroutine stack traces and terminate the process.

```bash
kill -QUIT <pid>
```

**Warning:** This terminates the process. Do not use it casually on production nodes. Prefer reproducing the problem in a test environment or using a debugger when process continuity is required.

### 10.3 Use the Race Detector

Run tests with:

```bash
go test -race ./...
```

The race detector can identify data races exercised by the test workload.

It does not prove that the application is race-free if the relevant code paths are not executed.

### 10.4 Use Delve

Install Delve:

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

Build a debugging binary without stripping symbols:

```bash
GOEXPERIMENT=jsonv2 go build \
  -gcflags="all=-N -l" \
  -o ./build_assets/V2bX-debug \
  .
```

If the application requires specific build tags, include the same tags in this command.

Start the debugger:

```bash
dlv exec ./build_assets/V2bX-debug -- --help
```

Replace `--help` with the actual application arguments needed to reproduce the problem.

Useful Delve commands:

```text
break main.main
continue
next
step
print variableName
locals
stack
goroutines
```

Use breakpoints in the relevant module to inspect configuration loading, panel responses, node initialization, or core lifecycle behavior.

### 10.5 CPU and Heap Profiling

If profiling is needed, use a controlled development build and an explicitly configured profiling mechanism.

Before using `net/http/pprof`, verify that the application actually exposes a pprof endpoint. Do not assume that such an endpoint is enabled by default.

Never expose profiling endpoints publicly on production nodes. CPU profiles, heap profiles, and goroutine dumps may contain sensitive runtime information.

## 11. Common Troubleshooting Scenarios

### 11.1 Application Exits Immediately

Check:

1. The executable architecture and permissions.
2. The Go runtime and build information.
3. The configuration file path.
4. Configuration parsing errors.
5. Required environment variables.
6. Core initialization errors.
7. systemd `ExecStart` and `WorkingDirectory`.

Useful commands:

```bash
file /path/to/V2bX
/path/to/V2bX --help
journalctl -u <service-name> -n 200 --no-pager
```

### 11.2 Node Does Not Start

Check:

1. Whether the node configuration is loaded.
2. Whether the required core was included at build time.
3. Whether the node's ports are already occupied.
4. Whether TLS certificate files exist and are readable.
5. Whether generated core configuration is valid.
6. Whether the underlying core reports a startup error.

### 11.3 Panel Data Is Not Updating

Check:

1. Panel DNS resolution and connectivity.
2. API authentication and permissions.
3. HTTP status codes and response bodies.
4. API response decoding.
5. The code path that applies panel data to node state.
6. Whether the relevant core receives the updated configuration.

### 11.4 Client Cannot Connect

Check:

1. DNS records and address selection.
2. TCP or UDP reachability.
3. Listening addresses and ports.
4. Firewall and NAT rules.
5. TLS certificate and SNI settings.
6. Client and server protocol compatibility.
7. The underlying core's logs.

### 11.5 Build Works Locally but Fails in CI

Compare:

- Go toolchain versions.
- `GOEXPERIMENT` settings.
- Build tags.
- Target OS and architecture.
- CGO settings and required system toolchains.
- Dependency cache and module download behavior.
- Linker flags and version-injection settings.

Reproduce the CI build command locally before changing application code.

## 12. Reporting a Bug

When opening an issue in the official repository, provide enough information to reproduce the problem.

Include:

- V2bX version or Git commit.
- Operating system and architecture.
- Go version, if the issue is build-related.
- Build tags and relevant build options.
- Selected proxy core and protocol.
- Sanitized configuration.
- Complete error messages and relevant logs.
- Steps to reproduce.
- Expected behavior.
- Actual behavior.
- Whether the problem is reproducible with the latest project version.

Remove passwords, tokens, private keys, and other sensitive information before submitting logs or configuration files.

Report issues and follow development discussions here:

https://github.com/proxynav/V2bX/issues

## 13. Further Reading

- [ProxyNav V2bX Repository](https://github.com/proxynav/V2bX)
- [Go Documentation](https://go.dev/doc/)
- [Go Command Documentation](https://pkg.go.dev/cmd/go)
- [Delve Debugger](https://github.com/go-delve/delve)
- [Xray-core](https://github.com/XTLS/Xray-core)
- [sing-box](https://github.com/SagerNet/sing-box)

---

**Debugging principle:** Identify the failing layer first, reproduce the issue with the same build configuration, and trace the execution path through the relevant V2bX module. Avoid changing multiple components at once, and verify each proposed fix with a reproducible test.
