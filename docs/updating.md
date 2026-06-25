# Updating

## Self-update (standalone installs)

```bash
vigil upgrade            # download + verify + replace this binary, then restart the service
vigil upgrade --check    # report only; don't change anything
vigil upgrade --yes      # skip the confirmation prompt
vigil self update        # identical to `vigil upgrade`
```

`vigil upgrade` checks the latest GitHub release, downloads the archive for your
OS/arch, **verifies it against `checksums.txt`**, replaces the running binary,
and — if a service/logon-task is installed — stops it around the swap and starts
it again.

```
 vigil upgrade
   ├─ dev build? ──────────► refuse (install a release); --check still works
   ├─ package-managed? ────► refuse → brew / scoop / apt (see below)
   └─ standalone:
        check latest ─ newer? ─ no ─► "already up to date"
                          │ yes
                   confirm (or --yes)
                          ▼
        stop service ─► download + verify checksums.txt ─► replace binary
                          ▼
        restart service  (macOS: re-grant Accessibility — signature changed)
```

## Package-manager installs

If you installed through a package manager, `vigil upgrade` refuses and points you
to the right tool:

| Installed via | Update with |
|---------------|-------------|
| Homebrew | `brew upgrade vigil` |
| Scoop | `scoop update vigil` |
| deb/rpm | `apt upgrade` / `dnf upgrade` |
| `go install` | re-run `go install …/cmd/vigil@latest` |

The interactive TUI shows an "update available" banner only for standalone
installs.

## Caveats

- **Dev builds** (`vigil version` shows `dev`) can't self-update — install a
  released build. `--check` still works.
- **macOS Accessibility (TCC).** Replacing the binary changes its signature, so
  macOS may drop the Accessibility grant. After updating, re-grant it
  (System Settings → Privacy & Security → Accessibility) and run `vigil doctor`.
- **Windows.** A running executable can't be deleted; the updater replaces it by
  renaming. The background **logon task** is stopped first so its copy is free.
- **Service binary path.** Self-update replaces the binary you ran
  (`vigil upgrade`). If your service points at a different path, re-run
  `vigil install` afterwards so it uses the updated binary.

[← Docs index](../README.md#documentation)
