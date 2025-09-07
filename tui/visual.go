package tui

import (
	"strings"

	key "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jlz22/listly/core"
)

var visualHighlightStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#5fa2ffff")). // red background
	Foreground(lipgloss.Color("#ffffff")).   // white text
	Bold(true)

type VisualKeyMap struct {
	Up               key.Binding
	Down             key.Binding
	UpFive           key.Binding
	DownFive         key.Binding
	NormalMode       key.Binding
	QuitNoWarning    key.Binding
	Delete           key.Binding
	Yank             key.Binding
	ToggleCompletion key.Binding
	JumpUp           key.Binding
	JumpDown         key.Binding
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
		{k.Up, k.Yank},         // first column
		{k.Down, k.NormalMode}, // second column
		{k.Delete, k.QuitNoWarning},
		{k.JumpUp, k.JumpDown},
	}
}

var DefaultVisualKeyMap = VisualKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	UpFive: key.NewBinding(
		key.WithKeys("K"),
		key.WithHelp("K", "up 5"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	DownFive: key.NewBinding(
		key.WithKeys("J"),
		key.WithHelp("J", "down 5"),
	),
	NormalMode: key.NewBinding(
		key.WithKeys("esc", "v"),
		key.WithHelp("esc/v", "normal mode"),
	),
	QuitNoWarning: key.NewBinding(
		key.WithKeys("ctrl+c"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d", "x"),
		key.WithHelp("d/x", "cut"),
	),
	Yank: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "yank"),
	),
	ToggleCompletion: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "mark done/not done"),
	),
	JumpUp: key.NewBinding(
		key.WithKeys("{"),
		key.WithHelp("{", "jump up"),
	),
	JumpDown: key.NewBinding(
		key.WithKeys("}"),
		key.WithHelp("}", "jump down"),
	),
}

func handleVisualInput(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultVisualKeyMap.Up):
			if m.cursor.row > 0 {
				m.cursor.row--
			}

		case key.Matches(msg, DefaultVisualKeyMap.Down):
			if m.cursor.row < len(m.data.list.Tasks)-1 {
				m.cursor.row++
			}

		case key.Matches(msg, DefaultVisualKeyMap.UpFive):
			if m.cursor.row > 4 {
				m.cursor.row -= 4
			} else {
				m.cursor.row = 0
			}

		case key.Matches(msg, DefaultVisualKeyMap.DownFive):
			if m.cursor.row < len(m.data.list.Tasks)-6 {
				m.cursor.row += 4
			} else {
				m.cursor.row = len(m.data.list.Tasks) - 1
			}

		case key.Matches(msg, DefaultVisualKeyMap.NormalMode):
			m = visualToNormal(m)

		case key.Matches(msg, DefaultVisualKeyMap.QuitNoWarning):
			return m, tea.Quit

		case key.Matches(msg, DefaultVisualKeyMap.Delete):
			copyBuff := copySelection(m)
			m.editInfo.copyBuff = copyBuff

			// delete the selection
			start := min(m.cursor.selStart, m.cursor.row)
			end := max(m.cursor.selStart, m.cursor.row)
			numToRemove := end - start
			done, notDone := core.SplitByCompletion(m.data.list)
			combined := append(notDone, done...)
			for i := 0; i <= numToRemove; i++ {
				taskId := combined[start+i].Id
				err := m.data.list.RemoveTask(taskId)
				if err != nil {
					panic(err)
				}
			}
			m = visualToNormal(m)
			m.cursor.row = min(m.data.list.Info.NumTasks-1, m.cursor.row)

		case key.Matches(msg, DefaultVisualKeyMap.Yank):
			copyBuff := copySelection(m)
			m.editInfo.copyBuff = copyBuff
			m = visualToNormal(m)

		case key.Matches(msg, DefaultVisualKeyMap.ToggleCompletion):
			start := min(m.cursor.selStart, m.cursor.row)
			end := max(m.cursor.selStart, m.cursor.row)
			done, notDone := core.SplitByCompletion(m.data.list)
			combined := append(notDone, done...)
			for i := start; i <= end; i++ {
				task := combined[i]
				task.Done = !task.Done
				if task.Done {
					m.data.list.Info.NumDone++
					m.data.list.Info.NumPending--
				} else {
					m.data.list.Info.NumDone--
					m.data.list.Info.NumPending++
				}
				m.cursor.row--
			}
			m.cursor.row = max(0, m.cursor.row)
			m = visualToNormal(m)

		case key.Matches(msg, DefaultNormalKeyMap.JumpUp):
			lastNotDone := m.data.list.Info.NumPending - 1
			c := m.cursor.row
			if c == lastNotDone+1 {
				m.cursor.row = lastNotDone
			} else if c > lastNotDone+1 {
				m.cursor.row = lastNotDone + 1
			} else {
				m.cursor.row = 0
			}

		case key.Matches(msg, DefaultNormalKeyMap.JumpDown):
			lastNotDone := m.data.list.Info.NumPending - 1
			c := m.cursor.row

			if c < lastNotDone {
				m.cursor.row = lastNotDone
			} else if c == lastNotDone {
				m.cursor.row = lastNotDone + 1
			} else {
				m.cursor.row = len(m.data.list.TaskIds) - 1
			}
		}
	}

	return m, nil
}

func renderVisualView(m model) string {
	lines := buildLines(m, true)

	sepBarIdx := m.data.list.Info.NumPending + 1
	start := min(m.cursor.selStart, m.cursor.row)
	end := max(m.cursor.selStart, m.cursor.row)
	numSelected := end - start

	idx := start + 1
	if start >= sepBarIdx {
		idx += 0 // skip the separating bar
	}
	for i := 0; i <= numSelected; {
		if idx == sepBarIdx {
			// skip the separating bar
			idx++
		}
		trimmed := strings.TrimRight(lines[idx][4:], " \t\n")
		lines[idx] = "    " + visualHighlightStyle.Render(trimmed) + "\n"
		i++
		idx++
	}

	return strings.Join(lines, "")
}

func copySelection(m model) []core.Task {
	// create buff with correct length
	start := min(m.cursor.selStart, m.cursor.row)
	end := max(m.cursor.selStart, m.cursor.row)
	copyBuff := make([]core.Task, end-start+1)

	done, notDone := core.SplitByCompletion(m.data.list)
	combined := append(notDone, done...)

	// fill buff
	for i := start; i <= end; i++ {
		copyBuff[i-start] = *m.data.list.Tasks[combined[i].Id]
	}
	return copyBuff
}

func visualToNormal(m model) model {
	m.cursor.selStart = -1
	m.mode = "normal"
	return m
}
