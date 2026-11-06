package todo

import (
	"errors"
	"slices"
)

type Service struct {
	todos []Todo
	repo  TodoRepository
}

func NewService(initial []Todo) *Service {
	return &Service{todos: slices.Clone(initial)}
}

func NewServiceWithRepository(repo TodoRepository) (*Service, error) {
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}

	todos, err := repo.Load()
	if err != nil {
		return nil, err
	}
	return &Service{todos: slices.Clone(todos), repo: repo}, nil
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

func (s *Service) Delete(id int) error {
	updated, err := DeleteTodo(s.todos, id)
	if err != nil {
		return err
	}
	s.todos = updated
	return nil
}

func (s *Service) Search(keyword string) []Todo {
	return SearchTodos(s.todos, keyword)
}

func (s *Service) Save() error {
	if s.repo == nil {
		return errors.New("this service does not have repository")
	}

	err := s.repo.Save(slices.Clone(s.todos))
	if err != nil {
		return err
	}
	return nil
}
