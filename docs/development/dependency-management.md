# Dependency Management

This document defines how dependencies are introduced, reviewed, updated, and maintained in **ProxyNav V2bX**.

The project uses more than one dependency mechanism:

- **Go Modules** for Go libraries and packages, tracked by `go.mod` and `go.sum`.
- **Git submodules** for source code maintained outside the main repository, including `sing-box_mod`.
- **Build toolchain and container images** for the compiler and runtime environment used to produce release artifacts.

These mechanisms have different versioning and update procedures. They must be reviewed and tested independently.

## 1. Repository and Dependency Baseline

- Official repository: <https://github.com/proxynav/V2bX>
- Original upstream fork: <https://github.com/wyx2685/V2bX>
- Go module manifest: `go.mod`
- Go module checksums: `go.sum`
- Git submodule configuration: `.gitmodules`

The repository currently declares:

```go
module github.com/proxynav/V2bX

go 1.25

toolchain go1.25.0
```

**Important:** The official Git repository URL and the Go module path are currently different. The existing module path is still used by internal imports and linker flags. Do not change it as part of an ordinary dependency update.

A module path migration is a separate project-wide change and must update all affected imports, build configuration, documentation, and release workflows together.

## 2. Dependency Categories

### 2.1 Direct Go dependencies

Direct dependencies are listed in the main `require` block of `go.mod`. They provide functionality directly used by the application.

The current dependency set includes important components such as:

| Dependency                              | Responsibility                       |
| --------------------------------------- | ------------------------------------ |
| `github.com/apernet/hysteria/core/v2`   | Hysteria core functionality          |
| `github.com/apernet/hysteria/extras/v2` | Hysteria supporting functionality    |
| `github.com/sagernet/sing`              | sing-related functionality           |
| `github.com/sagernet/sing-box`          | sing-box core                        |
| `github.com/xtls/xray-core`             | Xray core                            |
| `github.com/go-acme/lego/v4`            | ACME certificate automation          |
| `github.com/go-resty/resty/v2`          | HTTP client functionality            |
| `github.com/spf13/cobra`                | CLI framework                        |
| `github.com/spf13/viper`                | Configuration management             |
| `github.com/sirupsen/logrus`            | Logging                              |
| `go.uber.org/zap`                       | Structured logging                   |
| `golang.org/x/crypto`                   | Extended cryptographic functionality |
| `golang.org/x/sys`                      | Operating-system interfaces          |
| `google.golang.org/protobuf`            | Protocol Buffers                     |

This table is illustrative rather than an exhaustive dependency inventory. Always use the current `go.mod` as the authoritative source.

### 2.2 Indirect Go dependencies

Indirect dependencies are dependencies required by direct dependencies or other indirect dependencies.

They are normally marked with `// indirect` in `go.mod`.

Examples in the current dependency graph include:

- Cloud DNS provider SDKs used by ACME integrations.
- QUIC and networking libraries.
- sing-box supporting modules.
- Serialization and protocol libraries.
- Platform-specific system and networking packages.

Do not manually promote an indirect dependency to a direct dependency unless application source code imports it directly or there is another documented reason.

Likewise, do not manually remove indirect dependencies simply because they are not imported by application packages. They may be required by the selected versions of other modules.

### 2.3 Git submodules

The repository currently declares a Git submodule:

```ini
[submodule "sing-box_mod"]
    path = sing-box_mod
    url = https://github.com/wyx2685/sing-box_mod.git
```

Unlike Go modules, a Git submodule is pinned by a Git commit recorded in the parent repository.

The submodule is part of the source tree used by the build and must be initialized before compiling the project.

### 2.4 Toolchain and container dependencies

The project also relies on external build and runtime components:

- Go 1.25 toolchain.
- Alpine Linux base image in the current Dockerfile.
- `tzdata` and `ca-certificates` packages in the runtime image.
- GitHub Actions and other CI tooling, when used by repository workflows.

Changes to these components should be reviewed as supply-chain changes, not treated as ordinary Go library upgrades.

## 3. Initializing Dependencies

Clone the repository with its submodules:

```bash
git clone --recurse-submodules \
  https://github.com/proxynav/V2bX.git

cd V2bX
```

If the repository has already been cloned without submodules:

```bash
git submodule update --init --recursive
```

Verify the submodule state:

```bash
git submodule status --recursive
```

Download Go modules:

```bash
GOEXPERIMENT=jsonv2 go mod download
```

Verify module checksums:

```bash
go mod verify
```

