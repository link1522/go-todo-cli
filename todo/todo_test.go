package todo

import (
	"slices"
	"testing"
)

func TestAddTodo(t *testing.T) {
	tests := []struct {
		name    string
		todos   []Todo
		title   string
		id      int
		wantLen int
	}{
		{
			name:    "add todo to empty list",
			todos:   []Todo{},
			title:   "item 1",
			id:      1,
			wantLen: 1,
		},
		{
			name: "add todo to existing list",
			todos: []Todo{
				{
					ID:    1,
					Title: "item 1",
					Done:  false,
				},
			},
			title:   "item 2",
			id:      2,
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddTodo(tt.todos, tt.title, tt.id)

			if len(got) != tt.wantLen {
				t.Fatalf("expected %d todos, got %d", tt.wantLen, len(got))
			}

			last := got[len(got)-1]

			if last.ID != tt.id {
				t.Errorf("expected ID %d, got %d", tt.id, last.ID)
			}

			if last.Title != tt.title {
				t.Errorf("expected Title %s, got %s", tt.title, last.Title)
			}

			if last.Done {
				t.Errorf("expected Done false, got true")
			}

			for i, original := range tt.todos {
				actual := got[i]

				if actual != original {
					t.Errorf("todo at index %d changed: expected %+v, got %+v", i, original, actual)
				}
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	tests := []struct {
		name      string
		todos     []Todo
		deleteID  int
		want      []Todo
		wantFound bool
	}{
		{
			name: "delete first",
			todos: []Todo{
				{ID: 1, Title: "A", Done: false},
				{ID: 2, Title: "B", Done: false},
				{ID: 3, Title: "C", Done: false},
			},
			deleteID: 1,
			want: []Todo{
				{ID: 2, Title: "B", Done: false},
				{ID: 3, Title: "C", Done: false},
			},
			wantFound: true,
		},
		{
			name: "delete middle",
			todos: []Todo{
				{ID: 1, Title: "A", Done: false},
				{ID: 2, Title: "B", Done: false},
				{ID: 3, Title: "C", Done: false},
			},
			deleteID: 2,
			want: []Todo{
				{ID: 1, Title: "A", Done: false},
				{ID: 3, Title: "C", Done: false},
			},
			wantFound: true,
		},
		{
			name: "delete last",
			todos: []Todo{
				{ID: 1, Title: "A", Done: false},
				{ID: 2, Title: "B", Done: false},
				{ID: 3, Title: "C", Done: false},
			},
			deleteID: 3,
			want: []Todo{
				{ID: 1, Title: "A", Done: false},
				{ID: 2, Title: "B", Done: false},
			},
			wantFound: true,
		},
		{
			name: "delete not found",
			todos: []Todo{
				{ID: 1, Title: "A", Done: false},
				{ID: 2, Title: "B", Done: false},
				{ID: 3, Title: "C", Done: false},
			},
			deleteID: 99,
			want: []Todo{
				{ID: 1, Title: "A", Done: false},
				{ID: 2, Title: "B", Done: false},
				{ID: 3, Title: "C", Done: false},
			},
			wantFound: false,
		},
		{
			name:      "delete empty slice",
			todos:     []Todo{},
			deleteID:  1,
			want:      []Todo{},
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := DeleteTodo(tt.todos, tt.deleteID)

			if found != tt.wantFound {
				t.Errorf("expected %t, got %t", tt.wantFound, found)
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestDeleteTodo_MutatesInput(t *testing.T) {
	todos := []Todo{
		{ID: 1, Title: "A"},
		{ID: 2, Title: "B"},
		{ID: 3, Title: "C"},
	}

	before := slices.Clone(todos)

	_, found := DeleteTodo(todos, 2)

	if !found {
		t.Error("expected todo to be found")
	}

	if slices.Equal(todos, before) {
		t.Error("expected input slice elements to change")
	}
}

func TestDeleteTodo_ClonePreservesOriginal(t *testing.T) {
	original := []Todo{
		{ID: 1, Title: "A"},
		{ID: 2, Title: "B"},
		{ID: 3, Title: "C"},
	}

	before := slices.Clone(original)

	// TODO 1:
	// 複製 original
	// 將複製結果傳給 DeleteTodo
	// 刪除 ID = 2
	// 將回傳結果存入 updated、found
	updated, found := DeleteTodo(slices.Clone(original), 2)

	// TODO 2:
	// 驗證 found == true
	if !found {
		t.Error("found should be true")
	}

	// TODO 3:
	// 驗證 updated 是 [A, C]
	expectedUpdated := []Todo{
		{ID: 1, Title: "A"},
		{ID: 3, Title: "C"},
	}
	if !slices.Equal(updated, expectedUpdated) {
		t.Errorf("expected %+v, got %+v", expectedUpdated, updated)
	}

	// TODO 4:
	// 驗證 original 和 before 仍相同
	if !slices.Equal(original, before) {
		t.Error("original should be the same with before")
	}
}
