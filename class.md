我正在使用一個 `todo-cli` 專案學習 Go。請把我上傳的專案當成整門課的核心教材，先閱讀目前專案內容，再依照現有程式碼循序漸進教我 Go。

我的目標不是只學語法，而是透過同一個 Todo 專案逐步學會 Go 的實務開發方式，最後把它從簡單 CLI 演進成具有測試、Service、Repository、SQLite、REST API、Concurrency 與工程化結構的小型 Go Backend 專案。

## 課堂進度紀錄（最後更新：2026-10-29）

> **目前狀態：第 1 堂「完整單元測試 + Slice 深入」已完成。**
> **下次上課直接從第 2 堂「Go Error Handling」／Part 1／Task 01「不要吞掉檔案讀取錯誤」開始。**
> 第 2 堂的 Task 01 尚未實作；不要重講第 1 堂，也不要假設第 2 堂已完成。
> 本紀錄用於定位進度；程式碼仍以我下次上傳的最新專案版本為準。

### 第 1 堂已完成內容

- **Part 1：AddTodo 與基本單元測試**
  - 使用 `testing`、`t.Run`、table-driven test。
  - 驗證新增後長度、`ID`、`Title`、`Done == false`，及原有元素未被破壞。
- **Part 2–3：Slice、Backing Array 與 Mutation**
  - 理解 Slice Header、`len`、`cap`、`append`、共享 Backing Array 與值覆寫。
  - 能分析多個 Slice 共用同一底層陣列時，修改會如何互相影響。
- **Part 4：DeleteTodo 測試**
  - table-driven test 涵蓋刪除第一筆、中間、最後、找不到 ID、空 Slice。
  - 使用 `slices.Equal` 比較資料，並新增 `TestDeleteTodo_MutatesInput`。
- **Part 5–6：Slice Ownership／Task 04 Clone Before Mutation**
  - 理解 `backup := original` 只複製 Slice Header；`slices.Clone` 會建立獨立的元素儲存空間，但屬於淺複製。
  - 完成 `TestDeleteTodo_ClonePreservesOriginal`：先複製再刪除，驗證回傳 `[A, C]` 且原始資料不變。
  - 修正 `TestDeleteTodo_MutatesInput`：分開檢查 `found` 與資料變更，避免假陽性。
  - **目前維持 `DeleteTodo` 可以修改輸入 Backing Array 的 API 契約**；呼叫端有保留原始資料需求時自行 Clone。
- **Part 7／Task 05：CompleteTodo**
  - 完成已存在、找不到、重複完成三種 table-driven test。
  - 理解 `range` 的 `item` 是 Struct 元素的值複製；修改 `item.Done` 不會改到 `todos[i].Done`。
- **Part 7／Task 06：GetNextID**
  - 完成空 Slice、連續 ID、不連續 ID、亂序 ID 的測試。
  - 理解 `maxID + 1` 可能重用已刪除 ID，並討論資料庫唯一性、歷史關聯與並行風險。
- **Part 7／Task 07：SearchTodos**
  - 完成完整匹配、部分匹配、大小寫不敏感、找不到、空關鍵字等五種測試。
  - 已決定產品規格：**空字串關鍵字代表不篩選，回傳全部 Todo**。
- **Part 7／Task 08：SearchTodos 的資料不變性**
  - 完成 `TestSearchTodos_DoesNotMutateInput`，驗證多筆匹配順序正確、搜尋不修改原始 Slice。
  - 理解目前 `Todo` 為只有值型別欄位的 Struct，搜尋結果的元素複製不與輸入共用元素儲存；未來若包含 Pointer、Map、Slice 欄位，須再考慮淺複製限制。

### 最近一次程式碼與測試狀態

截至 2026-10-29 最近一次上傳的 `todo_test.go`，已包含：

