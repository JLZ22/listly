package tui

import (
	"strings"

	help "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jlz22/listly/core"
)

type confirmation struct {
	active  bool
	message string
}

type data struct {
	list core.List
	db   *core.DB
}

type cursor struct {
	row      int
	selStart int
}

type editInfo struct {
	textInput textinput.Model
	copyBuff  []core.Task // buffer for copied tasks
	dirty     bool
	taskId    int // the id of the task being edited
	location  int // where to insert the new task
}

type model struct {
	data         data
	cursor       cursor
	editInfo     editInfo
	confirmation confirmation
	mode         string
	vp           viewport.Model
	kmap         KeyMap
}

var titleStyle = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return lipgloss.NewStyle().BorderStyle(b).Padding(0, 0)
}()

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

	// initialize text input
	ti := textinput.New()
	ti.Placeholder = "Task Description"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40
	ti.Prompt = "    > [ ] "

	// load key-mappings
	pth, _ := db.GetKmapPath() // can ignore error here because LoadKmap will use defaults with bad path
	kmap, err := LoadKmap(pth)
	if err != nil {
		return model{}, err
	}

	return model{
		data: data{
			list: list,
			db:   db,
		},
		cursor: cursor{
			row:      0,
			selStart: -1,
		},
		editInfo: editInfo{
			textInput: ti,
			copyBuff:  make([]core.Task, 0),
			dirty:     false,
			taskId:    -1,
		},
		confirmation: confirmation{
			active:  false,
			message: "",
		},
		mode: "normal",
		vp:   viewport.New(0, 0),
		kmap: kmap,
	}, nil
}

func (m *model) ensureCursorVisible() {
	top := m.vp.YOffset
	bottom := top + m.vp.Height - 1
	if m.cursor.row < top {
		m.vp.YOffset = m.cursor.row - 1
	} else if m.cursor.row > bottom-8 {
		m.vp.YOffset = m.cursor.row - m.vp.Height + 9
	}
}

// Interpret input and update shared state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(makeHeader(m))
		footerHeight := lipgloss.Height(makeFooter(m, help.New().FullHelpView(DefaultNormalKeyMap.FullHelp())))
		verticalHeight := headerHeight + footerHeight
		m.vp.Width, m.vp.Height = msg.Width, msg.Height-verticalHeight

	case tea.KeyMsg:
		switch m.mode {
		case "normal":
			m, cmd = handleNormalInput(msg, m)
		case "insert":
			m, cmd = handleInsertInput(msg, m)
		case "visual":
			m, cmd = handleVisualInput(msg, m)
		}
	}

	// keep cursor in view & update content
	m.ensureCursorVisible()
	switch m.mode {
	case "normal":
		m.vp.SetContent(renderNormalView(m) + "\nEOF")
	case "insert":
		m.vp.SetContent(renderInsertView(m) + "\nEOF")
	case "visual":
		m.vp.SetContent(renderVisualView(m) + "\nEOF")
	}

	return m, cmd
}

func (m model) View() string {
	if m.confirmation.active {
		return "\n" + m.confirmation.message
	}
	out := makeHeader(m) + m.vp.View() + "\n\n"

	switch m.mode {
	case "normal":
		return out + makeFooter(m, help.New().FullHelpView(DefaultNormalKeyMap.FullHelp()))
	case "insert":
		return out + makeFooter(m, help.New().FullHelpView(DefaultInsertKeyMap.FullHelp()))
	case "visual":
		return out + makeFooter(m, help.New().FullHelpView(DefaultVisualKeyMap.FullHelp()))
	}
	return "Developer is a monkey - this shouldn't happen"
}

// convert the list into a bunch of lines
func buildLines(m model, includeCursor bool) []string {
	if m.data.list.Info.NumTasks == 0 {
		return []string{"\n No tasks in this list. Press \"n\" to add one.\n\n"}
	}
	lines := make([]string, 0, 1+m.data.list.Info.NumTasks+1) // + 1 for the dividing bar
	lines = append(lines, "\n  Todo:\n\n")

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
		lines = append(lines, "\n  Complete:\n\n") // add a blank line between pending and completed tasks
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

func makeHeader(m model) string {
	listName := m.data.list.Info.Name

	// mark dirty list with (*)
	if m.editInfo.dirty {
		listName += " (*)"
	}

	title := titleStyle.Render(listName)
	line := strings.Repeat("─", max(0, m.vp.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line) + "\n"
}

func makeFooter(m model, help string) string {
	line := strings.Repeat("─", max(0, m.vp.Width))
	return lipgloss.JoinVertical(lipgloss.Center, line, help)
}
