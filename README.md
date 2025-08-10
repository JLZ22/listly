# listly

`listly` is a CLI utility for organizing todo lists with Vim-like keybindings. You can create and manage different lists with CLI commands and make changes to tasks within them with Vim-style keybindings in a TUI. `listly` takes inspiration from Git's branching with respect to switching between and editing different todo lists (not the VC part).

Each todo list resembles a Git branch: there is always a current list that you're working on (unless you have no lists), similar to how Git always has a current branch checked out. Commands that operate on the “current list” implicitly affect this active list unless another list is explicitly specified.

This design allows seamless switching between multiple task contexts, enabling context-based task management or project-specific lists without losing track of your progress elsewhere. You can create new lists, switch between them, and keep tasks organized across different areas of your work or life while staying efficient with intuitive CLI commands and natural Vim-style keybindings.

## Installation

WIP

## Usage

### CLI

| Command                               | Description                                                                                          |
| ------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `listly open`                         | Open the TUI for the current list - fails if no current list exists.                                 |
| `listly open <list name>`             | Open the specified list in the TUI, and switch current list to it - fails if the list does not exist |
| `listly new <list name>`              | Create a new list with the specified name - fails if the list already exists                         |
| `listly switch <list name>`           | Switch to the specified list in the TUI - fails if the list does not exist                           |
| `listly show`                         | Print info about the current list and all tasks in it                                                |
| `listly show <list name>`             | Print info about the specified list and all tasks in it                                              |
| `listly list`                         | Print name of all lists and their task counts                                                        |
| `listly clean`                        | Remove all completed tasks from the current list                                                     |
| `listly clean <list name>`            | Remove all completed tasks from the specified list                                                   |
| `listly clean -a, --all`              | Remove all completed tasks from all lists                                                            |
| `listly rename <old name> <new name>` | Rename a list from <old name> to <new name> - fails if the new name already exists                   |
| `listly delete <list name>`           | Delete the specified list - fails if the list does not exist - requires confirmation                 |

### TUI Controls

| Key            | Action                                                             |
| -------------- | ------------------------------------------------------------------ |
| `j, up`        | Move down                                                          |
| `k, down`      | Move up                                                            |
| `h, left`      | Move left                                                          |
| `l, right`     | Move right                                                         |
| `n`            | Create a new task                                                  |
| `i`            | Edit current task                                                  |
| `D`            | Delete the current task and copy it                                |
| `space, Enter` | Toggle a task as done or not done                                  |
| `v`            | Toggle select mode                                                 |
| `d`            | Delete the selection and copy it                                   |
| `y`            | Copy the selected item(s) or current item if no selection          |
| `p`            | Paste the copied item(s) after the current item                    |
| `w`            | Save changes                                                       |
| `q`            | Quit the application - discard all changes, requiring confirmation |
