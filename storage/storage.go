package storage

import (
	"encoding/json"
	"os"
	"todo-cli/todo"
)

func LoadTodos() []todo.Todo {
	data, err := os.ReadFile("todos.json")
	if err != nil {
		return []todo.Todo{}
	}

	var todos []todo.Todo

	err = json.Unmarshal(data, &todos)

	if err != nil {
		return []todo.Todo{}
	}

	return todos
}

func SaveTodos(todos []todo.Todo) error {
	data, err := json.Marshal(todos)
	if err != nil {
		return err
	}

	err = os.WriteFile("todos.json", data, 0644)

	if err != nil {
		return err
	}

	return nil
}
