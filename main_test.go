package main

import (
	"bytes"
	"testing"
	"todo-cli/todo"
)

type fakeRepository struct {
	todos     []todo.Todo
	loadError error
	saveError error
}

func (r *fakeRepository) Load() ([]todo.Todo, error) {
	if r.loadError != nil {
		return nil, r.loadError
	}
	return r.todos, nil
}

func (r *fakeRepository) Save(todos []todo.Todo) error {
	if r.saveError != nil {
		return r.saveError
	}
	r.todos = todos
	return nil
}

func TestPrintTodos(t *testing.T) {
	tests := []struct {
		name     string
		todos    []todo.Todo
		expected string
	}{
		{
			name:     "empty todo",
			todos:    []todo.Todo{},
			expected: "",
		},
		{
			name:     "single unfinished",
			todos:    []todo.Todo{{ID: 1, Title: "Learn Go", Done: false}},
			expected: "1. [ ] Learn Go\n",
		},
		{
			name:     "single done",
			todos:    []todo.Todo{{ID: 1, Title: "Learn Go", Done: true}},
			expected: "1. [x] Learn Go\n",
		},
		{
			name: "multiple todos",
			todos: []todo.Todo{
				{ID: 1, Title: "Learn Go", Done: false},
				{ID: 2, Title: "Learn Testing", Done: true},
			},
			expected: "1. [ ] Learn Go\n2. [x] Learn Testing\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			printTodos(&buf, tt.todos)

			got := buf.String()
			if tt.expected != got {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestParseCommand(t *testing.T) {
	args := []string{"list"}
	command, err := parseCommand(args)
	if err != nil {
		t.Fatal(err)
	}

	if command != "list" {
		t.Errorf("expected 'list', got %s", command)
	}
}

func TestParseCommand_Failed(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		errorMessage string
	}{
		{
			name:         "missing command",
			args:         []string{},
			errorMessage: "missing command",
		},
		{
			name:         "unknown command",
			args:         []string{"unknown"},
			errorMessage: "unknown command",
		},
		{
			name:         "argument error",
			args:         []string{"list", "extra"},
			errorMessage: "argument error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := parseCommand(tt.args)
			if err == nil {
				t.Fatalf("expected error with message %s", tt.errorMessage)
			}

			if err.Error() != tt.errorMessage {
				t.Errorf("expected error message: %s, got %s", tt.errorMessage, err.Error())
			}

			if command != "" {
				t.Errorf("expected empty string, got %s", command)
			}
		})
	}
}

func TestRunCLI(t *testing.T) {
	todos := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
	}
	fakeRepo := &fakeRepository{
		todos: todos,
	}
	service, err := todo.NewServiceWithRepository(fakeRepo)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer

	err = runCLI("list", service, &buf)
	if err != nil {
		t.Fatal(err)
	}

	wanted := "1. [ ] Learn Go\n"
	if wanted != buf.String() {
		t.Errorf("expected %v, got %v", wanted, buf.String())
	}
}