The `GOEXPERIMENT=jsonv2` setting is included here to match the current project build configuration. Use the same environment setting when running project validation commands that need to reproduce the documented build environment.

## 4. Adding a New Go Dependency

Before introducing a dependency, maintainers should answer the following questions:

1. Is the functionality already available in the Go standard library or an existing project dependency?
2. Is the package actively maintained?
3. Is its license compatible with the project?
4. Does it introduce additional transitive dependencies?
5. Does it introduce platform-specific, CGO, or native-library requirements?
6. Does it affect the size, performance, security, or compatibility of the final binary?
7. Is it compatible with the project's Go version and build tags?

Prefer a focused, well-maintained dependency over a large framework when both can solve the problem.

### Add a dependency

First, import the package in application code. Then add the corresponding module:

```bash
go get example.com/vendor/package@v1.2.3
```

Replace the example module path and version with the actual dependency.

Use an explicit version instead of relying on an unspecified latest version.

Review the resulting changes:

```bash
git diff -- go.mod go.sum
```

Then normalize the module files:

```bash
GOEXPERIMENT=jsonv2 go mod tidy
```

Review the complete changes again:

```bash
git diff -- go.mod go.sum
```

Do not accept unexpected dependency removals, version changes, or large transitive dependency changes without understanding their cause.

## 5. Updating Dependencies

Dependency updates should be deliberate and reviewable.

### 5.1 Inspect available updates

List all modules and their selected versions:

```bash
go list -m all
```

Check which dependencies have newer versions available:

```bash
go list -m -u all
```

Inspect why a particular dependency is present:

```bash
go mod why -m example.com/vendor/package
```

Inspect the dependency graph when investigating version selection or unexpected transitive dependencies:

```bash
go mod graph
```

### 5.2 Update one dependency

For a targeted update:

```bash
go get example.com/vendor/package@v1.2.4
```

This is the preferred approach for changes to critical components because it makes the resulting diff easier to understand and troubleshoot.

### 5.3 Update all dependencies

A broad update can be performed with:

```bash
go get -u ./...
```

This may update many direct and indirect dependencies and can introduce breaking changes.

Use it only in a dedicated maintenance change. Do not combine a broad dependency refresh with unrelated feature development.

After updating:

```bash
GOEXPERIMENT=jsonv2 go mod tidy
go mod verify
```

Review all changes in `go.mod` and `go.sum` before proceeding.

### 5.4 Upgrade policy

Use the following general policy:

| Dependency type          | Recommended approach                            |
| ------------------------ | ----------------------------------------------- |
| Security fixes           | Prioritize and evaluate promptly                |
| Patch updates            | Prefer when they contain relevant fixes         |
| Minor updates            | Review changelog and compatibility              |
| Major updates            | Require explicit compatibility review           |
| Pre-release versions     | Require a clear technical justification         |
| Forked or custom modules | Review upstream changes and local compatibility |

Do not assume that semantic version numbers alone guarantee compatibility with this project's build tags, protocol implementations, or runtime behavior.

## 6. Critical Dependencies

Some dependencies have a much greater impact on application behavior than ordinary utility libraries.

Changes to these components require additional testing.

### 6.1 Xray Core

Dependency:

```text
github.com/xtls/xray-core
```

Before upgrading Xray Core, verify:

- Existing node protocol implementations still compile.
- Configuration structures and API usage remain compatible.
- TLS, Reality, transport, and other supported features work as expected.
- Existing connection and traffic-limit behavior is unaffected.
- The Xray-enabled build succeeds.

### 6.2 sing-box

Dependency:

```text
github.com/sagernet/sing-box
```

Before upgrading sing-box, verify:

- Core initialization and lifecycle behavior.
- Inbound and outbound configuration compatibility.
- Protocol registration and supported transports.
- TLS, ECH, QUIC, and other enabled features.
- Compatibility with the project's custom sing-box modifications.

The project also contains the `sing-box_mod` Git submodule. Changes to the Go module version and changes to the submodule revision must be reviewed together when they affect the same integration.

### 6.3 Hysteria

Dependencies:

```text
github.com/apernet/hysteria/core/v2
github.com/apernet/hysteria/extras/v2
```

Before upgrading Hysteria, verify:

- Hysteria core initialization.
- Authentication and user management.
- Transport and congestion-control behavior.
- QUIC-related build and runtime behavior.
- Compatibility with the project's Hysteria node configuration.

### 6.4 ACME and certificate dependencies

Dependency:

```text
github.com/go-acme/lego/v4
```

Before upgrading ACME-related modules, verify:

