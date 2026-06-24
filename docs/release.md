# Release process

Releases are tag-driven. The git tag is the single source of version truth;
`mta version` is stamped at build time via ldflags
(`internal/cli.version/commit/date`).

## Cutting a release

1. Move `[Unreleased]` items in `CHANGELOG.md` under a new `## [X.Y.Z]` section
   with the date.
2. Commit on `main`.
3. Tag and push:

   ```bash
   git tag vX.Y.Z
   git push origin vX.Y.Z
   ```

The `Release` workflow (`.github/workflows/release.yml`) then runs.

## Pipeline (hand-rolled, OSS)

> GoReleaser's `prebuilt` builder is Pro-only, and the macOS backend needs cgo
> (can't cross-compile from Linux), so the pipeline builds natively per-OS.

1. **build matrix** (one job per target): macOS arm64 + amd64 (cgo), Windows
   amd64, Linux amd64 + arm64 (`CGO_ENABLED=0`). Each job builds `./cmd/mta`,
   archives it as **`mta_<os>_<arch>.{tar.gz,zip}`** (version-less — required by
   the self-update contract), and Linux jobs also build **deb/rpm via `nfpm`**.
2. **release**: collects all artifacts, writes **`checksums.txt`** (sha256),
   extracts the tagged `CHANGELOG.md` section as the release body (plus
   `generate_release_notes`), and publishes the GitHub Release with every asset.
3. **brew-scoop** (best-effort): templates the Homebrew formula and Scoop
   manifest with the archive sha256s and pushes them to `simtabi/homebrew-tap`
   and `simtabi/scoop-bucket`.

## Self-update contract

`internal/selfupdate` downloads `mta_<os>_<arch>.<ext>` and validates it against
`checksums.txt`. Keep these aligned: the archive name template, the inner binary
name (`mta`/`mta.exe`), `checksums.txt`, and `vX.Y.Z` tags. Changing one side
means changing both.

## First-release prerequisites

1. Make the repo public.
2. Create `simtabi/homebrew-tap` and `simtabi/scoop-bucket` repositories.
3. Add a repo secret **`TAP_GITHUB_TOKEN`** (a repo-scoped PAT that can push to
   those two repos).
4. Push the first `vX.Y.Z` tag.

Every release carries a real, human-readable description sourced from the
CHANGELOG — never a bare "see changelog" stub. Keep GitHub Actions current via
Dependabot.

[← Docs index](../README.md#documentation)
