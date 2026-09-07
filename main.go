package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {
	todos := []Todo{}
	nextID := 1

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.SplitN(input, " ", 2)
		switch parts[0] {
		case "add":
			if len(parts) > 1 {
				todos = addTodo(todos, parts[1], nextID)
				nextID++
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
				if completeTodo(todos, id) {
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
				newTodos, success := deleteTodo(todos, id)
				todos = newTodos

				if success {
					fmt.Println("刪除成功")
				} else {
					fmt.Println("找不到此項目")
				}
			}
		case "exit":
			err := saveTodos(todos)
			if err != nil {
				fmt.Println("保存失敗:", err)
			}
			fmt.Println("Bye!")
			return
		}
	}
}

func printTodos(todos []Todo) {
	for _, todo := range todos {
		status := "[ ]"
		if todo.Done {
			status = "[x]"
		}

		fmt.Printf("%d. %s %s\n", todo.ID, status, todo.Title)
	}
}

func addTodo(todos []Todo, title string, id int) []Todo {
	return append(todos, Todo{
		ID:    id,
		Title: title,
		Done:  false,
	})
}

func completeTodo(todos []Todo, id int) bool {
	for i, todo := range todos {
		if id == todo.ID {
			todos[i].Done = true
			return true
		}
	}
	return false
}

func deleteTodo(todos []Todo, id int) ([]Todo, bool) {
	for i, todo := range todos {
		if id == todo.ID {
			return append(todos[:i], todos[i+1:]...), true
		}
	}
	return todos, false
}

func saveTodos(todos []Todo) error {
	data, err := json.Marshal(todos)
	if err != nil {
		return err
	}

	err = os.WriteFile("todos.json", data, 0644)

	if err != nil {
		return err
	}

	return nil
}
