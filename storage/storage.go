package storage

import (
	"encoding/json"
	"errors"
	"fmt"
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

		return nil, fmt.Errorf("read %s failed: %w", filename, err)
	}

	var todos []todo.Todo

	err = json.Unmarshal(data, &todos)

	if err != nil {
		return nil, fmt.Errorf("parse %s failed: %w", filename, err)
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
