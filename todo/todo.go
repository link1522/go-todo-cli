package todo

import (
	"errors"
	"fmt"
	"strings"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (todo *Todo) MarkDone() {
	todo.Done = true
}

func AddTodo(todos []Todo, title string, id int) []Todo {
	return append(todos, Todo{
		ID:    id,
		Title: title,
		Done:  false,
	})
}

func CompleteTodo(todos []Todo, id int) bool {
	for i, item := range todos {
		if id == item.ID {
			todos[i].MarkDone()
			return true
		}
	}
	return false
}

func DeleteTodo(todos []Todo, id int) ([]Todo, error) {
	for i, item := range todos {
		if id == item.ID {
			return append(todos[:i], todos[i+1:]...), nil
		}
	}
	return todos, &TodoNotFoundError{ID: id}
}

func GetNextID(todos []Todo) int {
	maxID := 0
	for _, item := range todos {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	return maxID + 1
}

func SearchTodos(todos []Todo, keyword string) []Todo {
	foundTodos := []Todo{}
	keyword = strings.ToLower(keyword)

	for _, item := range todos {
		title := strings.ToLower(item.Title)

		if strings.Contains(title, keyword) {
			foundTodos = append(foundTodos, item)
		}
	}

	return foundTodos
}

var ErrTodoNotFound = errors.New("todo not found")

type TodoNotFoundError struct {
	ID int
}

func (e *TodoNotFoundError) Error() string {
	return fmt.Sprintf("todo with ID %d not found", e.ID)
}

func (e *TodoNotFoundError) Unwrap() error {
	return ErrTodoNotFound
}
