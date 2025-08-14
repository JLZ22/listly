package tui

import (
	key "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type VisualKeyMap struct {
	Up            key.Binding
	Down          key.Binding
	NormalMode    key.Binding
	QuitNoWarning key.Binding
	Delete        key.Binding
	Yank          key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k VisualKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.NormalMode, k.QuitNoWarning, k.Delete, k.Yank}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k VisualKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Yank}, // first column
		{k.Down, k.NormalMode}, // second column
		{k.Delete, k.QuitNoWarning},
	}
}

var DefaultVisualKeyMap = VisualKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	NormalMode: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "return to normal mode"),
	),
	QuitNoWarning: key.NewBinding(
		key.WithKeys("ctrl+c"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d/x"),
		key.WithHelp("d/x", "cut selection"),
	),
	Yank: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "yank"),
	),
}

func handleVisualInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	return m, nil
}

func renderVisualView(m model) string {
	return "In Visual Mode"
}
