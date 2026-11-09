package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"todo-cli/storage"
	"todo-cli/todo"
)

func main() {
	if len(os.Args) < 2 {
		runInteractive()
		return
	}

	args := os.Args[1:]
	repo := storage.NewJSONRepository("todos.json")
	err := runCLI(args, repo, os.Stdout)
	if err != nil {
		fmt.Println(err)
	}
}

func runCLI(args []string, repo todo.TodoRepository, w io.Writer) error {
	if len(args) == 0 {
		return errors.New("missing command")
	}

	validCommands := []string{"list"}
	if !slices.Contains(validCommands, args[0]) {
		return fmt.Errorf("unknown command")
	}

	switch args[0] {
	case "list":
		if len(args) > 1 {
			return errors.New("argument error")
		}

		service, err := todo.NewServiceWithRepository(repo)
		if err != nil {
			return err
		}

		printTodos(w, service.List())
	default:
		return errors.New("unknown command")
	}

	return nil
}

func runInteractive() {
	todos, err := storage.LoadTodos()
	if err != nil {
		fmt.Println(err)
		return
	}
	nextID := todo.GetNextID(todos)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.SplitN(input, " ", 2)
		switch parts[0] {
		case "add":
			if len(parts) > 1 {
				todos = todo.AddTodo(todos, parts[1], nextID)
				nextID++

				if !save(todos) {
					continue
				}

				fmt.Println("新增成功")
			}
		case "list":
			printTodos(os.Stdout, todos)
		case "done":
			if len(parts) > 1 {
				id, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Println("ID 必須是數字")
					continue
				}
				if todo.CompleteTodo(todos, id) {
					if !save(todos) {
						continue
					}

					fmt.Println("完成成功")
				} else {
					fmt.Println("找不到此項目")
				}
			}
		case "delete":
			if len(parts) > 1 {
				id, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Println("ID 必須是數字")
					continue
				}

				newTodos, err := todo.DeleteTodo(todos, id)
				if err != nil {
					var notFoundErr *todo.TodoNotFoundError

					if errors.As(err, &notFoundErr) {
						fmt.Printf("找不到 Todo，ID: %d\n", notFoundErr.ID)
					} else {
						fmt.Println("刪除失敗: ", err)
					}

					continue
				}

				todos = newTodos

				if !save(todos) {
					continue
				}

				fmt.Println("刪除成功")
			}
		case "search":
			if len(parts) > 1 {
				foundTodos := todo.SearchTodos(todos, parts[1])
				if len(foundTodos) == 0 {
					fmt.Println("找不到任何項目")
					continue
				}
				printTodos(os.Stdout, foundTodos)
			}
		case "exit":
			fmt.Println("Bye!")
			return
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "讀取輸入失敗: ", err)
		}
	}
}

func printTodos(w io.Writer, todos []todo.Todo) {
	for _, todo := range todos {
		status := "[ ]"
		if todo.Done {
			status = "[x]"
		}

		fmt.Fprintf(w, "%d. %s %s\n", todo.ID, status, todo.Title)
	}
}

func save(todos []todo.Todo) bool {
	err := storage.SaveTodos(todos)
	if err != nil {
		fmt.Println("保存失敗:", err)
		return false
	}
	return true
}
