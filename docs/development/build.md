# Building ProxyNav V2bX

This document describes how to build **ProxyNav V2bX** from source, including local development builds, release builds, build tags, Docker images, and cross-compilation considerations.

The instructions are based on the official repository's current build configuration.

## 1. Project Repository

The official project repository is:

- **ProxyNav V2bX:** https://github.com/proxynav/V2bX
- **Original fork source:** https://github.com/wyx2685/V2bX

The official repository is the source of truth for new development and releases.

The current source retains some module paths and build settings inherited from the original project. Do not assume that the GitHub repository URL and the Go module path are identical.

## 2. Build Prerequisites

### 2.1 Go Toolchain

The current `go.mod` specifies:

```go
go 1.25
toolchain go1.25.0
```

Use Go 1.25 or a compatible newer toolchain.

Verify the installed version:

```bash
go version
go env GOOS GOARCH CGO_ENABLED
```

### 2.2 Git and Submodules

Clone the official repository with its submodules:

```bash
git clone --recurse-submodules https://github.com/proxynav/V2bX.git
cd V2bX
```

The project includes the `sing-box_mod` Git submodule.

If the repository has already been cloned without submodules, initialize them:

```bash
git submodule update --init --recursive
```

Verify the submodule state:

```bash
git submodule status
```

Builds involving sing-box functionality must use the submodule revision recorded by the main repository. Do not replace it with an arbitrary upstream revision.

### 2.3 Build Environment

The documented build configuration uses:

```bash
GOEXPERIMENT=jsonv2
```

The Dockerfile also sets:

```bash
CGO_ENABLED=0
```

The latter is part of the current Docker build configuration. Do not assume it is suitable for every future build variant or target platform without verification.

## 3. Understand the Build Configuration

V2bX uses Go build tags to select proxy core implementations and optional features.

### 3.1 Proxy Core Tags

The README identifies these core tags:

| Build tag | Purpose |
|---|---|
| `xray` | Include Xray core integration |
| `sing` | Include sing-box core integration |
| `hysteria2` | Include Hysteria2-related implementation |

The selected tags determine which implementations are compiled into the binary.

A build without a required core tag may succeed but will not contain that core's implementation.

### 3.2 Optional Feature Tags

The documented full build also enables:

| Build tag | Purpose |
|---|---|
| `with_quic` | QUIC-related functionality |
| `with_grpc` | gRPC-related functionality |
| `with_utls` | uTLS-related functionality |
| `with_wireguard` | WireGuard-related functionality |
| `with_acme` | ACME-related functionality |
| `with_gvisor` | gVisor-related functionality |

These tags should be kept consistent with the intended deployment and release configuration.

### 3.3 Full Build Tag Set

The current README documents the following full build tag set:

```text
sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor
```

Use this set when building the full-featured binary described by the project README.

Do not assume that every reduced combination is supported merely because the full build works. Test each combination used in development or release automation.

## 4. Build Locally

### 4.1 Prepare the Output Directory

The documented build command writes the binary to `build_assets/V2bX`.

Create the output directory:

```bash
mkdir -p build_assets
```

### 4.2 Build with the Project's Default Release Configuration

Set a version string:

```bash
version="dev"
```

Build the application:

```bash
GOEXPERIMENT=jsonv2 go build -v \
  -o build_assets/V2bX \
  -tags "sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor" \
  -trimpath \
  -ldflags "-X 'github.com/proxynav/V2bX/cmd.version=$version' -s -w -buildid="
```

This command follows the build configuration documented in the repository README.

The linker flags have the following purposes:

| Flag | Purpose |
|---|---|
| `-X` | Set the value of a string variable at link time |
| `-s` | Omit symbol table information |
| `-w` | Omit DWARF debugging information |
| `-buildid=` | Clear the Go build ID |

The `-X` argument uses the current module path and the version variable in `cmd`.

If the module path or version variable changes in a future revision, update the linker argument accordingly.

### 4.3 Verify the Binary

Check the generated file:

```bash
file build_assets/V2bX
```

Inspect embedded Go build information:

```bash
go version -m build_assets/V2bX
```

Confirm that the binary has the expected architecture and build configuration.

## 5. Development and Debug Builds

Release builds use `-s -w` to reduce binary size. These flags remove debugging information and should generally be omitted when investigating runtime problems.

### 5.1 Build Without Stripping Debug Information

