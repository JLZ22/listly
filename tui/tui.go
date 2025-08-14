package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jlz22/listly/core"
)

type confirmation struct {
	active  bool
	message string
}

type data struct {
	list core.List
	db  *core.DB
}

type cursor struct {
	row int
	selected map[int]struct{}
}

type editInfo struct {
	textInput textinput.Model
	copyBuff  []core.Task // buffer for copied tasks
	dirty bool
	taskId int // the id of the task being edited
	location int // where to insert the new task
}

type model struct {
	data data
	cursor cursor
	editInfo editInfo
	confirmation confirmation
	mode         string
	Error error
}

func (m model) Init() tea.Cmd {
	// No init I/O needed.
	return nil
}

// Initialize a new model to edit the list with the given name.
func NewModel(db *core.DB, listName string) (model, error) {
	list, err := db.GetList(listName)
	if err != nil {
		return model{}, err
	}

	ti := textinput.New()
	ti.Placeholder = "Task Description"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40
	ti.Prompt = "    > [ ] "

	return model{
		data: data{
			list: list,
			db:   db,
		},
		cursor: cursor{
			row:     0,
			selected: make(map[int]struct{}),
		},
		editInfo: editInfo{
			textInput: ti,
			copyBuff:  make([]core.Task, 0),
			dirty:    false,
			taskId:   -1,
		},
		confirmation: confirmation{
			active:  false,
			message: "",
		},
		mode: "normal",
		Error: nil,
	}, nil
}

// Interpret input and update shared state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.confirmation.active {
	case true:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "enter":
				m.confirmation.active = false
				m.confirmation.message = ""
				return m, tea.Quit
			case "n":
				m.confirmation.active = false
				m.confirmation.message = ""
				return m, nil
			}
		}
	case false:
		switch m.mode {
		case "normal":
			return handleNormalInput(msg, m)
		case "insert":
			return handleInsertInput(msg, m)
		case "visual":
			return handleVisualInput(msg, m)
		}
	}

	return m, nil
}

// Render the view based on current state.
func (m model) View() string {
	switch m.confirmation.active {
	case true:
		return "\n" + m.confirmation.message
	case false:
		switch m.mode {
		case "normal":
			return renderNormalView(m)
		case "insert":
			return renderInsertView(m)
		case "visual":
			return renderVisualView(m)
		}
	}

	return "Invalid Mode - Developer is a monkey."
}

// convert the list into a bunch of lines
func buildLines(m model, includeCursor bool) []string {
	if m.data.list.Info.NumTasks == 0 {
		return []string{"\n  " + m.data.list.Info.Name + " has no tasks. Press \"n\" to add one.\n\n"}
	}

	// render the list name
	lines := make([]string, 0, 1 + m.data.list.Info.NumTasks)
	listName := "\n  " + m.data.list.Info.Name

	// mark dirty list with (*)
	if m.editInfo.dirty {
		listName += " (*)"
	}
	listName += "\n\n"
	lines = append(lines, listName)

	// track the task count to know when to place the cursor
	i := 0

	// render incomplete tasks
	done, notDone := core.SplitByCompletion(m.data.list)
	for _, task := range notDone {
		str := "      [ ] " + task.Description + "\n"
		if i == m.cursor.row && includeCursor {
			str = "    > [ ] " + task.Description + "\n"
		}
		lines = append(lines, str)
		i++
	}

	// render horizontal bar and all of the completed tasks
	if len(done) > 0 {
		lines = append(lines, "\n    ==============================\n\n") // add a blank line between pending and completed tasks
		for _, task := range done {
			str := "      [x] " + task.Description + "\n"
			if i == m.cursor.row && includeCursor {
				str = "    > [x] " + task.Description + "\n"
			}
			lines = append(lines, str)
			i++
		}
	}

	return lines
}

func getNumNotDone(m model) int {
	_, notDone := core.SplitByCompletion(m.data.list)
	return len(notDone)
}

func getNumDone(m model) int {
	done, _ := core.SplitByCompletion(m.data.list)
	return len(done)
}