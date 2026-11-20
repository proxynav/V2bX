# Development Environment Setup

This guide explains how to prepare a local development environment for **ProxyNav V2bX**, including the required Go toolchain, Git repository setup, dependencies, build configuration, and recommended development tools.

The instructions are based on the current project source and build configuration.

## 1. Project Repository

The official project repository is:

- **ProxyNav V2bX:** https://github.com/proxynav/V2bX
- **Original fork source:** https://github.com/wyx2685/V2bX

All new development, documentation, issues, and pull requests should use the official ProxyNav repository.

The project is written in Go and provides a multi-core proxy node management implementation based on Xray and sing-box.

## 2. Prerequisites

Install the following tools before building the project.

### 2.1 Go

The current `go.mod` specifies:

- Go language version: `1.25`
- Toolchain: `go1.25.0`

Install Go 1.25 or a compatible newer toolchain.

Verify the installation:

```bash
go version
go env GOROOT GOPATH GOTOOLCHAIN
```

The project build instructions also use the `GOEXPERIMENT=jsonv2` environment variable. Keep this setting when following the documented build procedure.

### 2.2 Git

Git is required to clone the repository and initialize its submodules.

Verify the installation:

```bash
git --version
```

### 2.3 Optional Development Tools

The following tools are recommended:

| Tool               | Purpose                                              |
| ------------------ | ---------------------------------------------------- |
| Visual Studio Code | Lightweight Go development                           |
| GoLand             | Integrated Go development environment                |
| Delve              | Go debugging                                         |
| Make               | Optional build automation                            |
| GCC or Clang       | Native compilation when required by CGO dependencies |

A C compiler is not automatically required for every build. Whether one is needed depends on the selected dependencies, target platform, and `CGO_ENABLED` setting.

## 3. Clone the Repository

Clone the official repository:

```bash
git clone https://github.com/proxynav/V2bX.git
cd V2bX
```

### 3.1 Initialize Git Submodules

The repository contains a Git submodule named `sing-box_mod`.

Clone the repository together with its submodules:

```bash
git clone --recurse-submodules https://github.com/proxynav/V2bX.git
cd V2bX
```

If the repository has already been cloned without submodules, initialize them manually:

```bash
git submodule update --init --recursive
```

Verify the submodule state:

```bash
git submodule status
```

Make sure the submodule is checked out at the commit recorded by the main repository.

Do not replace the submodule with an arbitrary upstream version. Its pinned revision may contain changes required by the project.

### 3.2 Verify the Working Tree

```bash
git remote -v
git branch --show-current
git status --short
git rev-parse HEAD
```

This helps confirm that you are working in the intended repository and on the expected branch.

## 4. Understand the Project Build Requirements

Before building V2bX, understand the following project-specific settings.

### 4.1 Go Module Configuration

The current `go.mod` declares:

```go
module github.com/InazumaV/V2bX

go 1.25

toolchain go1.25.0
```

The module path is inherited from the original project and is currently different from the official GitHub repository URL.

This distinction is intentional in the current development setup:

- Git clone URL: `github.com/proxynav/V2bX`
- Current Go module path: `github.com/InazumaV/V2bX`

Do not independently rewrite the module path or internal imports during ordinary development.

A module-path migration should be handled as a separate coordinated change.

### 4.2 JSON v2 Experiment

The project's documented build command enables:

```bash
GOEXPERIMENT=jsonv2
```

Use this setting when building the project with the documented build configuration.

For example:

```bash
export GOEXPERIMENT=jsonv2
```

Verify the current value:

```bash
go env GOEXPERIMENT
```

If the build fails with errors related to experimental JSON packages or APIs, first verify that the expected Go toolchain and experiment setting are being used.

### 4.3 Conditional Compilation

V2bX uses Go build tags to select proxy cores and optional features.

The documented build configuration includes the following tags:

| Tag              | Purpose                                  |
| ---------------- | ---------------------------------------- |
| `xray`           | Include Xray core integration            |
| `sing`           | Include sing-box core integration        |
| `hysteria2`      | Include Hysteria2-related implementation |
| `with_quic`      | Enable QUIC-related functionality        |
| `with_grpc`      | Enable gRPC-related functionality        |
| `with_utls`      | Enable uTLS-related functionality        |
| `with_wireguard` | Enable WireGuard-related functionality   |
| `with_acme`      | Enable ACME-related functionality        |
| `with_gvisor`    | Enable gVisor-related functionality      |

