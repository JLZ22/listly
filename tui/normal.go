package tui

import (
	"strings"

	help "github.com/charmbracelet/bubbles/help"
	key "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jlz22/listly/core"
)

type NormalKeyMap struct {
	Up               key.Binding
	Down             key.Binding
	QuitWithWarning  key.Binding
	QuitNoWarning    key.Binding
	NewTask          key.Binding
	NewBefore        key.Binding
	NewAfter         key.Binding
	EditTask         key.Binding
	ClearAndEdit     key.Binding
	DeleteTask       key.Binding
	ToggleCompletion key.Binding
	EnableVisualMode key.Binding
	Yank             key.Binding
	PasteAfter       key.Binding
	PasteBefore      key.Binding
	Write            key.Binding
	JumpUp           key.Binding
	JumpDown         key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k NormalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Up, k.Down, k.QuitWithWarning, k.QuitNoWarning,
		k.NewTask, k.EditTask, k.DeleteTask, k.ToggleCompletion, k.EnableVisualMode,
		k.Yank, k.PasteAfter, k.PasteBefore, k.Write, k.JumpUp, k.JumpDown, k.NewBefore, k.NewAfter,
	}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k NormalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.QuitWithWarning},       // first column
		{k.Write, k.NewTask, k.EditTask},        // second column
		{k.ClearAndEdit, k.DeleteTask, k.Yank},  // third column
		{k.EnableVisualMode, k.PasteAfter, k.PasteBefore}, // fourth column
		{k.JumpUp, k.JumpDown, k.ToggleCompletion},        // fifth column
		{k.NewBefore, k.NewAfter},
	}
}

var DefaultNormalKeyMap = NormalKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	QuitWithWarning: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	QuitNoWarning: key.NewBinding( // don't offer help for no-warning quit
		key.WithKeys("ctrl+c"),
	),
	NewTask: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new task"),
	),
	NewBefore: key.NewBinding(
		key.WithKeys("O"),
		key.WithHelp("O", "new task before"),
	),
	NewAfter: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "new task after"),
	),
	EditTask: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "edit task"),
	),
	ClearAndEdit: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "clear and edit"),
	),
	DeleteTask: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "cut task"),
	),
	ToggleCompletion: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "mark done/not done"),
	),
	EnableVisualMode: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "visual mode"),
	),
	Yank: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "yank"),
	),
	PasteAfter: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "paste"),
	),
	PasteBefore: key.NewBinding(
		key.WithKeys("P"),
		key.WithHelp("P", "paste before"),
	),
	Write: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "write"),
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

func handleNormalInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultNormalKeyMap.Up):
			if m.cursor.row > 0 {
				m.cursor.row--
			}

		case key.Matches(msg, DefaultNormalKeyMap.Down):
			if m.cursor.row < len(m.data.list.Tasks)-1 {
				m.cursor.row++
			}

		case key.Matches(msg, DefaultNormalKeyMap.QuitWithWarning):
			if m.editInfo.dirty {
				m.confirmation.active = true
				m.confirmation.message = "You have unsaved changes. Are you sure you want to quit? (y/enter = yes, n = no)"
			} else {
				return m, tea.Quit
			}

		case key.Matches(msg, DefaultNormalKeyMap.QuitNoWarning):
			return m, tea.Quit

		case key.Matches(msg, DefaultNormalKeyMap.NewTask):
			m.editInfo.taskId = -1
			m.editInfo.location = m.data.list.Info.NumPending
			m.mode = "insert"

		case key.Matches(msg, DefaultNormalKeyMap.NewBefore):
			m.editInfo.taskId = -1
			m.editInfo.location = m.cursor.row
			m.mode = "insert"
		
		case key.Matches(msg, DefaultNormalKeyMap.NewAfter):
			m.editInfo.taskId = -1
			m.editInfo.location = m.cursor.row + 1
			m.mode = "insert"

		case key.Matches(msg, DefaultNormalKeyMap.EditTask):
			taskId := getCurrTaskId(m)
			m.editInfo.taskId = taskId
			m.editInfo.textInput.SetValue(m.data.list.Tasks[taskId].Description)
			m.mode = "insert"

		case key.Matches(msg, DefaultNormalKeyMap.ClearAndEdit):
			taskId := getCurrTaskId(m)
			m.editInfo.taskId = taskId
			m.mode = "insert"

		case key.Matches(msg, DefaultNormalKeyMap.DeleteTask):
			if len(m.data.list.Tasks) == 0 {
				return m, nil
			}
			taskId := getCurrTaskId(m)
			task := m.data.list.Tasks[getCurrTaskId(m)]
			m.editInfo.copyBuff = []core.Task{*task}
			m.data.list.RemoveTask(taskId)
			m.editInfo.dirty = true

			// fix cursor position
			m.cursor.row = max(0, m.cursor.row-1)

		case key.Matches(msg, DefaultNormalKeyMap.ToggleCompletion):
			// toggle completion
			currTaskId := getCurrTaskId(m)
			currTask := m.data.list.Tasks[currTaskId]
			currTask.Done = !currTask.Done
			m.editInfo.dirty = true

			// update cursor position
			_, notDone := core.SplitByCompletion(m.data.list)
			if currTask.Done {
				// keep cursor on last not done task
				m.cursor.row = max(0, min(m.cursor.row, len(notDone)-1))
			} else {
				// keep cursor on first done task
				if m.cursor.row < len(m.data.list.TaskIds)-1 {
					m.cursor.row++
				}
			}

		case key.Matches(msg, DefaultNormalKeyMap.EnableVisualMode):
			// TODO
			m.mode = "visual"

		case key.Matches(msg, DefaultNormalKeyMap.Yank):
			task := m.data.list.Tasks[getCurrTaskId(m)]
			m.editInfo.copyBuff = []core.Task{*task}

		case key.Matches(msg, DefaultNormalKeyMap.PasteAfter):
			m = pasteTasks(m, false)
			m.cursor.row++

		case key.Matches(msg, DefaultNormalKeyMap.PasteBefore):
			m = pasteTasks(m, true)

		case key.Matches(msg, DefaultNormalKeyMap.Write):
			err := m.data.db.SaveList(m.data.list)
			if err != nil {
				m.Error = err
				return m, tea.Quit
			}
			m.editInfo.dirty = false

		case key.Matches(msg, DefaultNormalKeyMap.JumpUp):
			lastNotDone := getNumNotDone(m) - 1
			c := m.cursor.row
			if c == lastNotDone + 1 {
				m.cursor.row = lastNotDone
			} else if c > lastNotDone + 1 {
				m.cursor.row = lastNotDone + 1
			} else {
				m.cursor.row = 0
			}

		case key.Matches(msg, DefaultNormalKeyMap.JumpDown):
			lastNotDone := getNumNotDone(m) - 1
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

func renderNormalView(m model) string {
	h := help.New()
	out := buildLines(m, true)
	return strings.Join(out, "") + "\n\n" + h.FullHelpView(DefaultNormalKeyMap.FullHelp())
}

func getCurrTaskId(m model) int {
	done, notDone := core.SplitByCompletion(m.data.list)
	combined := append(notDone, done...)
	return combined[m.cursor.row].Id
}

func pasteTasks(m model, before bool) model {
	if len(m.editInfo.copyBuff) == 0 {
		return m
	}

	// create new tasks with unique id's
	newTasks := make([]core.Task, len(m.editInfo.copyBuff))
	for i, task := range m.editInfo.copyBuff {
		t, err := m.data.list.NewTask(task.Description, task.Done)
		if err != nil {
			m.Error = err
			return m
		}
		newTasks[i] = t
	}

	// find the index of the task the cursor is on - this is diff from the cursor.row because the display order is not necessarily the same as the task list order
	taskId := getCurrTaskId(m)
	taskIndex := 0
	for i, id := range m.data.list.TaskIds {
		if id == taskId {
			taskIndex = i
			break
		}
	}

	// add new tasks to the list
	for i, task := range newTasks {
		pasteIdx := taskIndex + i + 1
		if before {
			pasteIdx = taskIndex + i
		}
		err := m.data.list.Insert(task, pasteIdx)
		if err != nil {
			m.Error = err
			return m
		}
	}

	m.editInfo.dirty = true
	return m
}
