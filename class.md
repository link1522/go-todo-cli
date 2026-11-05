我正在使用一個 `todo-cli` 專案學習 Go。請把我上傳的專案當成整門課的核心教材，先閱讀目前專案內容，再依照現有程式碼循序漸進教我 Go。

我的目標不是只學語法，而是透過同一個 Todo 專案逐步學會 Go 的實務開發方式，最後把它從簡單 CLI 演進成具有測試、Service、Repository、SQLite、REST API、Concurrency 與工程化結構的小型 Go Backend 專案。

## 課堂進度紀錄（最後更新：2026-11-05）

> **目前狀態：第 1 堂「完整單元測試 + Slice 深入」、第 2 堂「Go Error Handling」、第 3 堂「Struct 與 Method／Todo Service 基礎」均已完成。**
> **下次上課直接開始第 4 堂「Interface 與 Dependency Injection」／Part 1／Task 16「第一個 Interface」。**
> 第 4 堂尚未開始實作；不要重新布置或重講 Task 01–15 的內容。可在新課程必要時短暫回顧既有程式的關鍵 API。
> **任何進度及驗收結論都以當次使用者實際提供的最新專案內容為準；以下是續課定位，不代表新版本已自動整合。**

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

### 第 2 堂已完成內容

- **Part 1／Task 01：不要吞掉檔案讀取錯誤**
  - 將 `storage.LoadTodos()` 改為 `func LoadTodos() ([]todo.Todo, error)`。
  - 使用 `errors.Is(err, os.ErrNotExist)` 將「檔案不存在」視為首次使用，回傳空 Slice 與 `nil`。
  - 其他讀取錯誤與 JSON 解析錯誤向上傳遞；`main()` 載入失敗時停止，避免後續操作覆寫資料。
  - 已驗證三種情境：檔案不存在、合法 JSON、格式錯誤 JSON。
- **Part 2：Error Wrapping**
  - 學習 `fmt.Errorf()` 的 `%w` 與 `%v` 差異。
  - `storage.LoadTodos()` 對讀取錯誤使用 `read todos.json failed: %w`，對 JSON 解析錯誤使用 `parse todos.json failed: %w`。
  - 理解 `%w` 保留錯誤鏈，因此包裝後仍可使用 `errors.Is()` 判斷原始錯誤。
  - 理解 `%v` 只格式化錯誤文字，不保留可供 `errors.Is()`／`errors.As()` 遍歷的錯誤鏈。
- **Part 3：errors.As() 與 Pointer to Pointer**
  - 學習 `errors.As()` 可沿錯誤鏈尋找指定錯誤型別，並將找到的錯誤指派給目標變數。
  - 以 `var pathErr *os.PathError` 為例，理解 `pathErr` 是 `*os.PathError`，而 `&pathErr` 是 `**os.PathError`。
  - 已釐清 `errors.As(err, &notFoundErr)` 需要目標變數的地址，讓 `errors.As()` 能將找到的錯誤指標寫入該變數。
- **Part 4：Custom Error 與 DeleteTodo 整合**
  - 建立 `TodoNotFoundError`，包含 `ID int`，並以 Pointer Receiver 實作 `Error() string`。
  - 理解 Go 的隱式 Interface 實作：型別的方法集合符合 Interface 即自動實作，不需 `implements` 宣告。
  - 將 `DeleteTodo` 改為 `func DeleteTodo(todos []Todo, id int) ([]Todo, error)`；成功回傳更新後 Slice 與 `nil`，找不到時回傳原 Slice 與 `&TodoNotFoundError{ID: id}`。
  - `main.go` 使用 `errors.As(err, &notFoundErr)` 取得錯誤中的 Todo ID。
  - `TestDeleteTodo` 已改為驗證成功時 `err == nil`，失敗時驗證錯誤型別及 ID。
  - 已討論 Custom Error 的 Pointer Receiver、`*T`／`**T`，並以 Part 5 延伸到 Sentinel Error 與 `Unwrap()`。

### 第 2 堂／Part 5 已完成內容（Task 09）

- 定義固定錯誤值：`var ErrTodoNotFound = errors.New("todo not found")`。
- `TodoNotFoundError` 保留 `ID`、實作 `Error() string`，並透過 `Unwrap() error` 回傳 `ErrTodoNotFound`。
- 使用 `errors.Is()` 判斷「找不到 Todo」的錯誤語意；使用 `errors.As()` 取得 `*TodoNotFoundError` 與 ID。
- 使用 `fmt.Errorf("...: %w", err)` 包裝錯誤後，驗證 `errors.Is`／`errors.As` 仍能沿錯誤鏈工作；理解 `%v` 不保留錯誤鏈。
- 完成 `TestTodoNotFoundError_Error`、`TestTodoNotFoundError`、`TestTodoNotFoundErrorWrapped`。
- 已修正 `errors.As` 失敗後必須使用 `t.Fatal`／`t.Fatalf`，以免後續解參照 nil pointer 的測試風險。

