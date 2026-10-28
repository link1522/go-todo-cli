package todo

import "testing"

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
