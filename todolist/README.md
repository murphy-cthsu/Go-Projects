# TodoList CLI

This is a simple command-line interface (CLI) application for managing a todo list. It allows you to add, list, delete, and toggle tasks.

## Installation

1. Clone the repository:
    ```sh
    git clone <repository-url>
    cd todolist
    ```

2. Build the application:
    ```sh
    make build
    ```

## Usage

The CLI supports the following commands:

### Add a Task

To add a new task, use the `add` command with the `-title` flag:

```sh
./bin/todo add -title="Your task title"
```

### List Tasks

To list all tasks, use the `list` command:

```sh
./bin/todo list
```

### Delete a Task

To delete a task, use the `delete` command with the `-id` flag:

```sh
./bin/todo delete -id=1
```

### Toggle Task Completion

To toggle the completion status of a task, use the `toggle` command with the `-id` flag:

```sh
./bin/todo toggle -id=1
```

## Makefile Commands

- `make build`: Builds the application.
- `make test`: Runs the tests.

## Database

The application uses SQLite for storing tasks. The database file is named `tasks.db` and is created in the current directory.

## License

This project is licensed under the MIT License.
