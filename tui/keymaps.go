package tui

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	key "github.com/charmbracelet/bubbles/key"
)

var DefaultKeyMapConfig = map[string]map[string]string{
	"Shared": {
		"Up":               "k",
		"UpFive":           "K",
		"Down":             "j",
		"DownFive":         "J",
		"QuitNoWarning":    "ctrl+c",
		"Yank":             "y",
		"ToggleCompletion": " ",
		"JumpUp":           "{",
		"JumpDown":         "}",
	},
	"Normal": {
		"QuitWithWarning":  "q",
		"NewTask":          "n",
		"NewBefore":        "O",
		"NewAfter":         "o",
		"EditTask":         "i",
		"ClearAndEdit":     "x",
		"DeleteTask":       "d",
		"EnableVisualMode": "v",
		"PasteAfter":       "p",
		"PasteBefore":      "P",
		"Write":            "w",
	},
	"Insert": {
		"Discard": "esc",
		"Save":    "enter",
	},
	"Visual": {
		"NormalMode": "esc",
		"Delete":     "d",
	},
}

// Merges shared into specific, leaving specific's values in case of conflict.
func mergeKeys(shared map[string]string, specific map[string]string) map[string]string {
	out := make(map[string]string)

	// Copy shared
	for k, v := range shared {
		out[k] = v
	}

	// Override with specific
	for k, v := range specific {
		out[k] = v
	}

	return out
}

// ------------------------ Normal Mode Keymaps ------------------------

var normalCommands = []string{
	"Up", "UpFive", "Down", "DownFive", "QuitWithWarning", "QuitNoWarning",
	"NewTask", "NewBefore", "NewAfter", "EditTask", "ClearAndEdit", "DeleteTask",
	"ToggleCompletion", "EnableVisualMode", "Yank", "PasteAfter", "PasteBefore",
	"Write", "JumpUp", "JumpDown",
}

type NormalKeyMap struct {
	Up               key.Binding
	UpFive           key.Binding
	Down             key.Binding
	DownFive         key.Binding
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
		{k.Up, k.Down, k.QuitWithWarning},                 // first column
		{k.Write, k.NewTask, k.EditTask},                  // second column
		{k.ClearAndEdit, k.DeleteTask, k.Yank},            // third column
		{k.EnableVisualMode, k.PasteAfter, k.PasteBefore}, // fourth column
		{k.JumpUp, k.JumpDown, k.ToggleCompletion},        // fifth column
		{k.NewBefore, k.NewAfter},
	}
}

