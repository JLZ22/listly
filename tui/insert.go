package tui

import (
	"strings"

	key "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type InsertKeyMap struct {
	Discard       key.Binding
	QuitNoWarning key.Binding
	Save          key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k InsertKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Discard, k.Save}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k InsertKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Discard}, // first column
		{k.Save},
	}
}

var DefaultInsertKeyMap = InsertKeyMap{
	Discard: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "discard changes"),
	),
	QuitNoWarning: key.NewBinding(
		key.WithKeys("ctrl+c"),
	),
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save"),
	),
}

func handleInsertInput(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultInsertKeyMap.Discard):
			m = insertToNormal(m)
		case key.Matches(msg, DefaultInsertKeyMap.QuitNoWarning):
			return m, tea.Quit
		case key.Matches(msg, DefaultInsertKeyMap.Save):
			str := m.editInfo.textInput.Value()
			if len(str) == 0 {
				m = insertToNormal(m)
			} else {
				m = keepChanges(m)
			}
		}
	}

	updatedTextInput, cmd := m.editInfo.textInput.Update(msg)
	m.editInfo.textInput = updatedTextInput
	return m, cmd
}

func renderInsertView(m model) string {
	lines := buildLines(m, false)
	idx := min(m.editInfo.location, m.data.list.Info.NumPending) + 1

	if m.editInfo.taskId == -1 {
		// insert the view between two tasks because we're creating a new task
		before := lines[:idx]
		after := lines[idx:]
		lines = append(before, append([]string{m.editInfo.textInput.View() + "\n"}, after...)...)
	} else {
		// replace existing task with the view because we're editing it
		before := lines[:idx]
		after := lines[idx+1:] // skip the task that we're editing

		// put the text input between before and after
		lines = append(before, append([]string{m.editInfo.textInput.View() + "\n"}, after...)...)
	}

	return strings.Join(lines, "")
}

func insertToNormal(m model) model {
	m.editInfo.textInput.Reset()
	m.mode = "normal"
	m.editInfo.taskId = -1
	return m
}

func keepChanges(m model) model {
	idx := m.editInfo.location
	if m.editInfo.taskId == -1 {
		var taskIndex int
		if len(m.data.list.TaskIds) == 0 {
			taskIndex = 0
		} else {
			if idx > 0 {
				taskIndex = getTaskIndex(m, idx-1) + 1
			} else {
				taskIndex = getTaskIndex(m, idx)
			}
		}
		_, err := m.data.list.InsertNewTask(m.editInfo.textInput.Value(), taskIndex)
		if err != nil {
			return m
		}
		m.cursor.row = idx
	} else {
		m.data.list.EditTaskDescription(m.editInfo.taskId, m.editInfo.textInput.Value())
	}

	// cleanup and return to normal mode
	m.editInfo.dirty = true
	return insertToNormal(m)
}
