package main

import (
	"bytes"
	"errors"
	"slices"
	"strings"
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
		{
			name:   "done command",
			args:   []string{"done", "3"},
			wanted: CLICommand{Name: "done", ID: 3},
		},
		{
			name:   "delete command",
			args:   []string{"delete", "3"},
			wanted: CLICommand{Name: "delete", ID: 3},
		},
		{
			name:   "search command",
			args:   []string{"search", "Go"},
			wanted: CLICommand{Name: "search", Args: []string{"Go"}},
		},
		{
			name:   "search command 2",
			args:   []string{"search", "Learn Go"},
			wanted: CLICommand{Name: "search", Args: []string{"Learn Go"}},
		},
		{
			name: "search preserve whitespace",
			args: []string{"search", "  Learn Go  "},
			wanted: CLICommand{
				Name: "search",
				Args: []string{"  Learn Go  "},
			},
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
			if command.ID != tt.wanted.ID {
				t.Errorf("expected %d, got %d", tt.wanted.ID, command.ID)
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

		{
			name:         "done missing id",
			args:         []string{"done"},
			errorMessage: "done missing ID",
		},
		{
			name:         "done id is not a int",
			args:         []string{"done", "abc"},
			errorMessage: "ID must be an integer",
		},
		{
			name:         "done id should greater than zero (0)",
			args:         []string{"done", "0"},
			errorMessage: "ID must be greater than zero",
		},
		{
			name:         "done id should greater than zero (-1)",
			args:         []string{"done", "-1"},
			errorMessage: "ID must be greater than zero",
		},
		{
			name:         "done too many args",
			args:         []string{"done", "3", "extra"},
			errorMessage: "too many arguments",
		},

		{
			name:         "delete missing id",
			args:         []string{"delete"},
			errorMessage: "delete missing ID",
		},
		{
			name:         "delete id is not a int",
			args:         []string{"delete", "abc"},
			errorMessage: "ID must be an integer",
		},
		{
			name:         "delete 0",
			args:         []string{"delete", "0"},
			errorMessage: "ID must be greater than zero",
		},
		{
			name:         "delete -1",
			args:         []string{"delete", "-1"},
			errorMessage: "ID must be greater than zero",
		},
		{
			name:         "delete too many args",
			args:         []string{"delete", "3", "extra"},
			errorMessage: "too many arguments",
		},

		{
			name:         "search missing keyword",
			args:         []string{"search"},
			errorMessage: "search missing keyword",
		},
		{
			name:         "search empty keyword",
			args:         []string{"search", ""},
			errorMessage: "search empty keyword",
		},
		{
			name:         "search empty keyword (space)",
			args:         []string{"search", "   "},
			errorMessage: "search empty keyword",
		},
		{
			name:         "search too many arguments",
			args:         []string{"search", "Go", "extra"},
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
				t.Errorf("expected command name '', got %s", command.Name)
			}
			if command.ID != 0 {
				t.Errorf("expected command ID 0, got %d", command.ID)
			}
			if len(command.Args) > 0 {
				t.Errorf("expected length of command args to be 0, got %d", len(command.Args))
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

	if !strings.Contains(buf.String(), "2. Sleep 添加成功") {
		t.Errorf("unexpected stdout: %q", buf.String())
	}

	got := service.List()
	wanted := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
		{ID: 2, Title: "Sleep"},
	}

	if !slices.Equal(wanted, got) {
		t.Errorf("expcted %v, got %v", wanted, got)
	}

	if !slices.Equal(wanted, fakeRepo.todos) {
		t.Errorf("expcted %v, got %v", wanted, fakeRepo.todos)
	}
}

func TestRunCLI_Add_SaveFailed(t *testing.T) {
	todos := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
	}
	saveError := errors.New("save failed")
	fakeRepo := &fakeRepository{
		todos:     todos,
		saveError: saveError,
	}
	service, err := todo.NewServiceWithRepository(fakeRepo)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer

	err = runCLI(CLICommand{Name: "add", Args: []string{"Sleep"}}, service, &buf)
	if !errors.Is(err, saveError) {
		t.Fatal("expected save error")
	}

	if buf.Len() != 0 {
		t.Errorf("unexpected stdout: %q", buf.String())
	}

	got := service.List()
	wanted := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
	}

	if !slices.Equal(wanted, got) {
		t.Errorf("expcted %v, got %v", wanted, got)
	}

	if !slices.Equal(wanted, fakeRepo.todos) {
		t.Errorf("expcted %v, got %v", wanted, fakeRepo.todos)
	}
}

