# Configuration

Configuration is a single JSON file. Manage it with the CLI, the TUI settings
pane (`vigil` → `c`), or the guided wizard:

```bash
vigil config wizard      # guided interactive setup (engine, tz, schedule, graph)
vigil config init        # write defaults
vigil config path        # print the effective path
vigil config edit        # open in $EDITOR, validated on save
vigil config validate    # check without editing
vigil config show        # print the effective config (JSON)

# Non-interactive get/set (dotted keys, validated before save):
vigil config keys                          # list settable keys
vigil config get input.interval_seconds
vigil config set engine both
vigil config set input.method key

# Schedule windows (not part of get/set):
vigil schedule list
vigil schedule add --days Mon,Tue,Wed,Thu,Fri --start 08:00 --end 17:00
vigil schedule remove 1
vigil schedule clear
```

> Setting `engine` to `graph`/`both` requires `graph.client_id` to be set first
> (the change is validated), so configure `graph.*` before switching engines.

## File locations

| Scope | Config | Runtime (status/override/log/token) |
|-------|--------|-------------------------------------|
| user (Linux/macOS) | `~/.config/vigil/config.json` (`$XDG_CONFIG_HOME`) | `~/.local/state/vigil/` (`$XDG_STATE_HOME`) |
| user (Windows) | `%AppData%\vigil\config.json` | `%LocalAppData%\vigil\` |
| system (Linux) | `/etc/vigil/config.json` | `/etc/vigil/` |
| system (macOS) | `/Library/Application Support/vigil/` | same |
| system (Windows) | `C:\ProgramData\vigil\` | same |

> macOS deliberately uses `~/.config` / `~/.local/state` (not
> `~/Library/Application Support`) so the locations are predictable across
> platforms and easy to find.

> Multi-user note: with `--scope user`, runtime files live under your own state
> dir so the CLI can always write the override file the daemon watches. For a
> `--scope system` graph daemon, runtime files are root-owned; control commands
> then require matching privilege (or enable the optional loopback API via
> `control.port`).

## Schema

```jsonc
{
  "version": 1,
  "engine": "input",                 // "input" | "graph" | "both"
  "timezone": "America/New_York",    // IANA tz; "Local" uses the host zone
  "schedule": {
    "enabled": true,                 // false => only manual overrides apply
    "always": false,                 // true => active whenever the daemon runs
    "windows": [
      { "days": ["Mon","Tue","Wed","Thu","Fri"], "start": "08:00", "end": "17:00" }
    ]
  },
  "input": {
    "interval_seconds": 60,          // pulse cadence; must be in [5,300)
    "jitter_seconds": 25,            // +/- randomization for natural cadence
    "method": "mouse",              // "mouse" (real small move) | "key" (F15) | "zen"
    "prevent_sleep": true,           // hold a sleep/display assertion while active
    "move_pixels": 3                 // max nudge magnitude; each tick uses 1..N px, random axis
  },
  "graph": {
    "tenant_id": "common",           // "common" | "organizations" | tenant GUID
    "client_id": "",                 // your Entra public-client app id
    "availability": "Available",
    "activity": "Available",
    "expiration": "PT8H",            // ISO-8601 duration; re-asserted on refresh
    "refresh_minutes": 60
  },
  "control": { "port": 0 },          // 0 = file-based control plane
  "log": { "level": "info", "max_size_mb": 5, "max_backups": 3 }
}
```

### Schedules

- Times are `HH:MM` 24-hour in `timezone`. Days are `Mon`..`Sun`.
- A window whose `end` is **not after** `start` is treated as **overnight**
  (e.g. `22:00`–`06:00`): the morning segment belongs to the day after a listed
  start day.
- Multiple windows are OR-combined. Evaluation happens in the configured zone, so
  windows track DST. On the two DST-changeover days a window boundary may be off
  by up to the one-hour shift for that day; pick interval/jitter with a little
  margin if a hard boundary matters.

### Overrides ("at will")

Overrides take precedence over the schedule and persist across restarts:

```bash
vigil on                 # force active, indefinitely
vigil on --for 2h30m     # force active for a duration
vigil off --until 18:00  # force inactive until the next 18:00
vigil resume             # clear the override, follow the schedule again
```

The config is hot-reloaded; edits and overrides take effect within ~1 second
without restarting the service.

### Running & debugging

```bash
vigil run --dry-run      # run the loop but only log intended actions (no input/Graph)
vigil run --verbose      # debug-level logging
vigil status             # current daemon + service state (add --json for scripts)
vigil doctor             # capability/permission/Graph diagnostics
```

`--dry-run` is the easiest way to watch how the schedule and overrides resolve
without actually moving the mouse or touching presence.

[← Docs index](../README.md#documentation)
