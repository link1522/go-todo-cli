package todo

import (
	"errors"
	"slices"
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

func TestTodoService_Delete(t *testing.T) {
	service := NewService([]Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb"},
		{ID: 3, Title: "cccc"},
	})
	err := service.Delete(2)

	if err != nil {
		t.Fatalf("got error %v", err)
	}

	got := service.List()
	wanted := []Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 3, Title: "cccc"},
	}

	if !slices.Equal(got, wanted) {
		t.Errorf("expected %v, got %v", wanted, got)
	}
}

func TestTodoService_DeleteNotFound(t *testing.T) {
	service := NewService([]Todo{{ID: 1, Title: "aaaa"}})
	deleteId := 99
	err := service.Delete(deleteId)

	if !errors.Is(err, ErrTodoNotFound) {
		t.Errorf("expected error %v, got %v", ErrTodoNotFound, err)
	}

	var notFoundErr *TodoNotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Fatalf("expected error %v, got %v", notFoundErr, err)
	}

	if notFoundErr.ID != deleteId {
		t.Errorf("expected %d, got %d", deleteId, notFoundErr.ID)
	}

	wanted := []Todo{{ID: 1, Title: "aaaa"}}
	got := service.List()
	if !slices.Equal(wanted, got) {
		t.Errorf("expected %v, got %v", wanted, got)
	}
}

func TestTodoService_Search(t *testing.T) {
	service := NewService([]Todo{
		{ID: 1, Title: "Learn Go"},
		{ID: 2, Title: "Learn JS"},
		{ID: 3, Title: "Learn Rust"},
	})

	got1 := service.Search("Learn Go")
	wanted1 := []Todo{{ID: 1, Title: "Learn Go"}}
	if !slices.Equal(wanted1, got1) {
		t.Errorf("expected %v, got %v", wanted1, got1)
	}

	got2 := service.Search("Learn")
	wanted2 := []Todo{
		{ID: 1, Title: "Learn Go"},
		{ID: 2, Title: "Learn JS"},
		{ID: 3, Title: "Learn Rust"},
	}
	if !slices.Equal(wanted2, got2) {
		t.Errorf("expected %v, got %v", wanted2, got2)
	}

	got3 := service.Search("LEARN")
	wanted3 := []Todo{
		{ID: 1, Title: "Learn Go"},
		{ID: 2, Title: "Learn JS"},
		{ID: 3, Title: "Learn Rust"},
	}
	if !slices.Equal(wanted3, got3) {
		t.Errorf("expected %v, got %v", wanted3, got3)
	}

	got4 := service.Search("Python")
	wanted4 := []Todo{}
	if !slices.Equal(wanted4, got4) {
		t.Errorf("expected %v, got %v", wanted4, got4)
	}

	list := service.List()
	originalList := []Todo{
		{ID: 1, Title: "Learn Go"},
		{ID: 2, Title: "Learn JS"},
		{ID: 3, Title: "Learn Rust"},
	}
	if !slices.Equal(list, originalList) {
		t.Errorf("expected %v, got %v", originalList, list)
	}
}

func TestTodoService_SearchOwnership(t *testing.T) {
	service := NewService([]Todo{{ID: 1, Title: "Learn Go"}})

	matched := service.Search("Learn")
	matched[0].Title = "Learn TS"

	got := service.List()
	wanted := []Todo{{ID: 1, Title: "Learn Go"}}

	if !slices.Equal(got, wanted) {
		t.Errorf("expected %v, got %v", wanted, got)
	}
}

type fakeRepository struct {
	todos     []Todo
	loadError error
	saveError error
}

func (r *fakeRepository) Load() ([]Todo, error) {
	if r.loadError != nil {
		return nil, r.loadError
	}
	return r.todos, nil
}

func (r *fakeRepository) Save(todos []Todo) error {
	if r.saveError != nil {
		return r.saveError
	}
	r.todos = todos
	return nil
}

