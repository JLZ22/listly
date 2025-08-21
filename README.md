# listly

`listly` is a CLI utility for organizing todo lists with Vim-like keybindings. You can create and manage different lists with CLI commands and make changes to tasks within them with Vim-style keybindings in a TUI. `listly` takes inspiration from Git's branching with respect to switching between and editing different todo lists (not the VC part).

Each todo list resembles a Git branch: there is always a current list that you're working on (unless you have no lists), similar to how Git always has a current branch checked out. Commands that operate on the “current list” implicitly affect this active list unless another list is explicitly specified.

This design allows seamless switching between multiple task contexts, enabling context-based task management or project-specific lists without losing track of your progress elsewhere. You can create new lists, switch between them, and keep tasks organized across different areas of your work or life while staying efficient with intuitive CLI commands and natural Vim-style keybindings.

## Installation

Run the following command to install `listly` into your `$GOBIN` path, which defaults to `$GOPATH/bin` or `$HOME/go/bin` if the `GOPATH` environment variable is not set.

```bash
go install github.com/jlz22/listly@latest
```

## Usage

![Demo](assets/demo.gif)

[Slightly Longer YouTube Demo](https://youtu.be/s1b4MqS0Fhg)

### CLI

| Command                                        | Description                                                                                                |
| ---------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| `listly open [list name]`                      | Open the specified list in the TUI, and switch current list to it. Open current list if no list specified. |
| `listly new <list name> [other list names...]` | Create a new list with the specified name(s).                                                              |
| `listly switch <list name>`                    | Switch to the specified list in the TUI.                                                                   |
| `listly show [list name]`                      | Print info about the specified list and all tasks in it. Show current list if no list specified.           |
| `listly list`                                  | Print name of all lists and their task counts.                                                             |
| `listly clean [list name] ...`                 | Remove all completed tasks from the specified list(s). Clean current list if no list(s) specified.         |
| `listly clean -a, --all`                       | Remove all completed tasks from all lists.                                                                 |
| `listly rename <old name> <new name>`          | Rename a list from <old name> to <new name>                                                                |
| `listly delete <list name>`                    | Delete the specified list(s) - will ignore lists that do not exist.                                        |
| `listly import <file>`                         | Import tasks from a file. Supported formats: JSON, YAML.                                                   |
| `listly export <file> [list names...]`         | Export list(s) to a file. Exports current list if no list name specified. Supported formats: JSON, YAML.   |
| `listly auth`                                  | Add Google Gemini API key.                                                                                 |
| `listly generate <file>`                       | Generate todo lists from a prompt in a text file.                                                          |

### TUI Controls

| Key       | Action                                                             |
| --------- | ------------------------------------------------------------------ |
| `j, up`   | Move down                                                          |
| `k, down` | Move up                                                            |
| `n`       | Create a new task                                                  |
| `i`       | Edit current task                                                  |
| `d`       | Delete the current task and copy it                                |
| `space`   | Toggle a task as done or not done                                  |
| `v`       | Toggle visual mode                                                 |
| `d`       | Delete the selection and copy it                                   |
| `y`       | Copy the selected item(s) or current item if no selection          |
| `p`       | Paste the copied item(s) after the current item                    |
| `P`       | Paste the copied item(s) before the current item                   |
| `w`       | Save changes                                                       |
| `q`       | Quit the application - discard all changes, requiring confirmation |
| `{`       | Jump up                                                            |
| `}`       | Jump down                                                          |
| `o`       | New task after the cursor                                          |
| `O`       | New task before the cursor                                         |

### Getting a Gemini API Key

To use the `generate` command, you need to set up a Gemini API key. You can get one by following these steps:

1. Go to [Google AI Studio](https://aistudio.google.com/apikey).
2. Click on "Create API Key".
3. If prompted, select or create a project.
4. Copy the generated API key.
5. Run `listly auth` and follow the prompts to set your API key.

## Quirks / Issues

- TUI renders inconsistently when run on MacOS terminal as opposed to iTerm.
