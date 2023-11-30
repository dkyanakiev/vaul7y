# Changelog

## [0.0.5] - 2023-11-30


### Fixed
- Missing commands for 2 views
- Version command check would fail if missing `VAULT_TOKEN` or `VAULT_ADDR` is missing

## [0.0.3] - 2023-11-30

### Added
- Job filtering on secrets and mount views
- Better navigation options between views
- `vaul7y -v` to check the version
- Added a check and error out to prevent vaul7y from freezing if vault token and address are not set

### Fixed
- Error and Info modals tabbing out and changing focus
- Enter key constantly moving you to the Secret Engines view. Its due to the way Unix system recognize Enter and Ctrl+M
- Fixed an issue with watcher causing conflicts 
- Fixed logger to discard messages and not brake rendering when debugging is not enabled
