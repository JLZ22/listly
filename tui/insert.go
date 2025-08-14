package tui

import (
	"strings"

	help "github.com/charmbracelet/bubbles/help"
	key "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jlz22/listly/core"
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

func handleInsertInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultInsertKeyMap.Discard):
			m = returnToNormal(m)
		case key.Matches(msg, DefaultInsertKeyMap.QuitNoWarning):
			return m, tea.Quit
		case key.Matches(msg, DefaultInsertKeyMap.Save):
			str := m.editInfo.textInput.Value()
			if len(str) == 0 {
				m = returnToNormal(m)
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
	h := help.New()
	lines := buildLines(m, false)
	idx := min(m.editInfo.location, getNumNotDone(m))

	if m.editInfo.taskId == -1 {
		// insert the view between two tasks because we're creating a new task
		before := lines[:idx+1]
		after := lines[idx+1:]
		lines = append(before, append([]string{m.editInfo.textInput.View() + "\n"}, after...)...)
	} else {
		// replace existing task with the view because we're editing it
		before := lines[:idx+1]
		after := lines[idx+2:] // skip the task line
		lines = append(before, append([]string{m.editInfo.textInput.View() + "\n"}, after...)...)
	}

	return strings.Join(lines, "") + "\n\n" + h.FullHelpView(DefaultInsertKeyMap.FullHelp())
}

func returnToNormal(m model) model {
	m.editInfo.textInput.Reset()
	m.mode = "normal"
	m.editInfo.taskId = -1
	return m
}

func keepChanges(m model) model {
	idx := m.editInfo.location
	if m.editInfo.taskId == -1 {
		// If taskId is -1, we're creating a new task
		if idx == -1 {
			_, err := m.data.list.AddNewTask(m.editInfo.textInput.Value(), false)
			if err != nil {
				m.Error = err
				return m
			}
		} else {
			_, err := m.data.list.InsertNewTask(m.editInfo.textInput.Value(), idx)
			if err != nil {
				m.Error = err
				return m
			}
		}

		m.cursor.row = getNumNotDone(m) - 1
	} else {
		// editing existing task
		m.data.list.EditTaskDescription(m.editInfo.taskId, m.editInfo.textInput.Value())
		_, notDone := core.SplitByCompletion(m.data.list)
		m.cursor.row = len(notDone) - 1 // move cursor to the end of the list
	}

	// cleanup and return to normal mode
	m.editInfo.dirty = true
	return returnToNormal(m)
}