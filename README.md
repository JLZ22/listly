# listly

`listly` is a CLI utility for organizing todo lists with Vim-like keybindings. You can create and manage different lists with CLI commands and make changes to tasks within them with Vim-style keybindings in a TUI. `listly` takes inspiration from Git's branching with respect to switching between and editing different todo lists (not the VC part).

Each todo list resembles a Git branch: there is always a current list that you're working on (unless you have no lists), similar to how Git always has a current branch checked out. Commands that operate on the “current list” implicitly affect this active list unless another list is explicitly specified.

This design allows seamless switching between multiple task contexts, enabling project-specific lists without losing track of your progress elsewhere. You can create new lists, switch between them, and keep tasks organized across different areas of your work or life while staying efficient with intuitive CLI commands and natural Vim-style keybindings.

`listly` also includes a Google Gemini powered list generation feature that can create task lists based on a project description passed in via a text file. After the fact, I independently discovered that this features bears extremely close resemblance to Claude Code's "plan mode". 

## Table of Contents

- [Installation](#installation)
- [Demos](#demos)
  - [Basic Functionality](#basic-functionality)
  - [AI Powered List Generation](#ai-powered-list-generation-demo)
- [Usage](#usage)
  - [CLI](#cli)
  - [TUI Controls](#tui-controls)
  - [Getting a Gemini API Key](#getting-a-gemini-api-key)
- [Quirks / Issues](#quirks--issues)

## Demos

### Basic Functionality

![Demo](assets/demo.gif)

[Slightly Longer YouTube Demo](https://youtu.be/s1b4MqS0Fhg)

### AI Powered List Generation

![Gemini List Gen Demo](assets/generate_demo.gif)

[Youtube Version](https://youtu.be/Fy3c-LyJLRI)

## Installation

Run the following command to install `listly` into your `$GOBIN` path, which defaults to `$GOPATH/bin` or `$HOME/go/bin` if the `GOPATH` environment variable is not set.

```bash
go install github.com/jlz22/listly@latest
```

## Usage

### CLI

| Command                                        | Description                                                                                                |
| ---------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| `listly open [list name]`                      | Open the specified list in the TUI, and switch current list to it. Open current list if no list specified. |
| `listly new <list name> [other list names...]` | Create new list(s) with the specified name(s).                                                             |
| `listly switch <list name>`                    | Switch to the specified list in the TUI.                                                                   |
| `listly show [list name]`                      | Print info about the specified list and all tasks in it. Show current list if no list specified.           |
| `listly list`                                  | Print name of all lists and their task counts.                                                             |
| `listly clean [list names...]`                 | Remove all completed tasks from the specified list(s). Clean current list if no list(s) specified.         |
| `listly clean -a, --all`                       | Remove all completed tasks from all lists.                                                                 |
| `listly rename <old name> <new name>`          | Rename a list from <old name> to <new name>                                                                |
| `listly delete <list name>`                    | Delete the specified list(s) - will ignore lists that do not exist.                                        |
| `listly import <file>`                         | Import tasks from a file. Supported formats: JSON, YAML.                                                   |
| `listly export <file> [list names...]`         | Export list(s) to a file. Exports current list if no list name specified. Supported formats: JSON, YAML.   |
| `listly auth`                                  | Add Google Gemini API key.                                                                                 |
| `listly generate <file>`                       | Generate todo lists from a prompt in a text file.                                                          |
| `listly kmap set <file>` | Stores the specified file path as Listly’s custom key-map and automatically loads it on every run. |
| `listly kmap clear` | Removes the specified file path, reverting Listly to the default key-map. |
| `listly kmap show` | Outputs the file path that is stored for Listly's custom key-map. |

### TUI

#### Default Bindings

| Action                                                             | Official Name    | Mode                            | Key      |
| ------------------------------------------------------------------ | ---------------- | ------------------------------- | -------- |
| Move down                                                          | Down             | Shared - Normal, Visual         | `j`      |
| Move down 5 rows                                                   | DownFive         | Shared - Normal, Visual         | `J`      |
| Move up                                                            | Up               | Shared - Normal, Visual         | `k`      |
| Move up 5 rows                                                     | UpFive           | Shared - Normal, Visual         | `K`      |
| Create a new task                                                  | NewTask          | Normal                          | `n`      |
| Edit current task                                                  | EditTask         | Normal                          | `i`      |
| Delete the current task and copy it                                | DeleteTask       | Normal                          | `d`      |
| Delete selection in visual mode                                    | Delete           | Visual                          | `d`      |
| Clear and edit current task                                        | ClearAndEdit     | Normal                          | `x`      |
| Toggle a task as done or not done                                  | ToggleCompletion | Shared - Normal, Visual         | `space`  |
| Toggle visual mode                                                 | EnableVisualMode | Normal                          | `v`      |
| Copy the selected item(s) or current item if no selection          | Yank             | Shared - Normal, Visual         | `y`      |
| Paste the copied item(s) after the current item                    | PasteAfter       | Normal                          | `p`      |
| Paste the copied item(s) before the current item                   | PasteBefore      | Normal                          | `P`      |
| Save changes                                                       | Write            | Normal                          | `w`      |
| Quit the application - discard all changes, requiring confirmation | QuitWithWarning  | Normal                          | `q`      |
| Quit without confirmation                                          | QuitNoWarning    | Shared - Normal, Visual, Insert | `ctrl+c` |
| Jump up                                                            | JumpUp           | Shared - Normal, Visual         | `{`      |
| Jump down                                                          | JumpDown         | Shared - Normal, Visual         | `}`      |
| New task after the cursor                                          | NewAfter         | Normal                          | `o`      |
| New task before the cursor                                         | NewBefore        | Normal                          | `O`      |
| Discard changes                                                    | Discard          | Insert                          | `esc`    |
| Save changes in insert mode                                        | Save             | Insert                          | `enter`  |
| Back to normal mode.                                               | NormalMode       | Visual                          | `esc`    |

#### Custom Bindings

To import your own custom key-binds, you can use 

```
listly kmap set <file>
```

. The file **MUST** be a `.yaml` file that is formatted as `./assets/default_kmap.yaml` is. It **IS** case sensitive. Any commands (e.g. `QuitWithWarning`) that are not specified in your config file will be replaced with the default **UNLESS** that would create a duplicate binding in which case Listly will give you an error. Any commands that are not included in the "Official Name" column (e.g. `Quit`) will be ignored. 

Note: `./assets/default_kmap.yaml` is just an example for you. The defaults will not be changed if you modify this file. `./assets/toy_kmap.yaml` is an alternate mapping where many commands have swapped key-binds. This was created for fun and is not recommended for actual use. 

### Getting a Gemini API Key

To use the `generate` command, you need to set up a Gemini API key. You can get one by following these steps:

1. Go to [Google AI Studio](https://aistudio.google.com/apikey).
2. Click on "Create API Key".
3. If prompted, select or create a project.
4. Copy the generated API key.
5. Run `listly auth` and follow the prompts to set your API key.

## Quirks / Issues

- TUI renders inconsistently when run on MacOS terminal as opposed to iTerm.