// This function builds a NormalKeyMap from a config map under the naive assumption that
// all keys are present and valid.
func buildNormalKmapFromConfig(config map[string]string) (NormalKeyMap, error) {
	// validate the config
	for _, cmd := range normalCommands {
		if _, ok := config[cmd]; !ok {
			return NormalKeyMap{}, fmt.Errorf("missing key binding for command: %s", cmd)
		}
	}

	// Build the NormalKeyMap
	return NormalKeyMap{
		Up: key.NewBinding(
			key.WithKeys(config["Up"]),
			key.WithHelp(config["Up"], "up"),
		),
		Down: key.NewBinding(
			key.WithKeys(config["Down"]),
			key.WithHelp(config["Down"], "down"),
		),
		UpFive: key.NewBinding(
			key.WithKeys(config["UpFive"]),
			key.WithHelp(config["UpFive"], "up 5"),
		),
		DownFive: key.NewBinding(
			key.WithKeys(config["DownFive"]),
			key.WithHelp(config["DownFive"], "down 5"),
		),
		QuitWithWarning: key.NewBinding(
			key.WithKeys(config["QuitWithWarning"]),
			key.WithHelp(config["QuitWithWarning"], "quit"),
		),
		QuitNoWarning: key.NewBinding(
			key.WithKeys(config["QuitNoWarning"]),
		),
		NewTask: key.NewBinding(
			key.WithKeys(config["NewTask"]),
			key.WithHelp(config["NewTask"], "new task"),
		),
		NewBefore: key.NewBinding(
			key.WithKeys(config["NewBefore"]),
			key.WithHelp(config["NewBefore"], "new task before"),
		),
		NewAfter: key.NewBinding(
			key.WithKeys(config["NewAfter"]),
			key.WithHelp(config["NewAfter"], "new task after"),
		),
		EditTask: key.NewBinding(
			key.WithKeys(config["EditTask"]),
			key.WithHelp(config["EditTask"], "edit task"),
		),
		ClearAndEdit: key.NewBinding(
			key.WithKeys(config["ClearAndEdit"]),
			key.WithHelp(config["ClearAndEdit"], "clear and edit"),
		),
		DeleteTask: key.NewBinding(
			key.WithKeys(config["DeleteTask"]),
			key.WithHelp(config["DeleteTask"], "cut task"),
		),
		ToggleCompletion: key.NewBinding(
			key.WithKeys(config["ToggleCompletion"]),
			key.WithHelp(config["ToggleCompletion"], "mark done/not done"),
		),
		EnableVisualMode: key.NewBinding(
			key.WithKeys(config["EnableVisualMode"]),
			key.WithHelp(config["EnableVisualMode"], "visual mode"),
		),
		Yank: key.NewBinding(
			key.WithKeys(config["Yank"]),
			key.WithHelp(config["Yank"], "yank"),
		),
		PasteAfter: key.NewBinding(
			key.WithKeys(config["PasteAfter"]),
			key.WithHelp(config["PasteAfter"], "paste"),
		),
		PasteBefore: key.NewBinding(
			key.WithKeys(config["PasteBefore"]),
			key.WithHelp(config["PasteBefore"], "paste before"),
		),
		Write: key.NewBinding(
			key.WithKeys(config["Write"]),
			key.WithHelp(config["Write"], "write"),
		),
		JumpUp: key.NewBinding(
			key.WithKeys(config["JumpUp"]),
			key.WithHelp(config["JumpUp"], "jump up"),
		),
		JumpDown: key.NewBinding(
			key.WithKeys(config["JumpDown"]),
			key.WithHelp(config["JumpDown"], "jump down"),
		),
	}, nil
}

var DefaultNormalKeyMap, _ = buildNormalKmapFromConfig(mergeKeys(DefaultKeyMapConfig["Shared"], DefaultKeyMapConfig["Normal"]))

// ------------------------ Insert Mode Keymaps ------------------------

var insertCommands = []string{
	"Discard", "QuitNoWarning", "Save",
}

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

// This function builds a InsertKeyMap from a config map under the naive assumption that
// all keys are present and valid.
func buildInsertKmapFromConfig(config map[string]string) (InsertKeyMap, error) {
	// validate the config
	for _, cmd := range insertCommands {
		if _, ok := config[cmd]; !ok {
			return InsertKeyMap{}, fmt.Errorf("missing key binding for command: %s", cmd)
		}
	}

	return InsertKeyMap{
		Discard: key.NewBinding(
			key.WithKeys(config["Discard"]),
			key.WithHelp(config["Discard"], "discard changes"),
		),
		QuitNoWarning: key.NewBinding(
			key.WithKeys(config["QuitNoWarning"]),
		),
		Save: key.NewBinding(
			key.WithKeys(config["Save"]),
			key.WithHelp(config["Save"], "save"),
		),
	}, nil
}

var DefaultInsertKeyMap, _ = buildInsertKmapFromConfig(mergeKeys(DefaultKeyMapConfig["Shared"], DefaultKeyMapConfig["Insert"]))

// ------------------------ Visual Mode Keymaps ------------------------

var visualCommands = []string{
	"Up", "UpFive", "Down", "DownFive", "NormalMode", "QuitNoWarning",
	"Delete", "Yank", "ToggleCompletion", "JumpUp", "JumpDown",
}

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

