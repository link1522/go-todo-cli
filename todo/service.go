package todo

import "slices"

type Service struct {
	todos []Todo
}

func NewService(initial []Todo) *Service {
	return &Service{todos: slices.Clone(initial)}
}

func (s *Service) Add(title string) Todo {
	id := GetNextID(s.todos)
	todo := NewTodo(id, title)
	s.todos = append(s.todos, todo)
	return todo
}

func (s *Service) List() []Todo {
	return slices.Clone(s.todos)
}

func (s *Service) Complete(id int) error {
	found := CompleteTodo(s.todos, id)
	if !found {
		return &TodoNotFoundError{ID: id}
	}
	return nil
}
