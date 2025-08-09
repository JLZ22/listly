# tidy

`tidy` is a CLI utility for managing todo lists with vim-like keybindings. You can create, edit, and delete todo lists with CLI commands and make changes to tasks within those lists with vim-style keybindings. `tidy` takes inspiration from Git's branching with respect to switching between and editing different todo lists. 

Each todo list resembles a Git branch: there is always a current list that you're working on (unless you have no lists), similar to how Git always has a current branch checked out. Commands that operate on the “current list” implicitly affect this active list unless another list is explicitly specified. This design allows seamless switching between multiple task contexts, enabling workflows like context-based task management or project-specific lists without losing track of your progress. You can create new lists, switch between them, and keep tasks organized across different areas of your work or life, all while maintaining a familiar and efficient Vim-inspired interface.

## Installation

WIP

## Usage

### CLI

| Command                             | Description                                                                                          |
| ----------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `tidy open`                         | Open the TUI for the current list. If not available, create a new list named "untitled".             |
| `tidy open <list name>`             | Open the specified list in the TUI, and switch current list to it - fails if the list does not exist |
| `tidy new <list name>`              | Create a new list with the specified name - fails if the list already exists                         |
| `tidy switch <list name>`           | Switch to the specified list in the TUI - fails if the list does not exist                           |
| `tidy show`                         | Print info about the current list and all tasks in it                                                |
| `tidy show <list name>`             | Print info about the specified list and all tasks in it                                              |
| `tidy list`                         | Print name of all lists and their task counts                                                        |
| `tidy clean`                        | Remove all completed tasks from the current list                                                     |
| `tidy clean <list name>`            | Remove all completed tasks from the specified list                                                   |
| `tidy clean -a, --all`              | Remove all completed tasks from all lists                                                            |
| `tidy rename <old name> <new name>` | Rename a list from <old name> to <new name> - fails if the new name already exists                   |

### TUI Controls

| Key            | Action                                                             |
| -------------- | ------------------------------------------------------------------ |
| `j, up`        | Move down                                                          |
| `k, down`      | Move up                                                            |
| `h, left`      | Move left                                                          |
| `l, right`     | Move right                                                         |
| `n`            | Create a new item                                                  |
| `i`            | Edit current item                                                  |
| `D`            | Delete the current item and copy it                                |
| `space, Enter` | Toggle an item as done or not done                                 |
| `v`            | Toggle select mode                                                 |
| `d`            | Delete the selection and copy it                                   |
| `y`            | Copy the selected item(s) or current item if no selection          |
| `p`            | Paste the copied item(s) after the current item                    |
| `Esc`          | Quit the TUI or go back                                            |
| `w`            | Save changes                                                       |
| `q`            | Quit the application - discard all changes, requiring confirmation |