- `TestAddTodo`
- `TestDeleteTodo`
- `TestDeleteTodo_MutatesInput`
- `TestDeleteTodo_ClonePreservesOriginal`
- `TestCompleteTodo`
- `TestGetNextID`
- `TestSearchTodos`
- `TestSearchTodos_DoesNotMutateInput`

第 1 堂結束時，課堂回報 `go test -v ./todo` 與 `go vet ./todo` 已通過；當時在臨時副本以 Go 1.23.2 驗證，並未修改原始專案。**請在下次上課時重新檢查我實際上傳的檔案**，不要將舊測試結果當成新檔案已經通過。

先前提供的 `go-todo-cli-master.zip` 可能尚未包含後續單獨提交的最新版 `todo_test.go`，因此不能只以舊壓縮檔中的測試檔判斷進度。

### 下一次開課：第 2 堂／Part 1／Task 01 — 不要吞掉檔案讀取錯誤（尚未完成）

**問題：** 目前 `storage/storage.go` 的 `LoadTodos() []todo.Todo` 在讀取檔案失敗或 JSON 解析失敗時，都回傳空 Slice。這會把首次使用的正常情況與資料損壞、權限等錯誤混為一談；如果後續儲存，可能覆蓋原有資料。

**目標：** 修改 `LoadTodos` 的簽名為：

```go
func LoadTodos() ([]todo.Todo, error)
```

在 `storage/storage.go` 中：

1. 對 `os.ReadFile("todos.json")` 的錯誤進行判斷。
2. 用 `errors.Is(err, os.ErrNotExist)` 區分檔案不存在：此情況回傳空 Todo Slice 與 `nil` error，視為初次使用。
3. 其他讀取錯誤應回傳 error，不可默默當作空清單。
4. `json.Unmarshal` 失敗時回傳 error，不可吞掉 JSON 損壞。
5. 成功載入時回傳 `todos, nil`。

在 `main.go` 中：

1. 調整 `storage.LoadTodos()` 呼叫端以接收 `todos, err`。
2. 若載入出錯，顯示適當錯誤訊息並終止本次 CLI 執行；不應繼續執行可能覆寫資料的新增／儲存動作。
3. 暫不修改 `SaveTodos` 與其他 Todo domain 函式。

**教學重點：** Go `error` 回傳值、`nil` error、`errors.Is`、Error Propagation；`errors.New`、`fmt.Errorf("%w")` 與測試錯誤路徑留待後續階段深入。

**驗收：** 檢查 `storage/storage.go` 和 `main.go`，執行 `go test ./...`；請我提供兩份修改後檔案與測試結果，再進行 Code Review。只給必要提示／骨架，不要直接替我寫完整答案。

**思考題（尚未作答）：** 如果 `todos.json` 存在，但 JSON 缺少結尾 `}`，為什麼不能直接回傳空 Slice 並讓應用程式繼續執行？

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

- **有進度紀錄時：** 先閱讀本文「課堂進度紀錄」及當次提供的最新專案內容，從尚未完成的任務接續。
- **本次指定續課位置：** 第 2 堂「Go Error Handling」／Part 1／Task 01「不要吞掉檔案讀取錯誤」。第 1 堂已完成，不要重講或要求重做其 Task 04–08。
- **開課第一步：** 檢查最新 `storage/storage.go`、`main.go` 及相關測試，確認是否仍採用舊版 `LoadTodos() []todo.Todo`，再依實際狀態安排 Task 01；不依進度紀錄猜測新上傳檔案。
- **Task 01 完成後：** 先 Code Review 和驗證錯誤路徑，再循序學習 `errors.New`、`fmt.Errorf("%w")`、錯誤處理測試等；一次只給一個合理大小的任務。
- **若未來沒有任何進度紀錄：** 才從「第 1 堂：完整單元測試 + Slice 深入」開始，檢查專案結構並指定第一個小任務。
- 永遠以當次提供的專案內容為準，不要一次提供整堂課的完成版程式碼。