### 第 3 堂已完成內容（Task 10–15）

- **Part 1／Task 10：Method 與 Pointer Receiver**
  - 在 `Todo` 上新增 `(todo *Todo) MarkDone()`，修改原始 `Done` 欄位。
  - `CompleteTodo()` 改成 `todos[i].MarkDone()`；理解 `range` 的 `item` 是值複製，不能藉由 `item.MarkDone()` 修改原始 Slice 元素。
  - 完成 `TestTodo_MarkDone`，既有 `CompleteTodo` 測試保留。
- **Part 2／Task 11：Zero Value 與 Constructor Pattern**
  - 新增 `NewTodo(id int, title string) Todo`，預設 `Done == false`；`AddTodo` 改用 `NewTodo`。
  - 完成 `TestNewTodo`，理解 `bool` 的零值為 `false`，也理解 Go 不要求所有 Struct 必須有 Constructor。
- **Part 3／Task 12：Value Receiver**
  - 新增 `(todo Todo) Matches(keyword string) bool`，以 `strings.ToLower`／`strings.Contains` 實作不區分大小寫搜尋。
  - `SearchTodos()` 改為呼叫 `item.Matches(keyword)`，保持原本搜尋結果與順序；空關鍵字仍視為匹配全部。
  - 完成 `TestTodo_Matches`，理解 Value Receiver 的值複製，以及 Receiver 類型一致性與 Method Set 之後仍需深入學習。
- **Part 4／Task 13：Service 與 Slice Ownership**
  - 新增 `todo/service.go`，包含 `Service{todos []Todo}`、`NewService(initial []Todo) *Service`、`Add(title string) Todo`、`List() []Todo`。
  - `NewService` 與 `List` 使用 `slices.Clone` 隔離呼叫端與 Service 的元素儲存空間；理解淺複製限制。
  - `Add` 使用 `GetNextID`、`NewTodo` 並更新 `s.todos`；測試涵蓋新增 ID、初始 Slice 隔離、List 回傳值隔離。
- **Part 5／Task 14：Service.Complete 與 Error Propagation**
  - 新增 `(s *Service) Complete(id int) error`，重用 `CompleteTodo`；存在時回傳 `nil`，不存在時回傳 `&TodoNotFoundError{ID: id}`。
  - 驗證完成成功、找不到 Todo 時 `errors.Is`／`errors.As`、重複完成等情境。
  - 已改善 `t.Fatalf` 的安全檢查、成功案例的 `err == nil` 驗證，以及 `Add` 實際寫入與長度的檢查。
- **Part 6／Task 15：Service.Delete 與 Service.Search**
  - 新增 `(s *Service) Delete(id int) error`，呼叫 `DeleteTodo`，成功時將更新後的 Slice Header 指派回 `s.todos`；錯誤時直接傳遞，不更新狀態。
  - 新增 `(s *Service) Search(keyword string) []Todo`，重用 `SearchTodos(s.todos, keyword)`，沿用大小寫不敏感、順序不變等搜尋契約。
  - 測試涵蓋刪除成功、刪除不存在的 ID／錯誤資訊／狀態不變，以及完整、部分、大小寫不敏感、無結果搜尋。
  - 最近一次 Code Review 判定 Task 15 核心功能通過；曾在暫存副本額外驗證搜尋結果 Ownership 與空關鍵字搜尋，兩項皆通過，但未必已寫入使用者的正式測試檔。

### 最近一次專案狀態與整合注意事項