These tags determine which implementations and optional features are included in the binary.

A build that omits a required tag may compile successfully while lacking the functionality being investigated.

Use the same tags as the intended deployment when reproducing a build or runtime problem.

## 5. Download Go Dependencies

After cloning the repository and initializing its submodules, download the Go module dependencies:

```bash
go mod download
```

Verify downloaded dependencies:

```bash
go mod verify
```

Inspect the module configuration:

```bash
go list -m all
```

### 5.1 Network and Module Proxy Configuration

If dependency downloads fail, inspect the Go module environment:

```bash
go env GOPROXY GOSUMDB GOPRIVATE
```

The Go command uses module proxies and version control repositories to resolve dependencies.

If your network cannot access the default module proxy, configure an appropriate proxy for your environment.

For example:

```bash
go env -w GOPROXY=https://proxy.golang.org,direct
```

Use a trusted proxy and follow your organization's dependency security requirements.

### 5.2 Do Not Modify Dependencies Unnecessarily

Avoid running `go mod tidy` as part of routine environment setup.

It may change `go.mod` and `go.sum`, especially when build tags or platform-specific dependencies are involved.

Only update dependencies when the change is intentional. Review the resulting diff before committing it.

## 6. Recommended Project Structure

The repository contains the following major directories:

```text
V2bX/
├── .github/
├── api/
│   └── panel/
├── cmd/
├── common/
├── conf/
├── core/
├── example/
├── limiter/
├── node/
├── test_data/
├── sing-box_mod/
├── Dockerfile
├── go.mod
├── go.sum
└── main.go
```

The main development areas are:

| Directory       | Responsibility                       |
| --------------- | ------------------------------------ |
| `cmd/`          | Application commands and startup     |
| `conf/`         | Configuration structures and loading |
| `api/panel/`    | Panel API integration                |
| `core/`         | Proxy core integration               |
| `node/`         | Node configuration and lifecycle     |
| `limiter/`      | Connection and traffic limiting      |
| `common/`       | Shared utilities                     |
| `example/`      | Example configurations               |
| `test_data/`    | Test fixtures                        |
| `sing-box_mod/` | Git submodule used by the project    |

When modifying a component, inspect the surrounding package and its dependencies before introducing new abstractions.

## 7. IDE Setup

### 7.1 Visual Studio Code

Install Visual Studio Code and the official Go extension.

Recommended extensions:

- Go
- EditorConfig support, if the repository adopts EditorConfig
- GitLens, optionally

Open the repository root directory:

```bash
code .
```

Make sure the IDE uses the same Go toolchain as the terminal.

If the IDE does not recognize project dependencies, check the selected Go SDK and module configuration.

### 7.2 GoLand

GoLand can open the repository directly.

Recommended settings:

- Use Go modules for dependency management.
- Select Go 1.25 or a compatible newer toolchain.
- Enable automatic `gofmt` formatting.
- Configure the project root as the working directory.
- Ensure the terminal and IDE use the same environment variables.

### 7.3 Environment Variables in IDEs

The build environment should include:

```bash
GOEXPERIMENT=jsonv2
```

If you use a custom build configuration, ensure that the required build tags are also configured in the IDE's Go build settings.

## 8. Initial Build Verification

After preparing the environment, perform a basic verification.

### 8.1 Check Go Packages

```bash
go list ./...
```

This checks whether Go can discover the packages in the current module for the selected build environment.

### 8.2 Run Tests

```bash
GOEXPERIMENT=jsonv2 go test ./...
```

For verbose output:

```bash
GOEXPERIMENT=jsonv2 go test -v ./...
```

If a package has no test files, Go may report `[no test files]`. This is not a test failure.

### 8.3 Build the Application

Use the documented project build configuration:

```bash
mkdir -p build_assets

GOEXPERIMENT=jsonv2 go build -v \
  -o build_assets/V2bX \
  -tags "sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor" \
  -trimpath \
  -ldflags "-X 'github.com/InazumaV/V2bX/cmd.version=$version' -s -w -buildid="
```

Set `version` before running the command if you want to inject a version string:

```bash
version="dev"
```

For interactive debugging, omit the symbol-stripping flags `-s -w` so that debugging information remains available.

