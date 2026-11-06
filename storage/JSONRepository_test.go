package storage

import (
	"path/filepath"
	"slices"
	"testing"
	"todo-cli/todo"
)

func TestJSONRepository(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")

	jsonRepo := NewJSONRepository(path)

	todos := []todo.Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb"},
	}

	err := jsonRepo.Save(todos)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	got, err := jsonRepo.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if !slices.Equal(todos, got) {
		t.Errorf("expected %v, got %v", todos, got)
	}
}

func TestJSONRepository_DifferentPath(t *testing.T) {
	dir := t.TempDir()
	path1 := filepath.Join(dir, "todos1.json")
	path2 := filepath.Join(dir, "todos2.json")

	jsonRepo1 := NewJSONRepository(path1)

	todos1 := []todo.Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb"},
	}

	err := jsonRepo1.Save(todos1)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	got, err := jsonRepo1.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if !slices.Equal(todos1, got) {
		t.Errorf("expected %v, got %v", todos1, got)
	}

	jsonRepo2 := NewJSONRepository(path2)

	todos2 := []todo.Todo{
		{ID: 3, Title: "cccc"},
		{ID: 4, Title: "dddd"},
	}

	err = jsonRepo2.Save(todos2)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	got, err = jsonRepo2.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if !slices.Equal(todos2, got) {
		t.Errorf("expected %v, got %v", todos2, got)
	}

	got1, err := jsonRepo1.Load()
	if err != nil {
		t.Fatalf("Load repo1 failed: %v", err)
	}

	if !slices.Equal(todos1, got1) {
		t.Errorf("repo1 expected %v, got %v", todos1, got1)
	}
}
