package storage

import (
	"encoding/json"
	"errors"
	"os"
	"todo-cli/todo"
)

func LoadTodos() ([]todo.Todo, error) {
	filename := "todos.json"
	data, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []todo.Todo{}, nil
		}

		return nil, err
	}

	var todos []todo.Todo

	err = json.Unmarshal(data, &todos)

	if err != nil {
		return nil, err
	}

	return todos, nil
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