The command above follows the current README build configuration. If build tags, version variables, or linker settings change in the repository, update this guide accordingly.

### 8.4 Verify the Binary

```bash
file build_assets/V2bX
go version -m build_assets/V2bX
```

These commands help verify the resulting binary's architecture and embedded Go build information.

## 9. Cross-Platform Development

Go supports multiple target operating systems and architectures, but successful cross-compilation depends on the selected dependencies and build tags.

Check the current environment:

```bash
go env GOOS GOARCH CGO_ENABLED
```

For cross-compilation, explicitly set the target:

```bash
GOOS=linux GOARCH=amd64 GOEXPERIMENT=jsonv2 go build
```

Other possible targets include:

- Linux ARM64
- Linux ARM
- Windows AMD64
- macOS ARM64
- macOS AMD64

Before adding a target to release automation, verify that all selected cores and optional features support that target.

Some dependencies may require CGO or platform-specific toolchains. Do not assume that setting `CGO_ENABLED=0` is valid for every build variant.

## 10. Useful Development Commands

| Task                   | Command                                |
| ---------------------- | -------------------------------------- |
| Check Go version       | `go version`                           |
| Inspect Go environment | `go env`                               |
| Download dependencies  | `go mod download`                      |
| Verify dependencies    | `go mod verify`                        |
| List packages          | `go list ./...`                        |
| Format source          | `gofmt -w .`                           |
| Check formatting       | `gofmt -l .`                           |
| Run tests              | `GOEXPERIMENT=jsonv2 go test ./...`    |
| Run verbose tests      | `GOEXPERIMENT=jsonv2 go test -v ./...` |
| Run static analysis    | `go vet ./...`                         |
| Build application      | `GOEXPERIMENT=jsonv2 go build`         |
| Inspect binary         | `go version -m ./build_assets/V2bX`    |
| Check submodules       | `git submodule status`                 |

## 11. Troubleshooting Environment Problems

### 11.1 Go Version Mismatch

Symptoms:

- Build errors related to unsupported language features.
- Toolchain download failures.
- Errors involving experimental packages or APIs.

Check:

```bash
go version
go env GOTOOLCHAIN GOEXPERIMENT
```

Use a compatible Go 1.25 toolchain and enable the required experiment setting.

### 11.2 Missing Submodule Content

Symptoms:

- Missing package or source files under `sing-box_mod/`.
- Build errors referencing code expected from the submodule.

Fix:

```bash
git submodule update --init --recursive
git submodule status
```

### 11.3 Dependency Download Failures

Check:

```bash
go env GOPROXY GOSUMDB GOPRIVATE
```

Verify network access to the configured module proxy and the required source repositories.

### 11.4 Missing Core Implementation

Symptoms:

- Undefined symbols associated with a specific core.
- Missing functionality in a compiled binary.
- Build differences between local and CI environments.

Check the build tags and compare them with the documented build command.

### 11.5 CGO or Native Compiler Errors

Check:

```bash
go env CGO_ENABLED CC CXX
```

If the selected dependency requires native compilation, install the appropriate compiler and platform development libraries.

Do not disable CGO without first verifying that the selected dependencies support it.

### 11.6 Unexpected Module Path Errors

The current module path is inherited from the original repository:

```text
github.com/InazumaV/V2bX
```

Do not automatically replace it with:

```text
github.com/proxynav/V2bX
```

The official GitHub repository URL and the Go module path are separate concerns. A coordinated module migration is required to change the latter safely.

## 12. Development Resources

- [ProxyNav V2bX Repository](https://github.com/proxynav/V2bX)
- [Original V2bX Source](https://github.com/wyx2685/V2bX)
- [Go Documentation](https://go.dev/doc/)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Go Command Documentation](https://pkg.go.dev/cmd/go)
- [Delve Debugger](https://github.com/go-delve/delve)
- [Visual Studio Code Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)

---

**Development environment checklist**

- [ ] Go 1.25 or a compatible newer toolchain is installed.
- [ ] Git is installed.
- [ ] The official ProxyNav V2bX repository has been cloned.
- [ ] Git submodules have been initialized.
- [ ] Go dependencies have been downloaded and verified.
- [ ] `GOEXPERIMENT=jsonv2` is configured for project builds.
- [ ] Required proxy core and optional feature build tags are understood.
- [ ] The project packages can be listed successfully.
- [ ] Tests and the initial build have been verified.
