package main

import (
	"bytes"
	"errors"
	"slices"
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
	tests := []struct {
		name   string
		args   []string
		wanted CLICommand
	}{
		{
			name:   "list command",
			args:   []string{"list"},
			wanted: CLICommand{Name: "list"},
		},
		{
			name:   "add command",
			args:   []string{"add", "Learn Go"},
			wanted: CLICommand{Name: "add", Args: []string{"Learn Go"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := parseCommand(tt.args)

			if err != nil {
				t.Fatal(err)
			}

			if command.Name != tt.wanted.Name {
				t.Errorf("expected %s, got %s", tt.wanted.Name, command.Name)
			}
			if !slices.Equal(command.Args, tt.wanted.Args) {
				t.Errorf("expected %v, got %v", tt.wanted.Args, command.Args)
			}
		})
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
			name:         "list too many arguments",
			args:         []string{"list", "extra"},
			errorMessage: "too many arguments",
		},
		{
			name:         "add missing title",
			args:         []string{"add"},
			errorMessage: "add missing title",
		},
		{
			name:         "add empty title",
			args:         []string{"add", ""},
			errorMessage: "add empty title",
		},
		{
			name:         "add empty title (space)",
			args:         []string{"add", "   "},
			errorMessage: "add empty title",
		},
		{
			name:         "add too many arguments",
			args:         []string{"add", "Learn", "Go"},
			errorMessage: "too many arguments",
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

			if command.Name != "" {
				t.Errorf("expected empty command, got %s", command.Name)
			}
		})
	}
}

func TestRunCLI_List(t *testing.T) {
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

	err = runCLI(CLICommand{Name: "list"}, service, &buf)
	if err != nil {
		t.Fatal(err)
	}

	wanted := "1. [ ] Learn Go\n"
	if wanted != buf.String() {
		t.Errorf("expected %v, got %v", wanted, buf.String())
	}
}

func TestRunCLI_Add(t *testing.T) {
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

	err = runCLI(CLICommand{Name: "add", Args: []string{"Sleep"}}, service, &buf)
	if err != nil {
		t.Fatal(err)
	}

	got := service.List()
	wanted := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
		{ID: 2, Title: "Sleep"},
	}

	if !slices.Equal(wanted, got) {
		t.Errorf("expcted %v, got %v", wanted, got)
	}
}

func TestRunCLI_Add_SaveFailed(t *testing.T) {
	todos := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
	}
	fakeRepo := &fakeRepository{
		todos:     todos,
		saveError: errors.New("save failed"),
	}
	service, err := todo.NewServiceWithRepository(fakeRepo)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer

	err = runCLI(CLICommand{Name: "add", Args: []string{"Sleep"}}, service, &buf)
	if err == nil {
		t.Fatal("expected save error")
	}

	got := service.List()
	wanted := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
	}

	if !slices.Equal(wanted, got) {
		t.Errorf("expcted %v, got %v", wanted, got)
	}
}
