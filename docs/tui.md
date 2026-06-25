# Interactive TUI

Run `vigil` on a terminal (or `vigil tui` explicitly) to open the interactive UI.
Everything here is also available as plain CLI commands — see `vigil --help`.

## Navigation

The home screen is a **menu**:

- `↑`/`↓` (or `k`/`j`) — move the cursor
- `Enter` (or `→`/`l`) — select
- `Esc` — go back to the menu
- `q` / `Ctrl-C` — quit

A header shows the scope + version, a live **status strip** (state · engine ·
override · next change), and an **update banner** when a newer release exists.

## Screens

| Menu item | What it does |
|-----------|--------------|
| **Status** | Full daemon status (state, engine, drivers, override, next transition, last error) plus a tail of the log. |
| **Override** | Force the session active or inactive, or resume the schedule. "Force active" has timed presets (indefinite / 1h / 2h / 4h) matching the CLI's `on --for`. |
| **Schedule** | Edit the weekly active windows (see keys below). |
| **Settings** | Engine, input interval/jitter/method/move-pixels, prevent-sleep, timezone, and Graph `client_id`/`tenant_id`. |
| **Service** | Install / start / stop / restart / uninstall the background service (runs the matching `vigil` command). |
| **Account** | Microsoft Graph sign-in / status / sign-out (for the `graph` engine). |
| **Check for updates** | Self-update to the latest release (skipped on dev builds / package-manager installs). |
| **Help** | Key reference. |

### Schedule editor keys

`↑/↓` select a window · `Tab` cycle the focused field (days / start / end) ·
`1`–`7` toggle a day (Mon–Sun) when the days field is focused · `+`/`-` adjust the
focused time by 15 min · `a` add a window · `d` delete · `t` toggle enabled · `y`
toggle always-on · `s` save · `Esc` cancel.

### Settings editor keys

`↑/↓` select a row · `←/→` or `Space` change a value (cycle engine/method, ±
numbers, toggle booleans) · `Enter` edit a text field, `Enter` again to commit /
`Esc` to cancel the edit · `s` save (validated) · `Esc` close without saving.

## First run

With no config yet, the TUI opens a setup screen: press `w` for the guided
wizard or `i` to write defaults, then the menu loads.

[← Docs index](../README.md#documentation)
