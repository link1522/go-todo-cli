package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo-cli/storage"
	"todo-cli/todo"
)

func main() {
	todos := storage.LoadTodos()
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
			printTodos(todos)
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
				newTodos, success := todo.DeleteTodo(todos, id)

				if success {
					todos = newTodos

					if !save(todos) {
						continue
					}

					fmt.Println("刪除成功")
				} else {
					fmt.Println("找不到此項目")
				}
			}
		case "search":
			if len(parts) > 1 {
				foundTodos := todo.SearchTodos(todos, parts[1])
				if len(foundTodos) == 0 {
					fmt.Println("找不到任何項目")
					continue
				}
				printTodos(foundTodos)
			}
		case "exit":
			fmt.Println("Bye!")
			return
		}
	}
}

func printTodos(todos []todo.Todo) {
	for _, todo := range todos {
		status := "[ ]"
		if todo.Done {
			status = "[x]"
		}

		fmt.Printf("%d. %s %s\n", todo.ID, status, todo.Title)
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