```bash
mkdir -p build_assets

GOEXPERIMENT=jsonv2 go build -v \
  -o build_assets/V2bX-debug \
  -tags "sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor"
```

This produces a binary that is more suitable for debugging.

### 5.2 Build with Debug Compiler Settings

For more detailed source-level debugging:

```bash
GOEXPERIMENT=jsonv2 go build \
  -gcflags="all=-N -l" \
  -o build_assets/V2bX-debug \
  -tags "sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor"
```

The compiler flags disable optimizations and inlining to make debugger inspection easier.

Use this build with Delve when necessary.

### 5.3 Build a Specific Core Variant

For faster development and narrower testing, build only the required core implementation.

For example, an Xray-oriented build:

```bash
GOEXPERIMENT=jsonv2 go build \
  -o build_assets/V2bX-xray \
  -tags "xray"
```

A sing-box-oriented build:

```bash
GOEXPERIMENT=jsonv2 go build \
  -o build_assets/V2bX-sing \
  -tags "sing"
```

These are illustrative reduced builds. Additional feature tags may be required by the code paths or dependencies being tested.

Verify each reduced build against the current source before using it as a release artifact.

## 6. Validate the Source Before Building

Run the following checks before submitting a build-related change.

### 6.1 List Go Packages

```bash
go list ./...
```

If a package depends on conditional compilation, make sure the appropriate build environment and tags are selected.

### 6.2 Run Tests

```bash
GOEXPERIMENT=jsonv2 go test ./...
```

For verbose output:

```bash
GOEXPERIMENT=jsonv2 go test -v ./...
```

### 6.3 Run Static Analysis

```bash
go vet ./...
```

### 6.4 Verify Dependencies

```bash
go mod download
go mod verify
```

Review dependency changes before committing:

```bash
git diff -- go.mod go.sum
```

Avoid running `go mod tidy` as an automatic build preparation step. Only use it when dependency changes are intentional.

## 7. Docker Build

The repository includes a Dockerfile with a multi-stage build.

The current Docker build configuration:

- Uses `golang:1.25.0-alpine` as the builder image.
- Enables `CGO_ENABLED=0`.
- Sets `GOEXPERIMENT=jsonv2` for dependency download and compilation.
- Enables the documented core and optional feature tags.
- Uses Alpine as the runtime image.
- Installs timezone data and CA certificates.
- Places the binary at `/usr/local/bin/V2bX`.
- Uses `V2bX server --config /etc/V2bX/config.json` as the container entry point.

### 7.1 Prepare the Build Context

Make sure the Git submodule is initialized before building the image:

```bash
git submodule update --init --recursive
```

Docker builds from the files available in the build context. A missing submodule checkout can therefore cause build failures.

### 7.2 Build the Docker Image

From the repository root:

```bash
docker build -t proxynav/v2bx:dev .
```

### 7.3 Inspect the Image

```bash
docker image inspect proxynav/v2bx:dev
```

The Dockerfile's runtime configuration expects the application configuration at:

```text
/etc/V2bX/config.json
```

Provide the configuration through the appropriate container volume or deployment configuration.

### 7.4 Run a Test Container

For a basic startup test, provide a valid configuration file:

```bash
docker run --rm \
  -v "$PWD/config.json:/etc/V2bX/config.json:ro" \
  proxynav/v2bx:dev
```

This is a basic example. Production deployments may require additional volume mounts, network settings, capabilities, and persistent storage depending on enabled features.

### 7.5 Docker Build and Local Build Differences

The Dockerfile does not currently use the same linker flags as the README's local release command.

In particular, it does not inject the version through the documented `-X` argument or explicitly strip debug information using `-s -w`.

Keep these build paths in mind when comparing binaries produced locally and inside Docker.

## 8. Cross-Compilation

Go supports cross-compilation, but V2bX depends on multiple proxy cores and platform-specific libraries.

Do not assume that every Go-supported target is automatically supported by every selected dependency.

### 8.1 Inspect the Current Target

```bash
go env GOOS GOARCH CGO_ENABLED
```

### 8.2 Example: Linux AMD64

```bash
GOOS=linux GOARCH=amd64 \
GOEXPERIMENT=jsonv2 \
go build \
  -o build_assets/V2bX-linux-amd64 \
  -tags "sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor"
```

### 8.3 Example: Linux ARM64

