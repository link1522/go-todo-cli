package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"todo-cli/todo"
)

var _ todo.TodoRepository = (*JSONRepository)(nil)

type JSONRepository struct {
	path string
}

func NewJSONRepository(path string) *JSONRepository {
	return &JSONRepository{path: path}
}

func (r *JSONRepository) Load() ([]todo.Todo, error) {
	data, err := os.ReadFile(r.path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []todo.Todo{}, nil
		}

		return nil, fmt.Errorf("read %s failed: %w", r.path, err)
	}

	var todos []todo.Todo

	err = json.Unmarshal(data, &todos)
	if err != nil {
		return nil, fmt.Errorf("parse %s failed: %w", r.path, err)
	}

	return todos, nil
}

func (r *JSONRepository) Save(todos []todo.Todo) error {
	data, err := json.Marshal(todos)
	if err != nil {
		return fmt.Errorf("marshal todos for %q: %w", r.path, err)
	}

	err = os.WriteFile(r.path, data, 0644)
	if err != nil {
		return fmt.Errorf("write %s failed: %w", r.path, err)
	}

	return nil
}
