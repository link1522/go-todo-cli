package todo

type TodoRepository interface {
	Load() ([]Todo, error)
	Save(todos []Todo) error
}
