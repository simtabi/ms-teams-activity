package cli

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/simtabi/ms-teams-activity/internal/config"
	"github.com/simtabi/ms-teams-activity/internal/selfupdate"
	"github.com/simtabi/ms-teams-activity/internal/service"
	"github.com/spf13/cobra"
)

var (
	flagCheck bool
	flagYes   bool
	flagPurge bool
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Update mta to the latest release (alias of `self update`)",
	RunE:  func(_ *cobra.Command, _ []string) error { return doUpgrade() },
}

var selfCmd = &cobra.Command{
	Use:   "self",
	Short: "Manage the mta binary itself (update/install/uninstall)",
}

var selfUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update mta to the latest release",
	RunE:  func(_ *cobra.Command, _ []string) error { return doUpgrade() },
}

var selfInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Copy this binary to a standard location on PATH",
	RunE:  func(_ *cobra.Command, _ []string) error { return doSelfInstall() },
}

var selfUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the service and the mta binary (--purge also deletes config/data)",
	RunE:  func(_ *cobra.Command, _ []string) error { return doSelfUninstall() },
}

func doUpgrade() error {
	ctx := context.Background()
	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return err
	}
	ch := selfupdate.DetectChannel(exe)

	// Always allow a check; only block the apply for dev/package-managed.
	if flagCheck {
		info, err := selfupdate.Check(ctx, version)
		if err != nil {
			return err
		}
		printCheck(info, ch)
		return nil
	}

	if selfupdate.IsDev(version) {
		return fmt.Errorf("%w (use a released build, or `mta upgrade --check`)", selfupdate.ErrDevVersion)
	}
	if !ch.SelfUpdatable() {
		return fmt.Errorf("not self-updating: %s", ch.Advice())
	}

	info, err := selfupdate.Check(ctx, version)
	if err != nil {
		return err
	}
	if !info.Available {
		fmt.Printf("already up to date (%s)\n", version)
		return nil
	}
	if !flagYes && !confirm(fmt.Sprintf("update %s → %s?", info.Current, info.Latest)) {
		fmt.Println("cancelled")
		return nil
	}

	// Stop a running service so the binary file can be replaced cleanly.
	params, haveParams := serviceParamsBestEffort()
	wasRunning := false
	if haveParams {
		if st, err := service.StatusString(params); err == nil && st == "running" {
			wasRunning = true
			fmt.Println("stopping service...")
			_ = service.Stop(params)
		}
	}

	applied, err := selfupdate.Apply(ctx, version)
	if wasRunning {
		fmt.Println("restarting service...")
		_ = service.Start(params)
	}
	if err != nil {
		return err
	}
	fmt.Printf("updated to %s\n", applied.Latest)
	tccReminderIfNeeded()
	return nil
}

func printCheck(info selfupdate.Info, ch selfupdate.Channel) {
	if flagJSON {
		_ = printJSON(map[string]any{
			"current": info.Current, "latest": info.Latest,
			"available": info.Available, "channel": ch.SelfUpdatable(),
		})
		return
	}
	switch {
	case info.Latest == "":
		fmt.Println("no releases found")
	case info.Available:
		fmt.Printf("update available: %s → %s\n", info.Current, info.Latest)
		if !ch.SelfUpdatable() {
			fmt.Println(ch.Advice())
		}
	default:
		fmt.Printf("up to date (%s)\n", info.Current)
	}
}

func doSelfInstall() error {
	src, err := selfupdate.ExecutablePath()
	if err != nil {
		return err
	}
	dir, err := binInstallDir(scope())
	if err != nil {
		return err
	}
	dst := dir + string(os.PathSeparator) + binName()
	if err := copyFile(src, dst, 0o755); err != nil {
		return err
	}
	fmt.Printf("installed: %s\n", dst)
	fmt.Printf("ensure %s is on your PATH, then run `mta config init` and `mta install`.\n", dir)
	tccReminderIfNeeded()
	return nil
}

func doSelfUninstall() error {
	if !flagYes && !confirm("remove the mta service and binary?") {
		fmt.Println("cancelled")
		return nil
	}
	// Remove the service/task first (best effort).
	if params, ok := serviceParamsBestEffort(); ok {
		if err := service.Uninstall(params); err != nil {
			fmt.Fprintln(os.Stderr, "warning: service uninstall:", err)
		}
	}

	if flagPurge {
		if dir, err := config.ConfigDir(scope()); err == nil {
			_ = os.RemoveAll(dir)
		}
		if dir, err := config.RuntimeDir(scope()); err == nil {
			_ = os.RemoveAll(dir)
		}
		fmt.Println("removed config and runtime data")
	}

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		fmt.Printf("service removed. Delete the binary manually (it is running): %s\n", exe)
		return nil
	}
	if err := os.Remove(exe); err != nil {
		fmt.Fprintf(os.Stderr, "could not remove binary %s: %v\n", exe, err)
		return nil
	}
	fmt.Printf("removed binary: %s\n", exe)
	return nil
}

func init() {
	upgradeCmd.Flags().BoolVar(&flagCheck, "check", false, "only report whether an update is available")
	upgradeCmd.Flags().BoolVar(&flagYes, "yes", false, "do not prompt for confirmation")
	selfUpdateCmd.Flags().BoolVar(&flagCheck, "check", false, "only report whether an update is available")
	selfUpdateCmd.Flags().BoolVar(&flagYes, "yes", false, "do not prompt for confirmation")
	selfUninstallCmd.Flags().BoolVar(&flagYes, "yes", false, "do not prompt for confirmation")
	selfUninstallCmd.Flags().BoolVar(&flagPurge, "purge", false, "also delete config and runtime data")

	selfCmd.AddCommand(selfUpdateCmd, selfInstallCmd, selfUninstallCmd)
	rootCmd.AddCommand(upgradeCmd, selfCmd)
}
