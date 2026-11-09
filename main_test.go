package main

import (
	"bytes"
	"testing"
	"todo-cli/todo"
)

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

func TestRunCLI(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		hasError     bool
		errorMessage string
	}{
		{
			name:         "missing command",
			args:         []string{},
			hasError:     true,
			errorMessage: "missing command",
		},
		{
			name:         "unknown command",
			args:         []string{"unknown"},
			hasError:     true,
			errorMessage: "unknown command",
		},
		{
			name:         "argument error",
			args:         []string{"list", "extra"},
			hasError:     true,
			errorMessage: "argument error",
		},
		{
			name:         "normal",
			args:         []string{"list"},
			hasError:     false,
			errorMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runCLI(tt.args)
			gotError := err != nil

			if gotError != tt.hasError {
				if tt.hasError {
					t.Fatalf("expected error: %s", tt.errorMessage)
				} else {
					t.Errorf("it should not have error, but got %v", err)
				}
			}

			if gotError {
				if tt.errorMessage != err.Error() {
					t.Errorf("expected %s, got %s", tt.errorMessage, err.Error())
				}
			}
		})
	}
}