```bash
GOOS=linux GOARCH=arm64 \
GOEXPERIMENT=jsonv2 \
go build \
  -o build_assets/V2bX-linux-arm64 \
  -tags "sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor"
```

### 8.4 CGO Considerations

Some dependencies may require CGO or a target-specific compiler toolchain.

Before disabling CGO, inspect the dependency requirements and test the exact build configuration.

For example:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
GOEXPERIMENT=jsonv2 \
go build -o build_assets/V2bX-linux-amd64
```

This command is only appropriate when the selected dependencies and build tags support a CGO-disabled build.

For CGO-enabled cross-compilation, install and configure a suitable cross-compiler and any required native libraries.

### 8.5 Verify Cross-Compiled Artifacts

Use:

```bash
file build_assets/V2bX-linux-amd64
go version -m build_assets/V2bX-linux-amd64
```

Whenever possible, test the binary on the target operating system and architecture.

## 9. Reproducible Builds

To make build results easier to reproduce, record the following information for every release artifact:

- Git commit hash.
- Git branch or release tag.
- Go toolchain version.
- `GOEXPERIMENT` setting.
- Build tags.
- Target OS and architecture.
- CGO setting.
- Linker flags.
- Submodule commit revisions.

Useful commands:

```bash
git rev-parse HEAD
git branch --show-current
git submodule status
go version
go env GOOS GOARCH CGO_ENABLED GOEXPERIMENT
```

Do not rely only on a version label. Two binaries with the same version string may have been built from different commits or feature configurations.

## 10. Troubleshooting

### 10.1 Go Version or Toolchain Errors

Check:

```bash
go version
go env GOTOOLCHAIN GOEXPERIMENT
```

Use a compatible Go 1.25 toolchain and enable the required JSON experiment setting.

### 10.2 Missing Submodule Files

If compilation fails because code under `sing-box_mod` is missing:

```bash
git submodule update --init --recursive
git submodule status
```

Confirm that the submodule is checked out at the revision recorded by the main repository.

### 10.3 Undefined Symbols or Missing Core Implementations

Check the build tags.

A missing `xray`, `sing`, or `hysteria2` tag may exclude the corresponding implementation.

Compare the actual build command with the intended deployment configuration.

### 10.4 Linker Version Injection Fails

The current README uses:

```text
github.com/proxynav/V2bX/cmd.version
```

This is based on the module path and version variable in the current source.

If the linker reports that the variable cannot be found or the version is not injected, verify:

1. The package path declared by `go.mod`.
2. The actual variable name and type in `cmd/`.
3. Whether the variable is a string.
4. Whether the linker target uses the correct fully qualified package path.

Do not replace the module path with the official GitHub repository URL without a coordinated module migration.

### 10.5 CGO or Native Compilation Errors

Check:

```bash
go env CGO_ENABLED CC CXX
```

Install the appropriate compiler and platform development dependencies if required.

Do not blindly set `CGO_ENABLED=0` to work around native compilation errors.

### 10.6 Docker Build Fails

Check:

- Whether the submodule was initialized.
- Whether the Docker build context includes all required source files.
- Whether the selected Go builder image is compatible.
- Whether module downloads can reach the required repositories.
- Whether the Docker build uses the expected build tags.

Use plain build output for more detailed diagnostics:

```bash
docker build --progress=plain -t proxynav/v2bx:debug .
```

## 11. Build Checklist

Before publishing a binary or submitting a build-related change:

- [ ] The source comes from the official ProxyNav V2bX repository.
- [ ] Git submodules are initialized and checked out at the expected revisions.
- [ ] The Go toolchain is compatible with `go.mod`.
- [ ] `GOEXPERIMENT=jsonv2` is enabled.
- [ ] Required proxy core tags are selected.
- [ ] Required optional feature tags are selected.
- [ ] Dependencies are downloaded and verified.
- [ ] Relevant tests and static analysis pass.
- [ ] The target OS and architecture are correct.
- [ ] CGO requirements have been checked.
- [ ] Version injection and linker flags are correct.
- [ ] The resulting binary has been inspected and tested.

## 12. References

- [ProxyNav V2bX Repository](https://github.com/proxynav/V2bX)
- [Original V2bX Source](https://github.com/wyx2685/V2bX)
- [Go Documentation](https://go.dev/doc/)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Go Build Command](https://pkg.go.dev/cmd/go)
- [Docker Build Documentation](https://docs.docker.com/build/)