- Certificate issuance.
- DNS provider integrations.
- HTTP-01 and DNS-01 challenge behavior, where configured.
- Certificate renewal.
- TLS certificate loading and storage.

ACME integrations may pull in a large number of cloud-provider SDKs. Avoid introducing additional DNS provider dependencies unless the functionality is actually needed.

### 6.5 Logging and configuration libraries

The current dependency set includes both Logrus and Zap, as well as Cobra and Viper.

When changing these libraries:

- Follow the logging API already used by the surrounding component.
- Avoid introducing a second logging abstraction without a clear migration plan.
- Preserve existing configuration defaults and parsing behavior.
- Test configuration reload and error handling where applicable.
- Avoid unrelated logging or configuration refactoring during dependency upgrades.

## 7. Managing the `sing-box_mod` Submodule

The `sing-box_mod` submodule is maintained separately from Go module dependencies.

### 7.1 Inspect the current revision

```bash
git submodule status --recursive
```

Enter the submodule:

```bash
cd sing-box_mod
git status
git log -5 --oneline
```

Inspect its configured remote:

```bash
git remote -v
```

### 7.2 Update the submodule

Do not update the submodule automatically as part of routine Go dependency maintenance.

When an update is required, first review the target repository and select a known compatible commit.

For example:

```bash
git submodule update --remote sing-box_mod
```

This command follows the submodule's configured remote-tracking branch. It should not be used blindly in release preparation.

Review the resulting commit change in the parent repository:

```bash
git diff --submodule=log
git status
```

The parent repository records the selected submodule commit. Commit the parent repository's updated gitlink together with any required integration changes.

### 7.3 Reproduce the recorded revision

For a clean checkout, initialize the submodule at the revision recorded by the parent repository:

```bash
git submodule update --init --recursive
```

Do not replace the recorded revision with the latest upstream commit just because it is newer.

### 7.4 Submodule update requirements

Every submodule update should include:

- The reason for the update.
- The old and new commit IDs.
- A summary of relevant upstream changes.
- Compatibility analysis with the current sing-box integration.
- Successful builds using the project's supported build configuration.
- Runtime test results for affected functionality.

If the submodule and the corresponding Go modules must be updated together, treat them as one coordinated change.

## 8. Dependency Replacement and Forks

Go supports `replace` directives in `go.mod`, which can redirect a module dependency to another module path, version, or local directory.

A replacement may be useful when:

- A required fix has not yet been released upstream.
- A project-specific fork is necessary.
- A temporary compatibility patch is required.

However, replacements can make builds less predictable if they depend on local paths or undocumented forks.

Before introducing a replacement:

1. Document the reason.
2. Identify the original dependency and replacement source.
3. Record the relevant upstream issue or commit.
4. Explain whether the replacement is temporary or permanent.
5. Ensure CI and release builds can resolve it.
6. Define how and when the replacement should be removed.

Avoid machine-specific local replacements such as:

```go
replace example.com/vendor/package => ../local-package
```

unless they are strictly limited to local development and are not committed into the shared module configuration.

Do not introduce a replacement merely to silence a compilation error. Investigate the underlying API or version compatibility issue first.

## 9. Dependency Validation

Every dependency change must pass the relevant validation steps before it can be merged.

### 9.1 Module integrity

```bash
go mod verify
```

### 9.2 Module consistency

```bash
GOEXPERIMENT=jsonv2 go mod tidy
git diff --exit-code -- go.mod go.sum
```

The second command checks whether running `go mod tidy` would introduce additional changes. If it reports a diff, review the result rather than automatically accepting it.

### 9.3 Build validation

Use the project's documented full build configuration:

```bash
version="dependency-check"
mkdir -p build_assets

GOEXPERIMENT=jsonv2 go build \
  -v \
  -o build_assets/V2bX \
  -tags "sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor" \
  -trimpath \
  -ldflags "-X 'github.com/proxynav/V2bX/cmd.version=$version' -s -w -buildid="
```

Also test any reduced build configurations affected by the dependency change.

Do not assume that a dependency update is safe just because one build-tag combination succeeds.

### 9.4 Tests and static checks

Run the available project tests:

```bash
GOEXPERIMENT=jsonv2 go test ./...
```

Run static analysis:

```bash
GOEXPERIMENT=jsonv2 go vet ./...
```

If the repository has protocol-specific or integration tests, run the relevant tests as well.

A successful compile does not prove that a proxy core, transport, certificate workflow, or protocol integration works correctly at runtime.

### 9.5 Docker validation