func TestRunAppCommand_Error(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		errorMessage string
	}{
		{
			name:         "unknown command",
			args:         []string{"unknown"},
			errorMessage: "unknown command",
		},
		{
			name:         "add missing title",
			args:         []string{"add"},
			errorMessage: "add missing title",
		},
		{
			name:         "too many arguments",
			args:         []string{"list", "extra"},
			errorMessage: "too many arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			exitCode := runApp(tt.args, &stdout, &stderr)

			if exitCode != 1 {
				t.Errorf("expected 1, got %d", exitCode)
			}

			if stdout.String() != "" {
				t.Errorf("expected \"\", got %s", stdout.String())
			}

			if !strings.Contains(stderr.String(), tt.errorMessage) {
				t.Errorf("unexpected stderr: %q", stderr.String())
			}
		})
	}
}

func TestRunCLI_Done(t *testing.T) {
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

	err = runCLI(CLICommand{Name: "done", ID: 1}, service, &buf)
	if err != nil {
		t.Fatal(err)
	}

	wantedMsg := "1 已完成"
	if !strings.Contains(buf.String(), wantedMsg) {
		t.Errorf("unexpected stdout: %q", buf.String())
	}

	got := service.List()
	wanted := []todo.Todo{
		{ID: 1, Title: "Learn Go", Done: true},
	}
	if !slices.Equal(wanted, got) {
		t.Errorf("expected %v, got %v", wanted, got)
	}

	if !slices.Equal(wanted, fakeRepo.todos) {
		t.Errorf("expected %v, got %v", wanted, fakeRepo.todos)
	}
}

func TestRunCLI_Done_NotFound(t *testing.T) {
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

	targetID := 99
	err = runCLI(CLICommand{Name: "done", ID: targetID}, service, &buf)
	if !errors.Is(err, todo.ErrTodoNotFound) {
		t.Errorf("expect error %v, got %v", todo.ErrTodoNotFound, err)
	}
	var notFoundErr *todo.TodoNotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Fatalf("expect error %v, got %v", notFoundErr, err)
	}
	if notFoundErr.ID != targetID {
		t.Errorf("expect ID %d, got %d", targetID, notFoundErr.ID)
	}

	if buf.Len() != 0 {
		t.Errorf("unexpected stdout: %q", buf.String())
	}

	got := service.List()
	wanted := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
	}
	if !slices.Equal(wanted, got) {
		t.Errorf("expected %v, got %v", wanted, got)
	}
}

func TestRunCLI_Done_SaveFailed(t *testing.T) {
	todos := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
	}
	saveError := errors.New("save failed")
	fakeRepo := &fakeRepository{
		todos:     todos,
		saveError: saveError,
	}
	service, err := todo.NewServiceWithRepository(fakeRepo)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer

	err = runCLI(CLICommand{Name: "done", ID: 1}, service, &buf)
	if !errors.Is(err, saveError) {
		t.Fatal("expected save error")
	}
	if len(buf.String()) != 0 {
		t.Errorf("unexpected stdout: %q", buf.String())
	}

	got := service.List()
	wanted := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
	}
	if !slices.Equal(wanted, got) {
		t.Errorf("expected %v, got %v", wanted, got)
	}
	if !slices.Equal(wanted, fakeRepo.todos) {
		t.Errorf("expected %v, got %v", wanted, fakeRepo.todos)
	}
}

func TestParsePositiveID(t *testing.T) {
	id, err := parsePositiveID("3")
	if err != nil {
		t.Fatal(err)
	}
	if id != 3 {
		t.Errorf("expected %d, got %d", 3, id)
	}
}

func TestParsePositiveID_Failed(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wanted   int
		errorMsg string
	}{
		{
			name:     "not integer",
			input:    "abc",
			wanted:   0,
			errorMsg: "ID must be an integer",
		},
		{
			name:     "zero",
			input:    "0",
			wanted:   0,
			errorMsg: "ID must be greater than zero",
		},
		{
			name:     "-1",
			input:    "-1",
			wanted:   0,
			errorMsg: "ID must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := parsePositiveID(tt.input)
			if err == nil {
				t.Fatalf("expected error: %s", tt.errorMsg)
			}
			if id != 0 {
				t.Errorf("expected id 0 when failed, but got %d", id)
			}
			if !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("unexpected error message: %v", err)
			}
		})
	}
}

