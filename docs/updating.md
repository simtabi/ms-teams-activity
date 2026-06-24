# Updating

## Self-update (standalone installs)

```bash
mta upgrade            # download + verify + replace this binary, then restart the service
mta upgrade --check    # report only; don't change anything
mta upgrade --yes      # skip the confirmation prompt
mta self update        # identical to `mta upgrade`
```

`mta upgrade` checks the latest GitHub release, downloads the archive for your
OS/arch, **verifies it against `checksums.txt`**, replaces the running binary,
and — if a service/logon-task is installed — stops it around the swap and starts
it again.

## Package-manager installs

If you installed through a package manager, `mta upgrade` refuses and points you
to the right tool:

| Installed via | Update with |
|---------------|-------------|
| Homebrew | `brew upgrade mta` |
| Scoop | `scoop update mta` |
| deb/rpm | `apt upgrade` / `dnf upgrade` |
| `go install` | re-run `go install …/cmd/mta@latest` |

The interactive TUI shows an "update available" banner only for standalone
installs.

## Caveats

- **Dev builds** (`mta version` shows `dev`) can't self-update — install a
  released build. `--check` still works.
- **macOS Accessibility (TCC).** Replacing the binary changes its signature, so
  macOS may drop the Accessibility grant. After updating, re-grant it
  (System Settings → Privacy & Security → Accessibility) and run `mta doctor`.
- **Windows.** A running executable can't be deleted; the updater replaces it by
  renaming. The background **logon task** is stopped first so its copy is free.
- **Service binary path.** Self-update replaces the binary you ran
  (`mta upgrade`). If your service points at a different path, re-run
  `mta install` afterwards so it uses the updated binary.

[← Docs index](../README.md#documentation)
