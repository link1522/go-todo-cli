我正在使用一個 `todo-cli` 專案學習 Go。請把我上傳的專案當成整門課的核心教材，先閱讀目前專案內容，再依照現有程式碼循序漸進教我 Go。

我的目標不是只學語法，而是透過同一個 Todo 專案逐步學會 Go 的實務開發方式，最後把它從簡單 CLI 演進成具有測試、Service、Repository、SQLite、REST API、Concurrency 與工程化結構的小型 Go Backend 專案。

## 課堂進度紀錄（最後更新：2026-11-02）

> **目前狀態：第 1 堂「完整單元測試 + Slice 深入」已完成；第 2 堂「Go Error Handling」已完成 Part 1–4。**
> **下次上課從第 2 堂／Part 5「Sentinel Error 與 Custom Error 的設計取捨」開始。**
> 不要重講已完成的 Error、`errors.Is`、`fmt.Errorf("%w")`、`errors.As`、Pointer to Pointer 與 Custom Error 基礎。
> 本紀錄用於定位進度；程式碼仍以使用者下次提供的最新專案版本為準。

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
  - 已討論 Custom Error 的 Pointer Receiver、`*T`／`**T`，以及 Sentinel Error、`Unwrap()` 可以如何與 Custom Error 搭配；Sentinel Error 尚未實作。

### 最近一次程式碼狀態與待確認事項

最近一次收到的檔案為 `main.go`、`todo.go`、`todo_test.go`；另有先前更新的 `storage.go`。目前可確認：

- `storage.LoadTodos()` 回傳 `([]todo.Todo, error)`，使用 `errors.Is` 區分檔案不存在，並以 `%w` 包裝其他讀取錯誤及 JSON 解析錯誤。
- `todo.DeleteTodo()` 已回傳 `([]Todo, error)`，找不到項目時回傳 `*TodoNotFoundError`。
- `main.go` 已使用 `errors.As()` 取得 `TodoNotFoundError.ID`。
- `todo_test.go` 已加入 `errors.As()` 型別與 ID 驗證，並將 Slice 測試改為檢查 `err`。
- **待清理的小地方：** `TestDeleteTodo_ClonePreservesOriginal` 已先檢查 `err != nil` 並在錯誤時 `t.Fatalf`，後面又對同一個 `err` 呼叫 `errors.As()`；該第二次檢查是冗餘的，可以移除。
- **測試狀態：** 尚未收到這一版整合後 `go test ./...` 的最新執行結果；不可宣稱最新版本測試已通過。下次可先確認這項小清理與測試結果，不必重講前面課程。

### 下一次開課：第 2 堂／Part 5 — Sentinel Error 與 Custom Error 的設計取捨

**下次上課直接從本 Part 開始。** Part 1–4 不需重講。

預計討論與實作：

1. Sentinel Error 是什麼，以及如何使用 `errors.New()` 定義固定錯誤值。
2. `errors.Is()` 適合判斷固定錯誤語意的原因。
3. Custom Error 與 `errors.As()` 適合攜帶、取得結構化錯誤資訊的原因。
4. Custom Error 是否可以同時支援 `errors.Is()` 與 `errors.As()`。
5. `Unwrap() error` 如何讓 Custom Error 與 Sentinel Error 建立錯誤鏈關係。
6. 依 Todo CLI 的「Todo 不存在」情境，評估何時只用 Sentinel Error、何時使用 Custom Error，以及何時兩者搭配。

**下次開課的思考題（尚未作答）：**

假設已有 `var ErrTodoNotFound = errors.New("todo not found")`：

- 需求 A：只需要知道刪除失敗是因為 Todo 不存在，適合使用 `errors.Is()` 還是 `errors.As()`？
- 需求 B：除了知道 Todo 不存在，還需要取得使用者嘗試刪除的 ID，適合使用哪一種？
- 如果同時需要判斷錯誤種類及取得 ID，是否有必要讓 Custom Error 實作 `Unwrap()`？為什麼？

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

- **有進度紀錄時：** 先閱讀本文「課堂進度紀錄」及當次提供的最新專案內容，從尚未完成的任務接續。
- **本次指定續課位置：** 第 2 堂「Go Error Handling」／Part 5「Sentinel Error 與 Custom Error 的設計取捨」。第 1 堂及第 2 堂 Part 1–4 已完成，不要重講或要求重做。
- **開課第一步：** 從第 2 堂 Part 5 開始。先依當次最新檔案做必要的狀態確認；若尚未移除 `TestDeleteTodo_ClonePreservesOriginal` 中冗餘的 `errors.As()` 檢查，提醒清理並確認 `go test ./...` 結果，但不要因此重講 Part 1–4。
- **Part 5 開始後：** 先釐清 Sentinel Error、Custom Error、`errors.Is()`、`errors.As()` 與 `Unwrap()` 的設計取捨，再以 Todo CLI 為情境安排一個合理大小的實作任務；一次只給一個任務。
- **若未來沒有任何進度紀錄：** 才從「第 1 堂：完整單元測試 + Slice 深入」開始，檢查專案結構並指定第一個小任務。
- 永遠以當次提供的專案內容為準，不要一次提供整堂課的完成版程式碼。
