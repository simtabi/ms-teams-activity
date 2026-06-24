package cli

import "testing"

// TestCommandTree guards the CLI surface: every command has help text, names are
// unique, and the documented top-level commands stay wired up.
func TestCommandTree(t *testing.T) {
	seen := map[string]bool{}
	for _, c := range rootCmd.Commands() {
		name := c.Name()
		if seen[name] {
			t.Errorf("duplicate command %q", name)
		}
		seen[name] = true
		if c.Short == "" {
			t.Errorf("command %q has no Short description", name)
		}
		// A command must either do something or group subcommands.
		if c.RunE == nil && c.Run == nil && !c.HasSubCommands() {
			t.Errorf("command %q has neither an action nor subcommands", name)
		}
	}

	required := []string{
		"run", "install", "uninstall", "start", "stop", "restart",
		"on", "off", "resume", "status", "config", "schedule",
		"auth", "doctor", "self", "upgrade", "tui", "version",
	}
	for _, n := range required {
		if !seen[n] {
			t.Errorf("expected top-level command %q to be registered", n)
		}
	}
}

func TestSubcommandsRegistered(t *testing.T) {
	groups := map[string][]string{
		"config":   {"init", "get", "set", "keys", "wizard", "validate", "show", "path", "edit"},
		"schedule": {"list", "add", "remove", "clear"},
		"self":     {"update", "install", "uninstall"},
		"auth":     {"login", "status", "logout"},
	}
	for parent, subs := range groups {
		cmd, _, err := rootCmd.Find([]string{parent})
		if err != nil {
			t.Fatalf("find %q: %v", parent, err)
		}
		have := map[string]bool{}
		for _, c := range cmd.Commands() {
			have[c.Name()] = true
		}
		for _, s := range subs {
			if !have[s] {
				t.Errorf("%s: expected subcommand %q", parent, s)
			}
		}
	}
}

func TestVersionRuns(t *testing.T) {
	rootCmd.SetArgs([]string{"version", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}
	rootCmd.SetArgs(nil)
}
