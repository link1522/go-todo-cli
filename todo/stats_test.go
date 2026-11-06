package todo

import (
	"slices"
	"testing"
)

func TestCountPending(t *testing.T) {
	service := NewService([]Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb", Done: true},
		{ID: 3, Title: "cccc"},
	})

	wanted := 2
	got := CountPending(service)

	if wanted != got {
		t.Errorf("expect %d, got %d", wanted, got)
	}
}

type fakeTodoReader struct{ todos []Todo }

func (f *fakeTodoReader) List() []Todo {
	return slices.Clone(f.todos)
}

func TestCountPending_FakeReader(t *testing.T) {
	fakeReader := &fakeTodoReader{
		todos: []Todo{
			{ID: 1, Title: "aaaa"},
			{ID: 2, Title: "bbbb", Done: true},
			{ID: 3, Title: "cccc"},
			{ID: 4, Title: "dddd"},
		},
	}

	wanted := 3
	got := CountPending(fakeReader)

	if wanted != got {
		t.Errorf("expect %d, got %d", wanted, got)
	}
}

func TestCountPending_Empty(t *testing.T) {
	service := NewService([]Todo{})

	wanted := 0
	got := CountPending(service)

	if wanted != got {
		t.Errorf("expect %d, got %d", wanted, got)
	}
}

func TestStatsService(t *testing.T) {
	service := NewService([]Todo{
		{ID: 1, Title: "aaaa"},
		{ID: 2, Title: "bbbb", Done: true},
		{ID: 3, Title: "cccc"},
	})
	statsService := NewStatsService(service)

	wanted := 2
	got := statsService.PendingCount()

	if wanted != got {
		t.Errorf("expect %d, got %d", wanted, got)
	}
}

func TestStatsService_FakeReader(t *testing.T) {
	fakeReader := &fakeTodoReader{
		todos: []Todo{
			{ID: 1, Title: "aaaa"},
			{ID: 2, Title: "bbbb", Done: true},
			{ID: 3, Title: "cccc"},
			{ID: 4, Title: "dddd"},
		},
	}
	statsService := NewStatsService(fakeReader)

	wanted := 3
	got := statsService.PendingCount()

	if wanted != got {
		t.Errorf("expect %d, got %d", wanted, got)
	}
}
