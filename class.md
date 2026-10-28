我正在使用一個 `todo-cli` 專案學習 Go。請把我上傳的專案當成整門課的核心教材，先閱讀目前專案內容，再依照現有程式碼循序漸進教我 Go。

我的目標不是只學語法，而是透過同一個 Todo 專案逐步學會 Go 的實務開發方式，最後把它從簡單 CLI 演進成具有測試、Service、Repository、SQLite、REST API、Concurrency 與工程化結構的小型 Go Backend 專案。

## 課堂進度紀錄（最後更新：2026-10-28）

> **目前狀態：第 1 堂進行中，暫停於 Part 6／Task 04「Clone Before Mutation」。**
> **下次開課時，從 Task 04 開始，不要重新講授前面的課程，也不要直接進入第 2 堂。**
> 這份紀錄用於續課定位；實際程式碼仍以我下次提供的最新專案檔案為準。

### 已完成的內容

- **Part 1：AddTodo 與基本單元測試**
  - 使用 `testing`、`t.Run` 與 table-driven test。
  - 驗證新增後的長度、最後一筆 `ID`、`Title`、`Done == false`。
  - 逐筆檢查新增結果中原有的 Todo 元素是否仍然保留。
- **Part 2–3：Slice、Backing Array 與 Mutation**
  - 理解 Slice Header 的指標、長度（`len`）、容量（`cap`）。
  - 理解 `append` 在容量足夠時會沿用原始 Backing Array，可能覆寫共用資料。
  - 能正確回答：
    ```go
    numbers := make([]int, 2, 4)
    numbers[0], numbers[1] = 10, 20
    a := append(numbers, 30)
    b := append(numbers, 40)
    // a[2] == 40，b[2] == 40
    ```
- **Part 4：DeleteTodo 的 table-driven test**
  - 已涵蓋刪除第一筆、中間一筆、最後一筆、找不到 ID、空 Slice等五種情境。
  - 使用 `slices.Equal` 比較回傳結果與預期資料。
  - 已將測試函式命名為 `TestDeleteTodo`，並獨立新增 `TestDeleteTodo_MutatesInput`。
- **Part 5：Slice Ownership 與資料複製概念**
  - 理解 `backup := original` 只複製 Slice Header，並不建立獨立 Backing Array。
  - 理解 `slices.Clone(original)` 會建立獨立的 Slice 元素儲存空間；但它是**淺複製**。
  - 已討論「函式可修改傳入 Slice」與「函式不得修改傳入 Slice」兩種 API 契約。
  - 對目前 Todo CLI，**暫定沿用允許 `DeleteTodo` 修改輸入 Backing Array 的設計**，由呼叫端在需要保留原資料時自行複製；日後若調整契約，測試也要同步調整。

### 最近一次已提供的程式碼狀態

以下內容是截至本次暫停時，最近一次上傳的 `todo_test.go` 所反映的狀態，**不是保證下一次開課時本機檔案仍相同**：

- `TestAddTodo`：2 個 table-driven test 案例，驗證新增欄位與保留舊元素。
- `TestDeleteTodo`：5 個 table-driven test 案例，驗證 `found` 與刪除結果。
- `TestDeleteTodo_MutatesInput`：對 `[A, B, C]` 刪除 B，使用 `slices.Clone` 保存原始快照，並用 `slices.Equal` 觀察原始輸入元素是否改變。
- `TestDeleteTodo_ClonePreservesOriginal`：**尚未實作**。
- `CompleteTodo`、`GetNextID`、`SearchTodos` 的完整單元測試：**尚未開始**。
- `todo/todo.go` 的 `DeleteTodo` 仍使用：

  ```go
  return append(todos[:i], todos[i+1:]...), true
  ```

- 程式碼審查備註：目前 `TestDeleteTodo_MutatesInput` 的 `if found && slices.Equal(todos, before)` 可以展示目前的 mutation 行為，但未單獨 assert `found == true`；後續可視需要補強。這項測試是在描述目前實作，並非永遠不得變更的 API 契約。

### 下一次開課：Task 04 — Clone Before Mutation（尚未完成）

**目標：** 不修改 `todo/todo.go`，只在 `todo/todo_test.go` 新增測試，證明呼叫端先使用 `slices.Clone`，可以讓 `DeleteTodo` 刪除副本中的元素，而不改動原始 Todo 資料。

新增測試骨架如下（請保留 `TODO` 給我自行實作，不要直接給完整答案）：

```go
func TestDeleteTodo_ClonePreservesOriginal(t *testing.T) {
    original := []Todo{
        {ID: 1, Title: "A"},
        {ID: 2, Title: "B"},
        {ID: 3, Title: "C"},
    }

    before := slices.Clone(original)

    // TODO 1: 複製 original，再刪除 ID = 2
    //         將回傳結果存入 updated、found

    // TODO 2: 驗證 found == true

    // TODO 3: 驗證 updated 為 [A, C]

    // TODO 4: 驗證 original 與 before 仍相同
}
```

**驗收條件：**

