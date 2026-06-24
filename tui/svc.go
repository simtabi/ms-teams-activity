package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var serviceActions = []struct{ key, desc string }{
	{"install", "install + start the background service"},
	{"start", "start the service"},
	{"stop", "stop the service"},
	{"restart", "restart the service"},
	{"uninstall", "stop + remove the service"},
}

func (m model) updateService(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.mode = modeDashboard
	case "up", "k":
		if m.svcRow > 0 {
			m.svcRow--
		}
	case "down", "j":
		if m.svcRow < len(serviceActions)-1 {
			m.svcRow++
		}
	case "enter":
		action := serviceActions[m.svcRow].key
		m.mode = modeDashboard
		m.flash = "ran service " + action
		return m, m.execSelf(action)
	}
	return m, nil
}

func (m model) serviceView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("  Service actions") + "\n\n")
	var lines []string
	for i, a := range serviceActions {
		cursor := "  "
		row := a.key + helpStyle.Render("  — "+a.desc)
		if i == m.svcRow {
			cursor = selRowStyle.Render("▸ ")
			row = selRowStyle.Render(a.key) + helpStyle.Render("  — "+a.desc)
		}
		lines = append(lines, cursor+row)
	}
	b.WriteString(boxStyle.Render(strings.Join(lines, "\n")) + "\n")
	b.WriteString(helpStyle.Render("[↑/↓] choose  [enter] run  [esc] back"))
	return b.String()
}
