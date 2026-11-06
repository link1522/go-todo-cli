我正在使用一個 `todo-cli` 專案學習 Go。請把我上傳的專案當成整門課的核心教材，先閱讀目前專案內容，再依照現有程式碼循序漸進教我 Go。

我的目標不是只學語法，而是透過同一個 Todo 專案逐步學會 Go 的實務開發方式，最後把它從簡單 CLI 演進成具有測試、Service、Repository、SQLite、REST API、Concurrency 與工程化結構的小型 Go Backend 專案。

## 課堂進度紀錄（最後更新：2026-11-06）

> **目前狀態：第 1 堂「完整單元測試 + Slice 深入」、第 2 堂「Go Error Handling」、第 3 堂「Struct 與 Method／Todo Service 基礎」已完成；第 4 堂「Interface 與 Dependency Injection」的 Task 16–20 均已完成驗收。**
> **下次上課直接從「第 4 堂｜Task 21：資料一致性（Data Consistency）」開始。**
> Task 21 已講解需求、提供實作方向與思考題，但 **尚未收到 `AddAndSave()` 實作與測試，也尚未回答 Task 21 的思考題**。不要再從 Task 16 開始，亦不要把 Task 21 標記為完成。
> 可在新課開始時簡短回顧 Task 20 的 `Service`／Repository API，但不要重新布置 Task 01–20。
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

### 第 4 堂已完成內容（Task 16–20）

- **Part 1／Task 16：第一個 Interface — `TodoReader`**
  - 在 `todo/stats.go` 定義 `type TodoReader interface { List() []Todo }`，完成 `CountPending(reader TodoReader) int`，統計 `Done == false`。
  - 在 `todo/stats_test.go` 以真實 `*Service`、Fake Reader、空集合驗證結果；學員已理解 small interface 與 implicit implementation。
  - 釐清 Method Set：Interface 本身不會定義 Pointer Receiver；現有 `List()` 為 `func (s *Service) List() []Todo`，因此 `*Service` 滿足 `TodoReader`，`Service` 值型別不滿足。
  - `fakeTodoReader` 改用未匯出名稱，`List()` 回傳 `slices.Clone`。
- **Part 2／Task 17：Constructor Injection — `StatsService`**
  - 新增 `StatsService{reader TodoReader}`、`NewStatsService(reader TodoReader) *StatsService`、`PendingCount() int`，後者重用 `CountPending(s.reader)`。
  - 以真實 `*Service` 和 Fake Reader 測試通過；學員了解依賴 Interface 可降低耦合，且不同實作的初始化方式可不同。
  - Code Review 建議 Constructor 使用具名欄位 `&StatsService{reader: reader}`，屬風格建議、非必要修正。
- **Part 3／Task 18：Repository 契約**
  - 在 `todo/repository.go` 定義：

    ```go
    type TodoRepository interface {
        Load() ([]Todo, error)
        Save(todos []Todo) error
    }
    ```

  - 學員理解具體型別只要具有正確 Method Set 就能隱式實作介面；Package-Level Function 不等於某個型別的 Method。
  - 釐清空集合與 I/O／JSON 解析錯誤是不同情況，因此 Repository 讀取需要回傳 `error`。
- **Part 3／Task 19：`storage.JSONRepository` 與檔案測試**
  - 在 `storage` Package 建立 `JSONRepository{path string}`、`NewJSONRepository(path string)`、`Load()`、`Save([]todo.Todo) error`；由 Constructor 注入路徑，保留讀寫 JSON 的責任。
  - `Load()` 使用 `os.ReadFile`／`json.Unmarshal`；檔案不存在視為空資料，其他錯誤使用 `fmt.Errorf(... %w ...)` 包裝。
  - `Save()` 使用 `json.Marshal`／`os.WriteFile`；可選改善：寫入錯誤加上檔案路徑上下文。
  - 測試使用 `t.TempDir()`；已補齊 `Load()` 錯誤檢查及不同檔案路徑隔離驗證：A 儲存、B 儲存後重新讀取 A，確認不被覆寫。Task 19 已正式通過。
  - 可選的編譯期介面檢查：`var _ todo.TodoRepository = (*JSONRepository)(nil)`；不是必須。
- **Part 4／Task 20：將 Repository 注入 Todo Service**
  - 保留舊的 `NewService(initial []Todo)`，新增 `Service.repo TodoRepository`、`NewServiceWithRepository(repo TodoRepository) (*Service, error)`、`(s *Service) Save() error`。
  - Constructor 對 `nil` Repository 回傳 error；`repo.Load()` 失敗時向上傳遞 error；成功時使用 `slices.Clone(todos)` 建立獨立 Service 狀態。
  - `Save()` 對無 Repository 的 Service 回傳 error，且現在呼叫 `s.repo.Save(slices.Clone(s.todos))`，避免將內部 Backing Array 暴露給 Repository。
  - Fake Repository 可模擬 `loadError`、`saveError`；原本測試誤將 `slices.Equal(wanted, got)` 的條件寫反，已修正。
  - 最後一版測試 Fake Repository 故意以 `return r.todos, nil` 與 `r.todos = todos` 直接共享輸入，以驗證 Service 邊界複製有意義；`service.go` 的 Load／Save 兩側都已使用 Clone。
  - 全專案既有測試和 `go vet` 通過；審查者在暫存副本額外驗證 Load 後外部修改、Save 後 Service 修改，以及 Repository 在 Save 中修改輸入等情境，均通過；**這些額外測試未必已寫入學員的正式 `service_test.go`**。
  - 學員正確辨識共享 Backing Array 風險；需維持區分：Fake Repository 的 Slice 指派不是磁碟持久化，隔離成立是因 `slices.Clone`，不是因「已經 Save 就不會再變」。
  - 目前的 `Add()` 與 `Save()` 是分開的操作；如果 Save 失敗，記憶體狀態仍可能已變動。此一致性議題留到 Task 21。
  - Task 20 已正式通過，**下次不要重新驗收或要求重做**，除非使用者提交了實質不同的新版本。

