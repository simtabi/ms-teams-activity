# Input engine

The default engine. It injects small, periodic, human-like input that resets the
OS idle timer so Teams stays Available, and (where supported) holds a power
assertion to defer sleep/lock.

Enable it with `"engine": "input"` (or `"both"`).

## Methods

- `mouse` (default) — a tiny **real** relative move that immediately returns to
  origin. The offset is varied (1–3 px, random direction) and the interval is
  jittered, so movement is natural rather than mechanically identical each tick.
  Most reliable for apps with their own idle detection.
- `key` — taps **F15**, a key with no default action.
- `zen` — an in-place/zero-delta nudge. Least intrusive, but some idle detectors
  ignore it. Opt-in.

`interval_seconds` (default 60) must be below Teams' ~5-minute idle threshold and
below your OS auto-lock timeout. `jitter_seconds` randomizes the cadence.

## Per-OS behavior

```
                        ┌─ Windows ─► SendInput (real move / F15),
                        │             runs as a logon Scheduled Task
  reset the OS idle ────┼─ Linux ───► /dev/uinput virtual device
  timer every           │             (real kernel events; X11 + Wayland)
  interval ± jitter     └─ macOS ───► CGEventPost + IOPMAssertion —
                                      keeps Teams green but CANNOT reset the
                                      hardware idle timer (a lock still → Away)
```

### Windows
`SendInput` updates `GetLastInputInfo` reliably. The engine installs as a
**logon Scheduled Task** so it runs in your interactive session (a session-0
service cannot inject input). `prevent_sleep` uses `SetThreadExecutionState`.

### Linux
A `/dev/uinput` virtual mouse emits **real kernel events**, resetting idle under
X11 and Wayland. Requires uinput access (see installation). All methods map to a
tiny mouse move (the most reliable real event). The Teams Linux client is
web/PWA, which may apply its own tab-level idle in addition to OS idle.

### macOS
`CGEventPost` posts synthetic mouse/key events and the engine holds an
`IOPMAssertion` to defer display sleep. Two caveats:

1. **Accessibility permission is required** for `CGEventPost` to take effect —
   grant it to the `vigil` binary and verify with `vigil doctor`.
2. **Synthetic events do not reset the hardware idle timer.** Teams usually
   reads the synthetic-aware combined-session idle (so it stays green), but the
   screensaver/auto-lock uses the hardware timer. If auto-lock fires, Teams goes
   Away. Disable or lengthen auto-lock, or use the `graph` engine.

[← Docs index](../../README.md#documentation)