func TestTodoService_WithRepo(t *testing.T) {
	initTodos := []Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb"},
	}
	fakeRepo := fakeRepository{todos: initTodos}
	service, err := NewServiceWithRepository(&fakeRepo)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	got := service.List()
	if !slices.Equal(initTodos, got) {
		t.Errorf("expected %v, got %v", initTodos, got)
	}

	service.Add("cccc")
	err = service.Save()
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	got, err = fakeRepo.Load()
	wanted := []Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb"},
		{ID: 3, Title: "cccc"},
	}

	if !slices.Equal(wanted, got) {
		t.Errorf("expected %v, got %v", wanted, got)
	}
}

func TestTodoService_WithNilRepo(t *testing.T) {
	_, err := NewServiceWithRepository(nil)

	if err == nil {
		t.Errorf("expected error when passing nil repository")
	}
}

func TestTodoService_WithRepoLoadFailed(t *testing.T) {
	initTodos := []Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb"},
	}
	fakeRepo := fakeRepository{todos: initTodos, loadError: errors.New("Load failed")}
	_, err := NewServiceWithRepository(&fakeRepo)

	if err == nil {
		t.Errorf("expected error when loading failed")
	}
}

func TestTodoService_WithRepoSaveFailed(t *testing.T) {
	initTodos := []Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb"},
	}
	fakeRepo := fakeRepository{todos: initTodos, saveError: errors.New("Save failed")}
	service, err := NewServiceWithRepository(&fakeRepo)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	err = service.Save()
	if err == nil {
		t.Errorf("expected error when saving failed")
	}
}

func TestTodoService_MissRepo(t *testing.T) {
	service := NewService([]Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb"},
	})

	err := service.Save()

	if err == nil {
		t.Error("expected error when missing repository")
	}

	_, err = service.AddAndSave("cccc")

	if err == nil {
		t.Error("expected error when missing repository")
	}
}

func TestTodoService_ModifyLoadedTodos(t *testing.T) {
	initTodos := []Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb"},
	}
	fakeRepo := fakeRepository{todos: initTodos}
	service, err := NewServiceWithRepository(&fakeRepo)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	todos := service.List()
	todos[0].Title = "test"

	got := service.List()
	if !slices.Equal(initTodos, got) {
		t.Errorf("expected %v, got %v", initTodos, got)
	}
}

func TestTodoService_AddAndSave(t *testing.T) {
	initTodos := []Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb"},
	}
	fakeRepo := fakeRepository{todos: initTodos}
	service, err := NewServiceWithRepository(&fakeRepo)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	newTodo, err := service.AddAndSave("cccc")

	wanted := append(slices.Clone(initTodos), newTodo)
	got := service.List()

	if !slices.Equal(wanted, got) {
		t.Errorf("expected %v, got %v", wanted, got)
	}

	gotByRepo, err := service.repo.Load()
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if !slices.Equal(wanted, gotByRepo) {
		t.Errorf("expected %v, got %v", wanted, gotByRepo)
	}
}

func TestTodoService_AddAndSaveFailed(t *testing.T) {
	initTodos := []Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb"},
	}
	fakeRepo := fakeRepository{
		todos:     initTodos,
		saveError: errors.New("save failed"),
	}
	service, err := NewServiceWithRepository(&fakeRepo)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	_, err = service.AddAndSave("cccc")
	if err == nil {
		t.Error("expect got error when saving failed")
	}

	got := service.List()
	if !slices.Equal(initTodos, got) {
		t.Errorf("expected %v, got %v", initTodos, got)
	}

	gotByRepo, err := service.repo.Load()
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if !slices.Equal(initTodos, gotByRepo) {
		t.Errorf("expected %v, gotByRepo %v", initTodos, gotByRepo)
	}
}

func TestTodoService_AddAndSaveOwnership(t *testing.T) {
	repo := &fakeRepository{
		todos: []Todo{
			{ID: 1, Title: "Learn Go"},
		},
	}

	service, err := NewServiceWithRepository(repo)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.AddAndSave("Learn SQL")
	if err != nil {
		t.Fatal(err)
	}

	repo.todos[0].Title = "Modified By Repository"

	got := service.List()

	if got[0].Title != "Learn Go" {
		t.Errorf("Service data was modified: %s", got[0].Title)
	}
}
