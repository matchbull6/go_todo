# go_todo

A terminal task manager in Go. Tasks are stored in a JSON file.

## Usage

```sh
go_todo add "buy milk"                  # add a task
go_todo add "fix the bug" -p high       # flags can go anywhere
go_todo list                            # pending tasks
go_todo list --all                      # including done
go_todo done 3                          # mark task 3 as done
go_todo rm 3                            # delete task 3
go_todo --file work.json list           # use another file
go_todo help add                        # help for a command
```

The task file is chosen in this order:

1. the `--file` flag;
2. the `TODO_FILE` environment variable;
3. `.go_todo.json` in your home directory.

## Roadmap

- [x] Task model and JSON storage
- [x] Atomic save
- [x] Commands: `add`, `list`, `done`, `rm`
- [x] Flags: `--all`, `--priority`
- [x] Aligned output
- [ ] Task file choice: `--file`, `TODO_FILE`
- [ ] Tests
- [ ] Extras: `edit`, `clear`, due dates

## Build

```sh
go build -o go_todo.exe .
```
