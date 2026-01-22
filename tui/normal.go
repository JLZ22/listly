package tui

import (
	"strings"

	key "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jlz22/listly/core"
)

func handleNormalInput(msg tea.Msg, m model) (model, tea.Cmd) {
	switch m.confirmation.active {
	case true:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
				return m, tea.Quit

			case key.Matches(msg, key.NewBinding(key.WithKeys("n", "esc"))):
				m.confirmation.active = false
				m.confirmation.message = ""
				return m, nil

			case key.Matches(msg, key.NewBinding(key.WithKeys("y", "enter"))):
				return m, tea.Quit
			}
		}
	case false:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.kmap.Normal.Up):
				if m.cursor.row > 0 {
					m.cursor.row--
				}

			case key.Matches(msg, m.kmap.Normal.Down):
				if m.cursor.row < len(m.data.list.Tasks)-1 {
					m.cursor.row++
				}

			case key.Matches(msg, m.kmap.Normal.UpFive):
				if m.cursor.row > 4 {
					m.cursor.row -= 4
				} else {
					m.cursor.row = 0
				}

			case key.Matches(msg, m.kmap.Normal.DownFive):
				if m.cursor.row < len(m.data.list.Tasks)-6 {
					m.cursor.row += 4
				} else {
					m.cursor.row = len(m.data.list.Tasks) - 1
				}

			case key.Matches(msg, m.kmap.Normal.QuitWithWarning):
				if m.editInfo.dirty {
					m.confirmation.active = true
					m.confirmation.message = "You have unsaved changes. Are you sure you want to quit? (y/enter = yes, n = no)"
				} else {
					return m, tea.Quit
				}

			case key.Matches(msg, m.kmap.Normal.QuitNoWarning):
				return m, tea.Quit

			case key.Matches(msg, m.kmap.Normal.NewTask):
				m.editInfo.taskId = -1
				m.editInfo.location = m.data.list.Info.NumPending
				m.mode = "insert"

			case key.Matches(msg, m.kmap.Normal.NewBefore):
				if m.cursor.row >= m.data.list.Info.NumPending {
					m.editInfo.taskId = -1
					m.editInfo.location = m.data.list.Info.NumPending
					m.mode = "insert"
				} else {
					m.editInfo.taskId = -1
					m.editInfo.location = m.cursor.row
					m.mode = "insert"
				}

			case key.Matches(msg, m.kmap.Normal.NewAfter):
				if m.cursor.row >= m.data.list.Info.NumPending {
					m.editInfo.taskId = -1
					m.editInfo.location = m.data.list.Info.NumPending
					m.mode = "insert"
				} else {
					m.editInfo.taskId = -1
					m.editInfo.location = m.cursor.row + 1
					m.mode = "insert"
				}

			case key.Matches(msg, m.kmap.Normal.EditTask):
				if m.data.list.Info.NumTasks < 1 {
					m.editInfo.taskId = -1
					m.editInfo.location = m.data.list.Info.NumPending
					m.mode = "insert"
				} else {
					taskId := getTaskId(m, m.cursor.row)
					m.editInfo.taskId = taskId
					m.editInfo.textInput.SetValue(m.data.list.Tasks[taskId].Description)
					m.editInfo.location = m.cursor.row
					m.mode = "insert"
				}

			case key.Matches(msg, m.kmap.Normal.ClearAndEdit):
				if m.data.list.Info.NumTasks < 1 {
					m.editInfo.taskId = -1
					m.editInfo.location = m.data.list.Info.NumPending
					m.mode = "insert"
				} else {
					taskId := getTaskId(m, m.cursor.row)
					m.editInfo.taskId = taskId
					m.mode = "insert"
				}

			case key.Matches(msg, m.kmap.Normal.DeleteTask):
				if len(m.data.list.Tasks) == 0 {
					return m, nil
				}
				taskId := getTaskId(m, m.cursor.row)
				task := m.data.list.Tasks[getTaskId(m, m.cursor.row)]
				m.editInfo.copyBuff = []core.Task{*task}
				m.data.list.RemoveTask(taskId)
				m.editInfo.dirty = true

				// fix cursor position
				m.cursor.row = max(0, min(m.cursor.row, m.data.list.Info.NumTasks-1))

			case key.Matches(msg, m.kmap.Normal.ToggleCompletion):
				if m.data.list.Info.NumTasks < 1 {
					break
				}

				// toggle completion
				currTaskId := getTaskId(m, m.cursor.row)
				currTask := m.data.list.Tasks[currTaskId]
				currTask.Done = !currTask.Done
				m.editInfo.dirty = true

				// update cursor position and list info
				if currTask.Done {
					m.data.list.Info.NumPending -= 1
					m.data.list.Info.NumDone += 1
					// keep cursor on last not done task
					m.cursor.row = max(0, min(m.cursor.row, m.data.list.Info.NumPending-1))

				} else {
					m.data.list.Info.NumPending += 1
					m.data.list.Info.NumDone -= 1
					// keep cursor on first done task
					if m.cursor.row < m.data.list.Info.NumPending && m.data.list.Info.NumDone > 0 {
						m.cursor.row++
					}
				}

			case key.Matches(msg, m.kmap.Normal.EnableVisualMode):
				if m.data.list.Info.NumTasks > 0 {
					m.cursor.selStart = m.cursor.row
					m.mode = "visual"
				}

			case key.Matches(msg, m.kmap.Normal.Yank):
				if m.data.list.Info.NumTasks > 0 {
					task := m.data.list.Tasks[getTaskId(m, m.cursor.row)]
					m.editInfo.copyBuff = []core.Task{*task}
				}

			case key.Matches(msg, m.kmap.Normal.PasteAfter):
				m = pasteTasks(m, false)

			case key.Matches(msg, m.kmap.Normal.PasteBefore):
				m = pasteTasks(m, true)

			case key.Matches(msg, m.kmap.Normal.Write):
				err := core.WithDefaultDB(func(db *core.DB) error {
					return db.SaveList(m.data.list)
				})
				if err != nil {
					panic(err)
				}
				m.editInfo.dirty = false

			case key.Matches(msg, m.kmap.Normal.JumpUp):
				lastNotDone := m.data.list.Info.NumPending - 1
				c := m.cursor.row
				if c == lastNotDone+1 {
					m.cursor.row = lastNotDone
				} else if c > lastNotDone+1 {
					m.cursor.row = lastNotDone + 1
				} else {
					m.cursor.row = 0
				}

			case key.Matches(msg, m.kmap.Normal.JumpDown):
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
	}
	return m, nil
}