When a dependency change affects the release binary or runtime behavior, validate the Docker build too:

```bash
docker build -t proxynav/v2bx:dependency-check .
```

The current Dockerfile uses its own build command and build environment. Its linker flags are not identical to the README local release build command.

Do not assume that a successful local build automatically validates the Docker build.

## 10. Security and License Review

Before adding or upgrading a dependency, review:

- Project maintenance activity.
- Security advisories and known vulnerabilities.
- License compatibility.
- Dependency ownership and source repository.
- Release history and changelog.
- Transitive dependency changes.
- Whether the dependency downloads or executes additional code during builds.

For vulnerability inspection, maintainers may use the Go vulnerability database tooling:

```bash
govulncheck ./...
```

`govulncheck` is an external development tool, not currently a declared application dependency. Install and pin it through the project's development tooling policy before relying on it in CI.

A vulnerability report should be evaluated based on whether the affected code path is reachable and whether the vulnerable behavior is relevant to V2bX's deployment and runtime model.

Do not ignore a vulnerability solely because it is reported in an indirect dependency.

## 11. Dependency Update Pull Requests

Dependency updates should be submitted as focused pull requests.

A dependency update pull request should contain:

- The dependency name and old/new versions.
- The reason for the update.
- Relevant release notes or changelog links.
- Security advisory references, if applicable.
- Changes to transitive dependencies.
- Changes to `go.mod` and `go.sum`.
- Submodule revision changes, if applicable.
- Build and test results.
- Runtime validation results for affected functionality.

Avoid combining unrelated dependency updates with protocol features, configuration refactoring, or other functional changes.

For major upgrades to Xray, sing-box, Hysteria, or ACME components, consider separate pull requests or a clearly documented coordinated upgrade.

## 12. Rollback Procedure

If an upgrade causes build failures or runtime regressions, revert the dependency change to the last known-good state.

For Go module dependencies, restore the relevant `go.mod` and `go.sum` changes from the previous known-good commit, or explicitly select the previous version:

```bash
go get example.com/vendor/package@v1.2.3
GOEXPERIMENT=jsonv2 go mod tidy
```

For a Git submodule, restore the parent repository's previous gitlink:

```bash
git submodule update --init --recursive
```

If necessary, check out the known-good submodule revision recorded by the previous parent commit.

After rollback, repeat the build and test procedures. A successful rollback must restore both compilation and the affected runtime behavior.

## 13. Recommended Maintenance Practices

| Practice              | Recommendation                              |
| --------------------- | ------------------------------------------- |
| Dependency additions  | Require a clear technical reason            |
| Routine updates       | Prefer small, targeted updates              |
| Security updates      | Prioritize vulnerability remediation        |
| Core upgrades         | Require protocol and runtime validation     |
| Submodule updates     | Pin and review explicit commits             |
| `go.mod` and `go.sum` | Commit both files together                  |
| Local replacements    | Avoid committing machine-specific paths     |
| Broad upgrades        | Use dedicated maintenance changes           |
| Release builds        | Validate the documented build configuration |
| Rollbacks             | Preserve a known-good dependency state      |

### Automated dependency updates

GitHub Dependabot may be used to monitor Go modules, Docker base images, and GitHub Actions dependencies.

If enabled, automated update pull requests must still pass the same review and validation requirements described in this document.

Do not merge automated dependency updates solely because they were generated by Dependabot.

## 14. Maintainer Checklist

Before merging a dependency change:

- [ ] The dependency has a clear technical purpose.
- [ ] Its source, maintenance status, and license have been reviewed.
- [ ] Security implications have been considered.
- [ ] `go.mod` and `go.sum` contain only expected changes.
- [ ] Indirect dependency changes have been reviewed.
- [ ] Any `replace` directives are documented and justified.
- [ ] Submodule changes are pinned to reviewed commits.
- [ ] `go mod verify` succeeds.
- [ ] `go mod tidy` produces no unexplained changes.
- [ ] Relevant tests and static checks pass.
- [ ] The documented build configuration succeeds.
- [ ] Docker builds are validated when relevant.
- [ ] Runtime behavior has been tested for affected cores and protocols.
- [ ] Rollback information is available for high-impact upgrades.

## 15. References

- [ProxyNav V2bX Repository](https://github.com/proxynav/V2bX)
- [Original V2bX Repository](https://github.com/wyx2685/V2bX)
- [Go: Managing Dependencies](https://go.dev/doc/modules/managing-dependencies)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Go Vulnerability Management](https://go.dev/doc/security/vuln/)
