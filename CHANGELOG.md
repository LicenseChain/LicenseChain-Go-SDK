# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.2] - 2026-10-10

- Bumped `github.com/MicahParks/keyfunc/v3` to 3.8.2; module `go` directive is `1.25.0`.
- CI matrix updated to Go 1.25 / 1.26 / 1.27; goreleaser job uses Go 1.27.
- GitHub Actions bumped to `actions/checkout@v7`, `actions/setup-go@v7`, `goreleaser/goreleaser-action@v7`.

## [1.0.1] - 2026-10-10

- Relicensed public surface under Elastic License 2.0 (ELv2).
- Raised Go module requirement from 1.21 to 1.23 (later raised again in 1.0.2).
- Added Dependabot for Go modules and GitHub Actions.

## [1.0.0] - 2026-04-06

- Initial public SDK surface with repository documentation kit.
