package todo

import (
	"errors"
	"testing"
)

func TestTodoService_NextID(t *testing.T) {
	initial := []Todo{
		{ID: 7, Title: "aaaa"},
	}
	service := NewService(initial)

	wantedNextId := 8
	wantedNextTitle := "bbbb"
	next := service.Add(wantedNextTitle)

	if next.ID != wantedNextId {
		t.Errorf("expected id %d, got %d", wantedNextId, next.ID)
	}

	list := service.List()
	if len(list) != 2 {
		t.Fatalf("expected length of list is 2, but got %d", len(list))
	}

	last := len(list) - 1
	if list[last].ID != wantedNextId {
		t.Errorf("expected id %d, got %d", wantedNextId, list[last].ID)
	}
	if list[last].Title != wantedNextTitle {
		t.Errorf("expected title %s, got %s", wantedNextTitle, list[last].Title)
	}
}

func TestTodoService_NotModifyInitialTodos(t *testing.T) {
	initial := []Todo{
		{ID: 7, Title: "aaaa"},
	}
	service := NewService(initial)

	initial[0].ID = 999

	list := service.List()
	if list[0].ID != 7 {
		t.Error("initial slice cannot be modify")
	}
}

func TestTodoService_ListNotModifyTodos(t *testing.T) {
	initial := []Todo{
		{ID: 7, Title: "aaaa"},
	}
	service := NewService(initial)

	list := service.List()
	list[0].ID = 999

	list2 := service.List()

	if list2[0].ID != 7 {
		t.Error("list should not modify todos in service")
	}
}

func TestTodoService_Complete(t *testing.T) {
	service := NewService([]Todo{{ID: 1, Title: "aaaa"}})
	err := service.Complete(1)

	if err != nil {
		t.Fatalf("got error %v", err)
	}

	list := service.List()

	if !list[0].Done {
		t.Error("todo should be done")
	}
}

func TestTodoService_CompleteNotFound(t *testing.T) {
	service := NewService([]Todo{{ID: 1, Title: "aaaa"}})
	completeId := 99
	err := service.Complete(completeId)

	if !errors.Is(err, ErrTodoNotFound) {
		t.Errorf("expected %v, got %v", ErrTodoNotFound, err)
	}

	var notFoundErr *TodoNotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Fatalf("expected %v, got %v", notFoundErr, err)
	}

	if notFoundErr.ID != completeId {
		t.Errorf("expected %d, got %d", completeId, notFoundErr.ID)
	}
}

func TestTodoService_CompleteAlreadyDone(t *testing.T) {
	service := NewService([]Todo{{ID: 1, Title: "aaaa", Done: true}})
	err := service.Complete(1)

	if err != nil {
		t.Fatalf("got error %v", err)
	}

	list := service.List()

	if !list[0].Done {
		t.Error("todo should be done")
	}
}