### 最近一次專案狀態與整合注意事項

- **最新上傳紀錄（2026-11-06）**：最近一輪提交了新版 `/mnt/data/service.go` 和 `/mnt/data/service_test.go`，並通過 Task 20 審查。此前分批提交 `stats.go`、`stats_test.go`、`repository.go`、`JSONRepository.go`、`JSONRepository_test.go`。這些是歷次個別檔案；最初上傳的 `go-todo-cli.zip` 是 2026-11-05 的舊快照，**不能假設已自動同步所有新檔案**。
- **最新 Service API**：`NewService(initial []Todo)`、`NewServiceWithRepository(repo TodoRepository)`、`Add`、`List`、`Complete`、`Delete`、`Search`、`Save`。Service 包含 `todos []Todo` 與 `repo TodoRepository`；對建構載入的 Slice 及 Save 傳入的 Slice 都做 Clone。
- **其它新 API**：`TodoReader`、`CountPending`、`StatsService`、`NewStatsService`、`TodoRepository`、`storage.JSONRepository`、`NewJSONRepository`；Package 模組路徑依 `go.mod` 為 `todo-cli`，以當次實際版本為準。
- **Fake Repository 注意事項**：最後一版測試中的 `fakeRepository.Load()` 直接回傳內部 Slice，`Save()` 直接指派傳入 Slice。這是學員用以測試 Service Ownership 的刻意設計，不能誤判為缺陷；Service 必須自己守住兩側的資料所有權邊界。
- **尚可補強的測試**：目前正式 `TestTodoService_ModifyLoadedTodos` 主要驗證 `List()` 的回傳值不影響 Service，尚未直接持久記錄「修改 Repo 原始 Slice 不影響 Service」的獨立測試。此項屬非阻擋改善，不要因此否定已通過的 Task 20。
- **最近驗證結果與限制**：Task 20 最新整合暫存副本中 `go test -count=1 ./...` 與 `go vet ./...` 通過，另外的 Ownership 暫存測試也通過；驗證環境為 Go 1.23.2，需暫時調整副本的 `go.mod` 版本宣告，專案原訂 Go 1.25.5 尚未在該版本驗證。這是歷史紀錄，**不是下次新上傳版本的自動保證**。
- **仍有的非阻擋建議**：可考慮 `gofmt`、Go 習慣的小寫底線檔名（`json_repository.go`）、`Save()` 的錯誤上下文及更完整的負向測試；不要因此倒退重教 Task 16–20。

### 下一次開課：第 4 堂／Task 21 — 資料一致性（Data Consistency）

**下次直接從 Task 21 開始，不要重講 Task 16–20。Task 21 尚未實作。** 學員先前的思考方向是「Save 失敗就回復前一狀態」，已引導其思考更容易維護的方案：**先在候選狀態修改 → 儲存候選狀態 → 成功後才提交記憶體狀態**，避免需要回滾已修改的原始狀態。

1. **問題背景**：目前 `service.Add("Learn Go")` 後再呼叫 `service.Save()`；若 Save 失敗，記憶體已有新 Todo，但資料來源可能沒有成功儲存，造成不一致。
2. **Task 21 實作目標**：在 `todo/service.go` 新增 `func (s *Service) AddAndSave(title string) (Todo, error)`；**保留既有 `Add()`、`Save()` 方法及簽章**，不要在此階段整套重構。
3. **預計流程**：
   - 檢查 Repository 是否存在，沒有就回傳 error。
   - `candidate := slices.Clone(s.todos)` 建立候選資料狀態。
   - 使用已有的 `GetNextID(candidate)`、`NewTodo`，只把新 Todo 加到候選狀態。
   - 呼叫 `s.repo.Save(slices.Clone(candidate))`，保護候選狀態不被 Repository 意外修改。
   - 儲存失敗時回傳 error，不改 `s.todos`；只有儲存成功才以候選狀態更新 `s.todos` 並回傳新 Todo。
4. **至少需要的測試**：
   - 成功時新增 Todo，Service 和 Fake Repository 都有新資料。
   - Fake Repository Save 失敗時 `AddAndSave()` 回傳 error，Service 狀態不變、Repository 原資料不變。
   - 沒有 Repository 時回傳 error，Service 不變；可額外測試 Repository 故意修改傳入 Slice、原始資料 Ownership。
5. **重要界線**：上述設計可防止「已知 Save 失敗後仍更新 Service 記憶體」，但不能單靠這個流程保證實際磁碟儲存的原子性、交易性、並發安全或失敗後磁碟資料完全沒有改變。相關檔案安全寫入／資料庫 Transaction 後續再教。
6. **尚待學員回答的思考題**：**「為什麼『先建立候選狀態 → 儲存成功 → 更新記憶體』通常比『先修改記憶體 → 儲存失敗再 Rollback』更容易維護？」** 下次可先讓學員作答，再檢查實作；不要直接給完整完成答案。
7. 沿用「問題 → 原因 → Go 概念 → 實作 → 學員提交 → Code Review」的步調，以現有 Todo 專案為唯一實作基礎。先確認新上傳檔案的當前狀態，再安排後續工作。

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
