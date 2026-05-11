# Changelog

## [0.2.0] - 2026-05-11

### Added

- **Auth Methods view** — `ctrl-a` from any view opens a table of all enabled auth methods showing mount path, type, and description
- **Enriched header** — startup now fetches token TTL, attached policies, seal status, and cluster name; all are shown in the info bar alongside vault address and version
- **Secret version management** (KV v2 only):
  - `D` — soft-delete the current version (recoverable via rollback)
  - `R` — rollback to a previous version by entering a version number
  - `X` — permanently destroy a specific version (irreversible, requires confirmation)
- **Secret metadata editing** (`M`) — edit `max_versions`, `cas_required`, and `delete_version_after` for KV v2 secrets directly in the TUI
- **Policy editing** (`E` in policy ACL view) — edit policy rules in-place using the built-in text editor; save with `ctrl-w`, cancel with `esc`
- **Format selector for secret create/update** — pressing `ctrl-n` (new secret), `P` (patch), or `U` (update) now prompts for JSON or Key-Value format before opening the editor:
  - **JSON** — opens the existing JSON text editor (unchanged flow)
  - **Key-Value** — opens a form with labelled Key/Value fields; an **Add Key** button appends additional pairs; pairs with blank keys are skipped on save; form pre-populates with existing data on `P`/`U`
- **Policy ACL vim-style navigation** — when reading a policy file, `j`/`k` scroll one line, `d`/`u` scroll a half-page, `g`/`G` jump to top/bottom (arrow keys and Page Up/Down continue to work as before)

### Fixed

- **Blank secret view after rollback or metadata toggle** — `SecretObject()` now resets `ShowMetadata`, `ShowJson`, and `Editable` flags on every entry; previously these stale flags caused `Render()` to take the in-place refresh path on a cleared body slot, producing a blank view for all subsequent secrets
- **Event loop blocking on vault operations** — delete version, destroy versions, rollback, and metadata update calls are now dispatched to background goroutines with `QueueUpdateDraw` for UI callbacks; previously they ran on the event goroutine and froze the UI during the network call
- **DoneFunc state corruption after rollback/destroy prompts** — `promptRollback` and `promptDestroyVersion` now save and restore `TextInfoInput.Props.DoneFunc`; previously they permanently overrode it, breaking all subsequent text-input flows
- **Non-blocking `Draw()`** — `Draw()` now uses `select/default` so it never blocks when the draw channel is already signalled
- **Error propagation in vault KV client** — `Get()` and `GetMetadata()` now return the error to the caller instead of logging it and returning nil; callers can now surface fetch failures correctly

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