// This function builds a VisualKeyMap from a config map under the naive assumption that
// all keys are present and valid.
func buildVisualKmapFromConfig(config map[string]string) (VisualKeyMap, error) {
	// validate the config
	for _, cmd := range visualCommands {
		if _, ok := config[cmd]; !ok {
			return VisualKeyMap{}, fmt.Errorf("missing key binding for command: %s", cmd)
		}
	}

	return VisualKeyMap{
		Up: key.NewBinding(
			key.WithKeys(config["Up"]),
			key.WithHelp(config["Up"], "up"),
		),
		Down: key.NewBinding(
			key.WithKeys(config["Down"]),
			key.WithHelp(config["Down"], "down"),
		),
		UpFive: key.NewBinding(
			key.WithKeys(config["UpFive"]),
			key.WithHelp(config["UpFive"], "up 5"),
		),
		DownFive: key.NewBinding(
			key.WithKeys(config["DownFive"]),
			key.WithHelp(config["DownFive"], "down 5"),
		),
		NormalMode: key.NewBinding(
			key.WithKeys(config["NormalMode"]),
			key.WithHelp(config["NormalMode"], "normal mode"),
		),
		QuitNoWarning: key.NewBinding(
			key.WithKeys(config["QuitNoWarning"]),
		),
		Delete: key.NewBinding(
			key.WithKeys(config["Delete"]),
			key.WithHelp(config["Delete"], "cut"),
		),
		Yank: key.NewBinding(
			key.WithKeys(config["Yank"]),
			key.WithHelp(config["Yank"], "yank"),
		),
		ToggleCompletion: key.NewBinding(
			key.WithKeys(config["ToggleCompletion"]),
			key.WithHelp(config["ToggleCompletion"], "mark done/not done"),
		),
		JumpUp: key.NewBinding(
			key.WithKeys(config["JumpUp"]),
			key.WithHelp(config["JumpUp"], "jump up"),
		),
		JumpDown: key.NewBinding(
			key.WithKeys(config["JumpDown"]),
			key.WithHelp(config["JumpDown"], "jump down"),
		),
	}, nil
}

var DefaultVisualKeyMap, _ = buildVisualKmapFromConfig(mergeKeys(DefaultKeyMapConfig["Shared"], DefaultKeyMapConfig["Visual"]))

// ---------------------------------- TUI Keymap -------------------

type KeyMap struct {
	Normal NormalKeyMap
	Insert InsertKeyMap
	Visual VisualKeyMap
}

var DefaultKeyMap = KeyMap{
	Normal: DefaultNormalKeyMap,
	Insert: DefaultInsertKeyMap,
	Visual: DefaultVisualKeyMap,
}

func LoadKmap(pth string) (KeyMap, error) {
	// Read file
	data, err := os.ReadFile(pth)
	if err != nil {
		return DefaultKeyMap, nil
	}

	// Load yaml to map
	var config map[string]map[string]string
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return DefaultKeyMap, nil
	}

	// Merge shared with mode-specific
	normalKeys := mergeKeys(config["Shared"], config["Normal"])
	insertKeys := mergeKeys(config["Shared"], config["Insert"])
	visualKeys := mergeKeys(config["Shared"], config["Visual"])

	// Merge mode with default to fill in missing keys
	normalKeys = mergeKeys(DefaultKeyMapConfig["Normal"], normalKeys)
	insertKeys = mergeKeys(DefaultKeyMapConfig["Insert"], insertKeys)
	visualKeys = mergeKeys(DefaultKeyMapConfig["Visual"], visualKeys)

	// Check if any keys are overlapping within a mode
	err = checkConflicts(normalKeys, "normal")
	if err != nil {
		return KeyMap{}, err
	}
	err = checkConflicts(insertKeys, "insert")
	if err != nil {
		return KeyMap{}, err
	}
	err = checkConflicts(visualKeys, "visual")
	if err != nil {
		return KeyMap{}, err
	}

	// Populate key-map for each mode
	normalKmap, err := buildNormalKmapFromConfig(normalKeys)
	if err != nil {
		return DefaultKeyMap, nil
	}
	insertKmap, err := buildInsertKmapFromConfig(insertKeys)
	if err != nil {
		return DefaultKeyMap, nil
	}
	visualKmap, err := buildVisualKmapFromConfig(visualKeys)
	if err != nil {
		return DefaultKeyMap, nil
	}

	return KeyMap{
		Normal: normalKmap,
		Insert: insertKmap,
		Visual: visualKmap,
	}, nil
}

func checkConflicts(cfg map[string]string, mode string) error {
	seen := make(map[string]string)
	for cmd, keyStr := range cfg {
		if seen[keyStr] != "" {
			return fmt.Errorf("conflicting key binding in %s mode: %s and %s both bound to '%s'.\n", mode, cmd, seen[keyStr], keyStr)
		} else {
			seen[keyStr] = cmd
		}
	}
	return nil
}
