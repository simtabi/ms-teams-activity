# Installation

## Options

| Method | Command |
|--------|---------|
| Install script (macOS/Linux) | `curl -fsSL https://raw.githubusercontent.com/simtabi/vigil/main/scripts/install.sh \| sh` |
| Install script (Windows) | `irm https://raw.githubusercontent.com/simtabi/vigil/main/scripts/install.ps1 \| iex` |
| Homebrew | `brew install simtabi/tap/vigil` |
| Scoop | `scoop bucket add simtabi https://github.com/simtabi/scoop-bucket; scoop install vigil` |
| deb/rpm | download from releases, `sudo dpkg -i vigil_*.deb` / `sudo rpm -i vigil_*.rpm` |
| go install | `go install github.com/simtabi/vigil/cmd/vigil@latest` |
| Prebuilt archive | download `vigil_<os>_<arch>.{tar.gz,zip}` from the releases page |
| From source | `go build -o vigil ./cmd/vigil` |

The install scripts download the prebuilt binary and **verify its SHA-256**
against the release `checksums.txt`, falling back to a source build if the
download fails and Go is present.

> **go install / source on macOS** needs a C toolchain (Xcode CLT) because the
> macOS input backend uses cgo. Windows and Linux are pure Go.

## Building all targets yourself

`make dist` (or `./scripts/build-all.sh [version]`) builds ready-to-run binaries
for every target in `build/targets.txt`. The `dist/` root holds the bare,
self-describing binaries; all archives + deb/rpm are grouped under
`dist/archives/`:

```
dist/
  vigil_macos_arm64   vigil_linux_amd64   vigil_windows_amd64.exe   vigil_macos_universal   …
  checksums.txt
  archives/  vigil_<os>_<arch>.tar.gz · vigil_windows_<arch>.zip · *.deb · *.rpm · checksums.txt
```

It builds whatever your local toolchain supports (macOS targets need a C
compiler; everything else is pure-Go cross-compilation). CI runs the same script
on every push to `main`, so prebuilt binaries are always downloadable from the
latest run's Artifacts.

## Putting the binary on PATH

`vigil self install` copies the running binary to a standard location
(`~/.local/bin`, `/usr/local/bin` with `--scope system`, or
`%LOCALAPPDATA%\Programs\vigil` on Windows). `vigil self uninstall [--purge]`
removes the service and the binary (and, with `--purge`, config + data).

## First run

```bash
vigil config wizard   # guided setup (or: vigil config init for defaults)
vigil doctor          # verify capabilities and permissions
vigil install         # install + start the background service
vigil status          # check service + daemon state
```

## Daemon lifecycle

The background service is managed with these commands (they auto-select the
right per-OS mechanism — launchd / systemd / Windows service / Windows logon
task):

```bash
vigil install --init   # turnkey: write a default config if missing, then install + start
vigil start | stop | restart | status
vigil uninstall        # stop + remove the service (keeps the binary)
```

The install scripts can do the whole thing in one shot with `--with-service`
(`-WithService` on Windows). To remove everything later:

```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/simtabi/vigil/main/scripts/uninstall.sh | sh
curl -fsSL .../scripts/uninstall.sh | sh -s -- --purge   # also delete config + runtime data

# Windows (piped form can't pass flags)
irm https://raw.githubusercontent.com/simtabi/vigil/main/scripts/uninstall.ps1 | iex
./scripts/uninstall.ps1 -Purge                            # local file: also delete config + data

# Or directly, on any OS:
vigil self uninstall            # service + binary
vigil self uninstall --purge    # also config + runtime data
```

`--purge` (and `-Purge`) additionally delete the config and runtime/state
directories; without it, those are left in place.

## Scope: user vs system

`--scope user` (default) installs a per-user service that runs in your desktop
session. `--scope system` installs a machine-wide service.

**The input engine requires a desktop (GUI) session**, so it must be installed
with `--scope user`. `vigil install` refuses `input` + `--scope system`. A
system-wide service is appropriate for the `graph` engine, which is headless.

## Per-OS service mechanism

| OS | `input` engine | `graph` engine |
|----|----------------|----------------|
| macOS | LaunchAgent (user) | LaunchAgent (user) or LaunchDaemon (system) |
| Linux | `systemd --user` (run `loginctl enable-linger $USER` to persist when logged out) | systemd user or system |
| Windows | **logon Scheduled Task** (interactive session) | Windows service |

## Platform prerequisites

- **macOS** — grant the `vigil` binary **Accessibility** permission
  (System Settings → Privacy & Security → Accessibility). Because TCC keys on
  the binary's signature, re-grant after rebuilding an unsigned binary. Run
  `vigil doctor` to confirm. Synthetic input cannot reset the *hardware* idle
  timer, so disable or lengthen auto-lock to stay Available.
- **Linux** — `/dev/uinput` must exist (`sudo modprobe uinput`) and be writable
  by your user. Add a udev rule / group so the device is accessible without
  root, e.g.:

  ```
  KERNEL=="uinput", GROUP="input", MODE="0660", OPTIONS+="static_node=uinput"
  ```

  then add yourself to the `input` group and re-login.
- **Windows** — no special setup; the input engine installs as a logon task.

[← Docs index](../README.md#documentation)
