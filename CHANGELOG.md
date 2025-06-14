# Changelog

## [0.1.10] - 2025-06-15

## Fixed

App crashing when updating secret and pressing S key

## [0.1.9] - 2024-04-28

## Added

-- Vault cache token lookup

## [0.1.8] - 2024-04-28

## Fixed

-- Fixing issue with popups not being focused and requiring selection with mouse

## Added

-- Adding metadata view on secret objects

## [0.1.7] - 2024-04-24

## Added

-- Fallback method for mounts listing when user doesnt access to `sys/mounts`

## [0.1.5] - 2024-04-18

## Fixed

- Fixing issue where vaulty will error if no config file is provided.

## [0.1.4] - 2024-04-17

## Fixed

- Fixing issue where cli wont run if version check fails.

## [0.1.3] - 2024-03-14

## Fixed

- Dynamic version passed when building the binary.

## [0.1.2] - 2024-03-13

## Added

- Allows for custom config file to be passed over during using `-c` rather than using the default one

## [0.1.1] - 2024-01-24

## Added

- Additional error message when failing to create vault client

## Fixed

- Fixed loading for client key when using VAULT_CLIENT_KEY

## [0.1.0] - 2024-01-23

## Added
- Env variable loading in addition to a yaml 
- Namespace support for enterprise vault instances

## Fixed
- Minor bugfixes around navigation

## Changes
- Housekeeping change

## [0.0.7] - 2023-12-03

## Added
- Creation of new secrets and paths

## Fixed
- Formatting and layout for different views around secrets when editing and displaying json

## Changes
- Commands layout has `<` and `>` removed to improve readability


## [0.0.6] - 2023-12-03

## Added
- Support for both PUT and PATCH for KV2 secrets
    - Had to modify the default methods in the vault package.. I couldn't figure out a clean way to get rid of the wrapper
- Better key mappings
- Additional information pane to show edit mode and filters used to search

## Fixed
- Correctly scrolling to the top on secrets and policy view

## Changes
- Refactoring and restructuring to make navigation in the repo easier

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
