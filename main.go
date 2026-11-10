package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
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
	command, err := parseCommand(args)
	if err != nil {
		fmt.Println(err)
		return
	}

	repo := storage.NewJSONRepository("todos.json")
	service, err := todo.NewServiceWithRepository(repo)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = runCLI(command, service, os.Stdout)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func runCLI(command CLICommand, service *todo.Service, w io.Writer) error {
	switch command.Name {
	case "list":
		printTodos(w, service.List())
	case "add":
		todo, err := service.AddAndSave(command.Args[0])
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d. %s 添加成功\n", todo.ID, todo.Title)
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

func parseCommand(args []string) (CLICommand, error) {
	if len(args) == 0 {
		return CLICommand{}, errors.New("missing command")
	}

	switch args[0] {
	case "list":
		if len(args) > 1 {
			return CLICommand{}, errors.New("too many arguments")
		}
		return CLICommand{Name: "list"}, nil
	case "add":
		if len(args) > 2 {
			return CLICommand{}, errors.New("too many arguments")
		}
		if len(args) < 2 {
			return CLICommand{}, errors.New("add missing title")
		}
		if strings.TrimSpace(args[1]) == "" {
			return CLICommand{}, errors.New("add empty title")
		}
		return CLICommand{Name: "add", Args: []string{args[1]}}, nil
	default:
		return CLICommand{}, errors.New("unknown command")
	}
}

type CLICommand struct {
	Name string
	Args []string
}
