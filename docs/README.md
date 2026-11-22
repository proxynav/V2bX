# V2bX Documentation

Welcome to the V2bX documentation.

This documentation provides information about installing, configuring, operating, and developing V2bX. It is intended for users, developers, and project maintainers.

## Documentation Structure

The documentation is organized into the following sections.

### Getting Started

Guides for installing and getting started with V2bX.

- [Installation](getting-started/installation.md) — Installation instructions and system requirements.
- [Quick Start](getting-started/quick-start.md) — Get V2bX up and running quickly.
- [First Configuration](getting-started/first-configuration.md) — Create your first V2bX configuration.

### Architecture

Technical documentation describing the internal architecture and design of V2bX.

- [Architecture Overview](architecture/overview.md) — High-level system architecture and major components.
- [Project Structure](architecture/project-structure.md) — Directory layout and package responsibilities.
- [Application Lifecycle](architecture/application-lifecycle.md) — Application startup, initialization, runtime, and shutdown.
- [Configuration System](architecture/configuration-system.md) — Configuration loading, parsing, validation, and updates.
- [Core Abstraction](architecture/core-abstraction.md) — Core interfaces and implementation design.
- [Data Flow](architecture/data-flow.md) — How configuration and node data flow through the application.

### Development

Guides and references for contributors and developers.

- [Development Environment](development/environment.md) — Set up a local Go development environment.
- [Build Guide](development/build.md) — Build V2bX from source.
- [Testing Guide](development/testing.md) — Run tests and measure code coverage.
- [Debugging Guide](development/debugging.md) — Debug V2bX during development.
- [Coding Style](development/coding-style.md) — Code conventions and development practices.
- [Dependency Management](development/dependency-management.md) — Manage Go modules and project dependencies.

### Proxy Cores

Documentation about supported proxy cores and how to extend core integrations.

- [Core Overview](cores/overview.md) — Supported cores and their responsibilities.
- [sing-box Core](cores/sing-box.md) — sing-box integration and implementation details.
- [Adding a New Core](cores/adding-a-new-core.md) — Guide to implementing a new proxy core.

### Protocols

Protocol-specific configuration and implementation documentation.

- [Protocol Overview](protocols/overview.md) — Supported protocols and their general architecture.
- [VMess](protocols/vmess.md)
- [VLESS](protocols/vless.md)
- [Trojan](protocols/trojan.md)
- [Shadowsocks](protocols/shadowsocks.md)
- [Hysteria2](protocols/hysteria2.md)
- [AnyTLS](protocols/anytls.md)

Each protocol guide may cover configuration options, core compatibility, implementation details, and usage examples.

### Operations and Troubleshooting

Guides for operating, monitoring, and troubleshooting V2bX.

- [Configuration Reference](operations/configuration-reference.md) — Complete configuration options and descriptions.
- [Logging](operations/logging.md) — Logging configuration and diagnostics.
- [Troubleshooting](operations/troubleshooting.md) — Common issues and troubleshooting procedures.
- [Performance](operations/performance.md) — Performance considerations and optimization.

### Architecture Decision Records

Records of important technical decisions and their rationale.

- [Architecture Decisions](decisions/README.md) — Index of architecture decision records.

ADRs document the context, alternatives, decisions, and consequences of significant architectural changes.

## Contributing to Documentation

Documentation improvements are welcome.

When contributing documentation, please follow these guidelines:

1. Write documentation in clear and concise English.
2. Use Markdown formatting consistently.
3. Keep examples synchronized with the actual implementation.
4. Include practical examples where appropriate.
5. Update this index when adding new documentation.
6. Avoid documenting behavior that is not currently implemented.
7. Review internal links before submitting a pull request.

For code contribution guidelines, see [CONTRIBUTING.md](../CONTRIBUTING.md).

## Project Resources

- [GitHub Repository](https://github.com/proxynav/V2bX)
- [README](../README.md)
- [Contributing Guidelines](../CONTRIBUTING.md)
- [Issue Tracker](https://github.com/proxynav/V2bX/issues)
- [Go Documentation](https://go.dev/doc/)

---

*This documentation is maintained alongside the V2bX source code. Some sections may be incomplete or under development.*