func TestRunCLI_Delete(t *testing.T) {
	todos := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
		{ID: 2, Title: "Learn Python"},
		{ID: 3, Title: "Learn Rust"},
	}
	fakeRepo := &fakeRepository{
		todos: todos,
	}
	service, err := todo.NewServiceWithRepository(fakeRepo)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer

	err = runCLI(CLICommand{Name: "delete", ID: 2}, service, &buf)
	if err != nil {
		t.Fatal(err)
	}

	wantedMsg := "2 刪除成功"
	if !strings.Contains(buf.String(), wantedMsg) {
		t.Errorf("unexpected stdout: %q", buf.String())
	}

	got := service.List()
	wanted := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
		{ID: 3, Title: "Learn Rust"},
	}
	if !slices.Equal(wanted, got) {
		t.Errorf("expected %v, got %v", wanted, got)
	}

	if !slices.Equal(wanted, fakeRepo.todos) {
		t.Errorf("expected %v, got %v", wanted, fakeRepo.todos)
	}
}

func TestRunCLI_Delete_NotFound(t *testing.T) {
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

	targetID := 99
	err = runCLI(CLICommand{Name: "delete", ID: targetID}, service, &buf)
	if !errors.Is(err, todo.ErrTodoNotFound) {
		t.Errorf("expect error %v, got %v", todo.ErrTodoNotFound, err)
	}
	var notFoundErr *todo.TodoNotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Fatalf("expect error %v, got %v", notFoundErr, err)
	}
	if notFoundErr.ID != targetID {
		t.Errorf("expect ID %d, got %d", targetID, notFoundErr.ID)
	}

	if buf.Len() != 0 {
		t.Errorf("unexpected stdout: %q", buf.String())
	}

	got := service.List()
	wanted := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
	}
	if !slices.Equal(wanted, got) {
		t.Errorf("expected %v, got %v", wanted, got)
	}
}

func TestRunCLI_Delete_SaveFailed(t *testing.T) {
	todos := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
	}
	saveError := errors.New("save failed")
	fakeRepo := &fakeRepository{
		todos:     todos,
		saveError: saveError,
	}
	service, err := todo.NewServiceWithRepository(fakeRepo)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer

	err = runCLI(CLICommand{Name: "delete", ID: 1}, service, &buf)
	if !errors.Is(err, saveError) {
		t.Fatal("expected save error")
	}
	if len(buf.String()) != 0 {
		t.Errorf("unexpected stdout: %q", buf.String())
	}

	got := service.List()
	wanted := []todo.Todo{
		{ID: 1, Title: "Learn Go"},
	}
	if !slices.Equal(wanted, got) {
		t.Errorf("expected %v, got %v", wanted, got)
	}
	if !slices.Equal(wanted, fakeRepo.todos) {
		t.Errorf("expected %v, got %v", wanted, fakeRepo.todos)
	}
}

func TestRunCLI_Search(t *testing.T) {
	tests := []struct {
		name      string
		todos     []todo.Todo
		keyword   string
		wantedMsg string
	}{
		{
			name: "search matched",
			todos: []todo.Todo{
				{
					ID:    1,
					Title: "Learn Go",
					Done:  true,
				},
				{
					ID:    2,
					Title: "Sleep",
				},
				{
					ID:    3,
					Title: "Learn TS",
				},
			},
			keyword:   "Learn",
			wantedMsg: "1. [x] Learn Go\n3. [ ] Learn TS\n",
		},
		{
			name: "search not found",
			todos: []todo.Todo{
				{
					ID:    1,
					Title: "Learn Go",
					Done:  true,
				},
				{
					ID:    2,
					Title: "Sleep",
				},
				{
					ID:    3,
					Title: "Learn TS",
				},
			},
			keyword:   "apple",
			wantedMsg: "找不到任何項目\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeRepo := &fakeRepository{
				todos: slices.Clone(tt.todos),
			}
			service, err := todo.NewServiceWithRepository(fakeRepo)
			if err != nil {
				t.Fatal(err)
			}
			var buf bytes.Buffer

			err = runCLI(CLICommand{Name: "search", Args: []string{tt.keyword}}, service, &buf)
			if err != nil {
				t.Fatal(err)
			}

			if got := buf.String(); got != tt.wantedMsg {
				t.Errorf("expected %q, got %q", tt.wantedMsg, got)
			}

			if !slices.Equal(service.List(), tt.todos) {
				t.Errorf("expected %v, got %v", tt.todos, service.List())
			}
			if !slices.Equal(fakeRepo.todos, tt.todos) {
				t.Errorf("expected %v, got %v", tt.todos, fakeRepo.todos)
			}
		})
	}
}