func renderNormalView(m model) string {
	out := buildLines(m, true)
	return strings.Join(out, "")
}

func getTaskId(m model, displayIdx int) int {
	done, notDone := core.SplitByCompletion(m.data.list)
	combined := append(notDone, done...)
	return combined[displayIdx].Id
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
			panic(err)
		}
		newTasks[i] = t
	}

	// add new tasks to the list
	taskIndex := 0
	if m.data.list.Info.NumTasks-1 > 0 {
		taskIndex = getTaskIndex(m, min(m.cursor.row, m.data.list.Info.NumPending))
	}
	for i, task := range newTasks {
		pasteIdx := min(taskIndex+i+1, len(m.data.list.TaskIds))
		if before {
			pasteIdx = taskIndex + i
		}
		err := m.data.list.Insert(task, pasteIdx)
		if err != nil {
			panic(err)
		}
	}

	// fix cursor
	if before {
		m.cursor.row -= len(newTasks) - 1
	} else {
		m.cursor.row += len(newTasks)
	}

	m.editInfo.dirty = true
	return m
}

// find the index of the task with the given id
func getTaskIndex(m model, displayIdx int) int {
	id := getTaskId(m, displayIdx)
	for i, taskId := range m.data.list.TaskIds {
		if taskId == id {
			return i
		}
	}
	return -1 // should never happen, but just in case
}