- **最近提供的檔案**：`service.go`、`service_test.go`；先前另有新版 `todo.go`、`todo_test.go`。最初上傳的 `go-todo-cli.zip` 是較早的快照，未必含最新的 Service 檔案。下次以使用者當次提供的完整專案或最新檔案為準，不能直接假設 ZIP 已同步更新。
- 目前最後提交的 Service API：`NewService`、`Add`、`List`、`Complete`、`Delete`、`Search`。`Service` 的 `todos` 是未匯出欄位，`List` 回傳 Clone；`Search` 依賴現有 `SearchTodos` 建立獨立結果 Slice。
- Todo Domain：`Todo` 含 `ID`、`Title`、`Done`；已提供 `MarkDone`、`Matches`、`NewTodo`、`AddTodo`、`CompleteTodo`、`DeleteTodo`、`GetNextID`、`SearchTodos`、`ErrTodoNotFound`、`TodoNotFoundError`。
- **歷史測試紀錄**：最近一次對 Task 15 的審查，在整合暫存副本中排除舊版 `todo_test.go` 的重複 Service 測試後，`go test -count=1 ./...` 與 `go vet ./...` 通過。當時使用 Go 1.23.2，僅在暫存副本調整 `go.mod` 版本要求；原專案指定 Go 1.25.5，未在該版本驗證。這是過去的驗收結果，**不是對下次新上傳版本的保證**。
- **非阻擋事項**：確認本機 `todo_test.go` 已移除搬至 `service_test.go` 的三個同名 Service 測試，避免重複宣告；可在適當時機補上正式的 Search Ownership 測試，並以 `gofmt` 統一格式。不需因這些小清理倒退重講第 3 堂。

### 下一次開課：第 4 堂／Part 1 — Interface 與 Dependency Injection（Task 16）

**下次上課直接從第 4 堂 Part 1 開始。** 第 1–3 堂都已完成，Task 16 尚未提交任何實作。

1. 先根據當次實際提供的最新專案，確認 `Service.List() []Todo` 存在，以及 Service 使用 Pointer Receiver；僅做必要的狀態檢查。
2. 用目前的 Todo Service 解釋 Interface 描述行為契約、隱式實作、small interface 與 Method Set，不需要先引入 Repository／資料庫。
3. **Task 16：新增 `todo/stats.go`、`todo/stats_test.go`**，定義：

   ```go
   type TodoReader interface {
       List() []Todo
   }

   func CountPending(reader TodoReader) int {
       // TODO: 回傳 Done == false 的 Todo 數量
   }
   ```

4. 要讓學員自行實作 `CountPending`，測試至少涵蓋：使用真實 `*Service`；建立有 `List() []Todo` 的 Fake Reader；使用 Fake Reader 計算；空集合回傳 0。
5. 思考題：若 `List()` 定義為 `(s *Service) List() []Todo`，`var reader TodoReader = &Service{}` 可以成立，但 `var reader TodoReader = Service{}` 不成立；以 Method Set 解釋原因。
6. 按照既有節奏「問題 → 原因 → Go 概念 → 實作」安排一個適中任務，不要直接提供整份 `stats.go` 與 `stats_test.go` 的完成答案。待使用者提交後，再作 Code Review。
7. 第 4 堂後續方向依序為 Interface / Implicit Implementation / Method Set / Dependency Injection，再考慮小型 Repository Pattern、`TodoRepository`、`JSONRepository`；不要在 Task 16 提前整套實作。

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
   - error、Error Propagation
   - errors.Is、errors.As
   - errors.New、Sentinel Error
   - fmt.Errorf("%w")、Error Wrapping、Unwrap
   - Custom Error 與自訂錯誤型別
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

- **有進度紀錄時：** 先閱讀本文「課堂進度紀錄」及使用者當次提供的最新專案內容，從尚未完成的任務接續；不要直接假設先前檔案仍是最新版本。
- **本次指定續課位置：** **第 4 堂「Interface 與 Dependency Injection」／Part 1／Task 16「第一個 Interface」**。第 1、2、3 堂及 Task 01–15 已完成，不重新授課或要求重做。
- **開課第一步：** 確認當次專案的 `Service.List() []Todo`、`Service` Receiver 和相關測試；如果遇到歷史快照與新檔案不一致，說明差異，必要時以最新檔案為準。
- **接著開始 Task 16：** 解釋 `TodoReader` 小介面與 `*Service` 的隱式實作；安排 `stats.go`／`stats_test.go` 的 `CountPending` 與 Fake Reader 練習。實作由學員完成，提交後再驗收。
- **不得提前跳到 Repository／SQLite：** 先完成 Interface、Method Set、基本 Dependency Injection 的學習；維持每次只進行一個適中任務。
- **舊版整合提示：** 如因 `todo_test.go`、`service_test.go` 含同名測試導致編譯失敗，提醒刪除舊版重複測試，而不是要求重新完成已通過的 Service 設計；曾有 CRLF／gofmt 格式差異。
- **若未來沒有任何進度紀錄：** 才從「第 1 堂：完整單元測試 + Slice 深入」開始。
- 永遠以當次提供的專案內容為準，保留循序漸進的教學節奏，不要一次提供整堂課的完成版程式碼。