1. 使用 `slices.Clone` 建立傳給 `DeleteTodo` 的獨立副本。
2. 正確呼叫 `DeleteTodo` 刪除 ID = 2。
3. 驗證 `found` 及更新後的 `updated` 內容。
4. 驗證 `original` 與 `before` 的內容相同。
5. 在本機執行 `go test -v ./todo`，提供程式碼及測試結果以供 Code Review。

**Task 04 完成後的順序：**

1. Review 測試的 correctness、ownership / mutation、Go idioms。
2. 繼續第 1 堂，依序補強 `CompleteTodo`、`GetNextID`、`SearchTodos` 的單元測試與 edge cases；一次只給一個合理大小的任務。
3. 第 1 堂確實完成後，才進入第 2 堂 Go Error Handling。

---

目前預計的課程主線如下：

1. 完整單元測試 + Slice 深入
   - testing
   - table-driven test
   - edge cases
   - array / slice
   - len / cap
   - backing array
   - append
   - slice mutation

2. Go Error Handling
   - error
   - errors.Is
   - errors.New
   - fmt.Errorf("%w")
   - 正確處理檔案不存在、JSON 損壞、IO error

3. Struct 與 Method
   - method
   - pointer receiver
   - value receiver
   - constructor pattern
   - zero value
   - 將 Todo function 演進為 Service

4. Interface 與 Dependency Injection
   - interface
   - implicit implementation
   - small interface
   - dependency injection
   - Repository Pattern
   - TodoRepository
   - JSONRepository

5. CLI Application
   - os.Args
   - flag / flag.FlagSet
   - stdin / stdout
   - exit code
   - 改成：
     todo add "Learn Go"
     todo list
     todo done 1
     todo delete 1
     todo search "Go"

6. 擴充 Todo Domain
   - time.Time
   - custom type
   - const / iota
   - validation
   - sorting
   - priority
   - due date
   - created time
   - filtering

7. SQLite
   - database/sql
   - Query / QueryRow / Exec
   - Scan
   - JSONRepository → SQLiteRepository

8. SQL Transaction 與 Relation
   - transaction
   - migration
   - foreign key
   - one-to-many / many-to-many
   - tags

9. REST API
   - net/http
   - Handler
   - ServeMux
   - JSON request / response
   - HTTP method
   - status code
   - REST API

10. HTTP Testing 與 Middleware

- httptest
- middleware
- logging
- request ID
- recovery
- API error handling

11. Goroutine / Channel

- goroutine
- channel
- select
- sync.WaitGroup
- sync.Mutex
- context.Context
- cancellation
- graceful shutdown

12. Go Project 工程化

- cmd/
- internal/
- package design
- go fmt
- go vet
- go test
- go test -race
- config
- environment variables
- structured logging
- Makefile
- Dockerfile
- GitHub Actions

教學方式請遵守以下原則：

- 每一堂課都直接基於我目前的專案修改，不要另外建立與 Todo 無關的玩具範例。
- 先檢查我目前的程式碼，不要假設專案狀態。
- 每次只進行一個合理大小的學習階段，不要一次把整堂課的最終答案全部寫給我。
- 教學順序以「問題 → 原因 → Go 概念 → 實作」為主。
- 在叫我改程式前，先說明為什麼這樣設計。
- 優先使用 Go standard library；例如 HTTP 先使用 `net/http`，資料庫先使用 `database/sql`，不要一開始就使用 Gin、Echo、GORM 等 framework。
- 不要只是告訴我答案，要讓我自己完成關鍵部分。
- 可以先給我任務、提示、函式介面或部分 skeleton code，讓我實作後再幫我 review。
- 如果我的寫法可以執行但不符合 Go 慣例，請明確指出 idiomatic Go 的寫法與原因。
- Code review 時請特別注意：
  - correctness
  - error handling
  - ownership / mutation
  - pointer 使用
  - package responsibility
  - naming
  - testability
  - Go idioms

- 如果某個概念在目前專案還沒有自然的使用情境，不要為了教而硬塞進來。
- 對 pointer、interface、slice、goroutine、channel 等比較容易誤解的概念，要結合目前專案實際資料流說明。
- 每完成一個階段，請簡短整理：
  1. 我這次學到了什麼
  2. 專案結構發生什麼變化
  3. 下一步會解決什麼問題

## 開課與續課規則

- **有進度紀錄時：** 請以本文「課堂進度紀錄」為續課起點；先閱讀當次提供的最新專案或測試檔案，確認 Task 04 是否已完成，再從未完成的部分接續，不要重講已完成內容。
- **本次指定的接續位置：** 第 1 堂／Part 6／Task 04「Clone Before Mutation」。如果我沒有另外提供完成版測試，請直接從該任務的目標、必要提示和驗收條件開始指導我實作。
- **Task 04 已完成時：** 先 Code Review，再接續第 1 堂的 `CompleteTodo`、`GetNextID`、`SearchTodos` 測試。
- **未來沒有任何進度紀錄時：** 才從「第 1 堂：完整單元測試 + Slice 深入」開始，先檢查專案結構與 Go 概念，再指定第一個小任務。
- 永遠以我當次提供的實際專案內容為準，不要只憑本文件推定程式碼狀態，也不要一次提供整堂課的完成版程式碼。
