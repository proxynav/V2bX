# Release Process

This document describes the release process for V2bX, including release preparation, validation, tagging, publishing, and post-release verification.

## 1. Overview

Releases should be prepared from a stable, reviewed state of the source code. All release-related changes must be validated before publishing a new version.

The release process aims to ensure:

- Reproducible and consistent builds.
- Proper validation of source code and dependencies.
- Clear versioning and release notes.
- Traceability between published binaries and source code.
- Reliable release artifacts for supported platforms and architectures.

## 2. Release Preparation

Before creating a release, ensure that the intended changes are ready.

### 2.1 Review Changes

- Review all changes intended for the release.
- Ensure pending fixes and features have been merged into the release branch.
- Confirm there are no known critical bugs or unresolved regressions.
- Verify that configuration examples and documentation are up to date.

### 2.2 Check Dependencies

Review dependency changes, particularly those involving:

- sing-box
- Hysteria2
- Go modules
- Build tags and compile-time options

Ensure dependency updates are intentional and compatible with the supported build configuration.

### 2.3 Confirm Version Information

Determine the release version before starting the release build.

Use semantic versioning where applicable:

- `MAJOR`: incompatible changes.
- `MINOR`: backward-compatible functionality.
- `PATCH`: backward-compatible bug fixes.

Release versions should use the `v` prefix, for example:

- `v1.2.0`
- `v1.2.1`
- `v1.3.0`

Pre-release versions may use suffixes such as `-beta.1` or `-rc.1`.

## 3. Validation and Testing

Run the relevant tests and checks before building release artifacts.

### 3.1 Run Tests

```bash
go test ./...
```

Review the results and resolve any unexpected failures before proceeding.

### 3.2 Verify the Build

V2bX uses build tags to control core and feature support. Ensure the release build uses the intended build tags and Go environment.

Example full-featured build:

```bash
GOEXPERIMENT=jsonv2 go build -v \
  -trimpath \
  -tags "sing xray hysteria2 with_quic with_grpc with_utls with_wireguard with_acme with_gvisor" \
  -ldflags "-s -w" \
  -o V2bX .
```

Confirm that the binary is generated successfully and that its version information is correct.

### 3.3 Smoke Testing

Where possible, test the release binary in a clean environment.

Recommended checks:

- Binary starts successfully.
- Configuration parsing and validation work as expected.
- Supported cores initialize successfully.
- Panel API connectivity works.
- Node configuration can be loaded.
- User synchronization works.
- Traffic reporting works.
- Supported protocols can establish connections.
- Service restart and shutdown behave correctly.

## 4. Create the Release Tag

Once validation is complete, create a release tag referencing the exact commit intended for publication.

Check the working tree first:

```bash
git status
git log -1 --oneline
```

Create and push the tag:

```bash
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin vX.Y.Z
```

Replace `vX.Y.Z` with the actual release version.

Release tags should not be moved or reused. If a release has a problem after publication, publish a new version instead of replacing the existing tag.

## 5. Build Release Artifacts

Release binaries should be built from the tagged commit rather than from an uncommitted working tree.

Check out the release tag:

```bash
git checkout vX.Y.Z
```

Build the release binary using the approved build configuration.

For multi-platform releases, build each supported operating system and architecture using the project's release build matrix.

Verify the resulting artifacts:

- Correct target OS and architecture.
- Correct core and feature support.
- Successful binary execution.
- Correct version information.
- No unintended debug or temporary files.

## 6. Publish the GitHub Release

Create a GitHub Release corresponding to the pushed tag.

The release should include:

### Release Title

Use the version number as the release title:

`V2bX vX.Y.Z`

### Release Notes

Release notes should summarize:

- New features.
- Bug fixes.
- Core and dependency updates.
- Breaking changes.
- Configuration changes.
- Known issues, if any.

### Release Assets

Upload the required release binaries and supporting files.

If applicable, include checksums for verifying downloaded artifacts.

## 7. Post-Release Verification

After publication:

- Verify that the release page is accessible.
- Confirm all expected assets are uploaded.
- Verify checksums, if provided.
- Test downloading the published artifacts.
- Confirm the release tag points to the intended commit.
- Perform a basic installation or upgrade test using the published release.

If a critical issue is discovered, document it and prepare a corrective release.

## 8. Release Checklist

Use this checklist before publishing a release:

- [ ] Changes reviewed and merged.
- [ ] Dependencies reviewed.
- [ ] Version determined.
- [ ] Tests passed.
- [ ] Release build completed successfully.
- [ ] Smoke tests completed.
- [ ] Release tag created and pushed.
- [ ] Artifacts built from the tagged commit.
- [ ] Release notes prepared.
- [ ] GitHub Release published.
- [ ] Assets verified.
- [ ] Post-release checks completed.

## 9. Release Integrity

Every published release should be traceable to a specific source commit and tag.

Do not modify published artifacts silently or reuse an existing release tag. Any correction to published content should be documented and, when necessary, published as a new release version.
