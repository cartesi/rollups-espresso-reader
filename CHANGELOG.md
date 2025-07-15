<!-- markdownlint-disable MD024 -->
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0]

### Added

- Support for Node v2's snapshot feature

### Changed

- Update Node v2 dependency to `v2.0.0-alpha.6`
- Update Espresso SDK version to `20250623`
- Special handling to avoid issues when reading blockchain data from certain free tier services like Alchemy's
- L1 input processing now always starting from initial Espresso block's L1 finalized block
- Fixed CORS issue for Espresso Reader endpoints

## [0.3.0]

### Added

- Infrastructure and instructions to execute Espresso Reader + Node locally
- Database transaction atomicity, ensuring saved state is always consistent and recoverable

### Changed

- Update Node v2 dependency to `v2.0.0-alpha.4` (DataAvailability configuration now specified during app deployment)
- Update Espresso SDK to comply with API `/v1` upgrade

## [0.2.3-node-20250128]

### Changed

- Improved epoch closure procedure (avoids edge case that can lead to skipping L1 blocks)
- Added CI e2e tests with local environment using espresso-dev-node

## [0.2.2-node-20250128]

### Added

- CI binaries generation for multiple platforms

### Changed

- Fixed epoch closure procedure
- Fixed startup policy when there is no record of the last processed espresso block
- Default URI endpoint updated to `0.0.0.0:8080`

## [0.2.1-node-20250128]

### Changed

- Binary generation in CI now supporting multiple platforms

## [0.2.0-node-20250128]

### Changed

- Adapt to new node version (2025-01-28)
- Safer parsing of EIP-712 structures

## [0.1.0]

### Added

- Initial EspressoReader implementation


<!-- markdownlint-disable MD053 -->
[Unreleased]: https://github.com/cartesi/rollups-espresso-reader/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/cartesi/rollups-espresso-reader/releases/tag/v0.4.0
[0.3.0]: https://github.com/cartesi/rollups-espresso-reader/releases/tag/v0.3.0
[0.2.3-node-20250128]: https://github.com/cartesi/rollups-espresso-reader/releases/tag/v0.2.3-node-20250128
[0.2.2-node-20250128]: https://github.com/cartesi/rollups-espresso-reader/releases/tag/v0.2.2-node-20250128
[0.2.1-node-20250128]: https://github.com/cartesi/rollups-espresso-reader/releases/tag/v0.2.1-node-20250128
[0.2.0-node-20250128]: https://github.com/cartesi/rollups-espresso-reader/releases/tag/v0.2.0-node-20250128
[0.1.0]: https://github.com/cartesi/rollups-espresso-reader/releases/tag/v0.1.0
