我正在使用一個 `todo-cli` 專案學習 Go。請把我上傳的專案當成整門課的核心教材，先閱讀目前專案內容，再依照現有程式碼循序漸進教我 Go。

我的目標不是只學語法，而是透過同一個 Todo 專案逐步學會 Go 的實務開發方式，最後把它從簡單 CLI 演進成具有測試、Service、Repository、SQLite、REST API、Concurrency 與工程化結構的小型 Go Backend 專案。

## 課堂進度紀錄（最後更新：2026-11-12）

> **目前狀態：第 1–4 堂已完成；第 5 堂 CLI Application 進行中，Task 22–28 均已完成驗收。** `list`、`add` 子命令、`parseCommand()`、`runCLI()`、`runApp()`，以及 stdout／stderr 分流、Exit Code 與必要測試均已完成。
> **下次上課直接從「第 5 堂 — Task 29：新增 `done` 子命令」開始。** 不要重講 Task 01–28、不要重新要求 Task 28 的測試補強。
> **Task 29 起點**：先延續學員尚未回答的設計思考題——`done 3` 的 ID 應在 `parseCommand()` 驗證後保留字串（A），還是修改 `CLICommand` 以儲存已解析的整數 ID（B）？先讓學員回答與討論取捨，再開始小步實作 `parseCommand()` 的 `done` 分支；**暫不要求同時實作 `runCLI()`**。
> **版本注意**：本檔是續課進度紀錄；舊 `go-todo-cli.zip` 尚未自動整合後續單獨上傳的 `main.go`／`main_test.go`。任何進度與驗收都須以當次能取得的最新程式為準。

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

### 第 4 堂／Task 21 已完成內容 — 資料一致性（Data Consistency）

- **Task 21：`AddAndSave()` 設計與實作已通過驗收。** 在 `todo/service.go` 新增 `func (s *Service) AddAndSave(title string) (Todo, error)`；保留原有 `Add()`、`Save()` 與 `TodoRepository` 介面不變。
- **流程**：先檢查 Repository → `candidate := slices.Clone(s.todos)` → 使用 `NewTodo(GetNextID(candidate), title)` 建立新 Todo 並加入候選 Slice → `s.repo.Save(slices.Clone(candidate))` → Save 成功後才更新 `s.todos` 並回傳新 Todo。失敗時以 `fmt.Errorf("save todos failed: %w", err)` 保留錯誤鏈，且不提交候選狀態。
- **測試**：已提交 `TestTodoService_AddAndSave`、`TestTodoService_AddAndSaveFailed`，並在 `TestTodoService_MissRepo` 涵蓋無 Repository 情境；已檢查成功時 Service／Repository 皆含新增 Todo、失敗時兩邊狀態維持原樣，相關測試通過。
- **Slice Ownership**：學員已能正確說明 `slices.Clone` 的目的，以及 Slice Header 與 Backing Array 的共享問題。Repository 收到 `slices.Clone(candidate)`，因此在目前 `Todo` 欄位為值型別的前提下，Repository 修改自己持有的 Slice 不會直接修改 Service 的資料。
- **簡化結論**：學員正確回答：將最後的 `s.todos = slices.Clone(candidate)` 改為 `s.todos = candidate` 仍不會造成 Service／Repository 共用 Backing Array；因為傳給 `Save()` 的是另一份 Clone。最後一次 Clone 可省略，前兩個 Clone 應依資料所有權邊界保留。**上傳版本仍使用第三次 Clone；這是安全但冗餘的寫法，不必強迫回頭修改。**
- **思考題已回答**：學員指出「先修改記憶體再 Save」可能使其他使用者／操作在儲存成功前讀到新資料，甚至儲存失敗後仍讀到未持久化狀態。已補充：候選狀態解決的是失敗後不提交 Service 狀態，**不等同** Mutex／並發安全，也不保證檔案寫入原子性或資料庫交易性。
- **非阻擋補強**：建議將 `TestTodoService_AddAndSaveOwnership` 正式加入測試，以驗證 Save 成功後修改 `fakeRepository.todos[0]` 不會影響 `service.List()`。課堂已提供此測試範例，但**尚未確認學員已加入正式測試檔**。此項不影響 Task 21 已驗收完成的結論。

### 最近一次專案狀態與整合注意事項

- **最新提交（2026-11-12）**：學員先後上傳新版 `main.go`、多版 `main_test.go`，完成 Task 28 與相關測試補強。`main()` 的子命令分支會呼叫 `runApp()` 並透過 `os.Exit(exitCode)` 回傳程序狀態；`runApp()` 依序執行 `parseCommand`、建立 JSONRepository／Service、`runCLI`，錯誤輸出 stderr 並回傳 1，成功回傳 0；無參數時保留 `runInteractive()`。最新 `main_test.go` 已補強成功輸出、Fake Repository 狀態、`errors.Is`、儲存失敗不顯示成功訊息與各個子測試獨立的 Buffer。**續課應優先使用這些單獨上傳的最新程式檔**，不要以舊 ZIP 覆蓋。
- **歷史提交（2026-11-07）**：學員曾另外上傳 `service.go` 與 `service_test.go`，完成 Task 21 `AddAndSave()`。此記錄為歷史背景，非下次課程起點。
- **最新 Service API**：`NewService(initial []Todo)`、`NewServiceWithRepository(repo TodoRepository)`、`Add`、`List`、`Complete`、`Delete`、`Search`、`Save`、`AddAndSave`。Service 包含 `todos []Todo` 與 `repo TodoRepository`。
- **其他既有 API**：`TodoReader`、`CountPending`、`StatsService`、`NewStatsService`、`TodoRepository`、`storage.JSONRepository`、`NewJSONRepository`；Module 路徑為 `todo-cli`，依最新實際專案為準。
- **Fake Repository 行為**：`Load()` 直接回傳 `r.todos`，`Save()` 成功時直接指派傳入的 Slice；此處是有意設計，用來測試 Service 的 Ownership 保護，不必把 Fake 一律改成 Clone。
- **前次驗證**：第 4 堂 Task 21 的 Code Review 回覆表示，`go test -count=1 ./...` 與 `go vet ./...` 通過；因執行環境為 Go 1.23.2，僅在暫存副本調整 `go.mod` 版本宣告，原專案 Go 1.25.5 宣告未變。這是歷史驗證紀錄，不保證下次修改後仍通過。
- **後續注意**：下次從 Task 29 的 `done` 子命令開始；先核對最新程式與學員的設計選擇，再進入 `parseCommand()` 的增量修改，不要因舊 ZIP 缺少 Task 28 而倒退重教。

### 第 5 堂原始學習路線（歷史規劃；以本文最末 2026-11-12 續課更新為準）

**Task 22–28 已完成，不要再把下列初期規劃當成下一次作業。** 本節保留第 5 堂的整體演進目標供參考。

1. **開課先檢查**：使用者當次提供的最新 `main.go`、`go.mod`，確認目前命令互動流程、Service／Repository 的初始化方式、錯誤輸出，以及已有的測試。不假設歷史 ZIP 是最新版。
2. **第一階段主題**：先理解 `os.Args`、CLI 子命令（subcommand）與參數驗證；必要時介紹 `flag`／`flag.FlagSet`，但依當前程式複雜度漸進安排。
3. **逐步目標**：讓使用者自己將操作發展為 `todo add "Learn Go"`、`todo list`、`todo done 1`、`todo delete 1`、`todo search "Go"`；不可一次要求完成全部命令。
4. **後續主題**：標準輸入／輸出、`stderr`、exit code、測試 CLI 行為與分離命令解析和 Domain Service 的責任。先用 Go Standard Library，不急著導入 Cobra 等第三方 Framework。
5. **教學流程**：先解釋當前 `main.go` 的問題及設計選擇，再交付小型任務／提示／必要的 Skeleton；等學員提交後才 Code Review。延續「問題 → 原因 → Go 概念 → 實作 → 驗收」節奏。
6. **注意界線**：`AddAndSave` 的基本資料一致性已學過；並發安全、原子檔案寫入、SQL Transaction 留待日後自然遇到時再深入，不要在第 5 堂第一個任務硬塞。

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
- **本次指定續課位置：** **第 5 堂 — Task 29「新增 `done` 子命令」**。Task 01–28 已完成核心驗收，不要重新授課或要求重做。
- **開課第一步：** 核對當次可取得的最新 `main.go`／`main_test.go`（單獨上傳檔可能比 ZIP 更新），確認 `parseCommand()`、`runCLI()` 與 `runApp()` 的現況。
- **接著開始 Task 29：** 從已提出但尚未回答的 ID 建模設計思考題（A：驗證後仍保留字串、執行時再轉換；B：`CLICommand` 保存已轉換整數）續課。討論選擇後僅布置 `parseCommand()` 對 `done` 的解析與驗證，不要直接交付完整答案或同時要求整合 `runCLI()`。
- **不提前跳到 SQLite 或 REST API：** 先完成 CLI 的命令解析、參數驗證與責任分層，保持單一課程階段的適當工作量。
- **舊版整合提示：** 如 `todo_test.go`、`service_test.go` 含同名測試導致編譯失敗，提醒移除重複舊測試，而不是重做已通過的 Service；留意 CRLF／gofmt 差異。
- **若未來沒有任何進度紀錄：** 才從「第 1 堂：完整單元測試 + Slice 深入」開始。
- 永遠以當次提供的專案內容為準，循序漸進教學，不要一次提供整堂課的完成版程式碼。


---

## Latest continuation update (2026-11-09) — Lesson 5 / CLI Application

### Completed in this session: Tasks 22–25

- **Task 22 — CLI subcommand entry point:** introduced `os.Args`, command selection, and the first `list` subcommand. Extracted command handling from `main()` into `runCLI(args []string, ...)` so it can be called directly in tests.
- **Task 22.1 — Argument validation:** added the empty-argument guard to prevent `args[0]` from panicking; validates the command and its argument count before initializing the Service.
- **Task 22.2 — Output abstraction:** changed `printTodos` to accept `io.Writer`, use `fmt.Fprintf`, and receive `os.Stdout` in normal execution. The learner understands that `*bytes.Buffer` satisfies `io.Writer` because its method set includes the pointer-receiver `Write` method; assigning a `bytes.Buffer` value to `io.Writer` does not trigger automatic address-taking.
- **Task 23 — Unit test for CLI output:** created `TestPrintTodos` using `bytes.Buffer` and exact output comparison, including newline characters.
- **Task 24 — Table-driven tests:** refactored `TestPrintTodos` into cases for an empty slice, one unfinished Todo, one completed Todo, and multiple Todos. Uses `t.Run(tt.name, ...)`; the learner correctly explained that named subtests make failures easier to locate and can be run individually.
- **Task 25 — Fake Repository and CLI tests:** added a local `fakeRepository` implementing `todo.TodoRepository`; `TestRunCLI` verifies `list` output through `bytes.Buffer`; `TestRunCLI_Errors` verifies missing command, unknown command, and extra arguments. The learner correctly identified that Fake Repository can supply deterministic Todo data without accessing the real JSON file.
- **Prior verification note:** the previous review reported `go test ./...` and `go vet ./...` passing in a temporary copy. The available toolchain was Go 1.23.2 while the project declares Go 1.25.5, so only the temporary copy's Go version declaration was adjusted for that check; the original project files were not changed by that verification.

### Current submitted CLI state (based on the separately uploaded `main.go` and `main_test.go`)

- `main()` keeps the interactive mode when no command-line arguments are provided. Otherwise it creates `storage.NewJSONRepository("todos.json")` and calls `runCLI(os.Args[1:], repo, os.Stdout)`.
- Current signature: `runCLI(args []string, repo todo.TodoRepository, w io.Writer) error`.
- `runCLI` validates an empty argument list, currently recognizes only `list`, rejects extra arguments for `list`, constructs `todo.Service` with `todo.NewServiceWithRepository(repo)`, and prints `service.List()` through the supplied writer.
- `runInteractive()` remains the earlier Scanner-based implementation and still uses the legacy storage/domain functions directly. It has not yet been migrated to the Service-based CLI path.
- `main_test.go` contains `fakeRepository`, the four `TestPrintTodos` table cases, a successful `TestRunCLI` case, and table-driven `TestRunCLI_Errors` cases.
- The latest `main.go` and `main_test.go` were uploaded separately in this conversation. The `go-todo-cli.zip` source snapshot may not contain these exact latest files. At the next lesson, inspect the files the user provides and do not assume the ZIP's source files are newer than the separately uploaded files.

### Design decision recorded for Task 26

- The learner prefers `runCLI` to receive `*todo.Service` rather than `todo.TodoRepository`.
- This is accepted for the current Todo CLI: the CLI should call the application's Service API and should not construct the Service or know which Repository backs it. Do not add another CLI-facing Interface solely for hypothetical future replacement; introduce one only if a concrete need appears.
- Tests can still use `fakeRepository`: construct a real Service with `todo.NewServiceWithRepository(fakeRepo)`, then inject that Service into `runCLI`.

### Historical Task 26 plan (completed; not the next lesson)

**Archive only:** Task 26 has already passed; do not assign this section again. The bullets below document the original Task 26 requirements, not the next lesson.

1. Extract command parsing and validation into `parseCommand(args []string) (string, error)`. It should reject missing commands, unknown commands, and invalid argument counts without initializing the Repository or Service.
2. Change the command runner to `runCLI(command string, service *todo.Service, w io.Writer) error`. It should execute the parsed command, call Service methods, and write output; it should neither receive a Repository nor construct a Service.
3. Keep dependency assembly in `main()`: parse and validate first; only after successful parsing create `JSONRepository`, create `todo.Service` with `todo.NewServiceWithRepository`, then call `runCLI`.
4. The order matters because `NewServiceWithRepository` loads data during construction. Invalid commands should be rejected before any repository load is attempted.
5. In tests, create a Service backed by `fakeRepository`, then test `runCLI` with that Service and a `bytes.Buffer`. Test `parseCommand` independently with table-driven cases.
6. Keep the current interactive `runInteractive()` path unchanged unless a specific Task 26 requirement makes a change necessary. The focus is the subcommand path's dependency direction and responsibility boundaries.



---

## 歷史續課更新（2026-11-10）— 第 5 堂 Task 26–27 完成，當時預定 Task 28（已完成）

### 已完成：Task 26 — CLI 與 Service 解耦

- 將命令解析獨立為 `parseCommand()`，先驗證命令，之後才建立 Repository 和 Service，避免無效命令觸發 Repository Load。
- `runCLI()` 改成由外部傳入 `*todo.Service` 和 `io.Writer`；不在 `runCLI()` 建立 Repository 或 Service。
- 使用 `fakeRepository` 建立真的 Service 來測試 `runCLI()`，並獨立測試命令解析。
- **Task 26 已完成驗收**；無需要求學員重做。

### 已完成：Task 27 — 帶參數的 CLI 子命令

- 加入 `type CLICommand struct { Name string; Args []string }`；`parseCommand(args []string) (CLICommand, error)` 將命令名稱與參數分開回傳。
- `parseCommand()` 支援 `list`（零個額外參數）與 `add "Learn Go"`（剛好一個非空白標題）；拒絕缺少命令、未知命令、參數過多、`add` 缺少標題、空字串與只有空白的標題。
- 使用 `strings.TrimSpace(args[1]) == ""` 驗證純空白標題，但保留原始標題字串作為 `CLICommand.Args[0]`。
- `runCLI(command CLICommand, service *todo.Service, w io.Writer) error` 執行 `add` 時呼叫 `service.AddAndSave(command.Args[0])`，成功時透過 `fmt.Fprintf(w, ...)` 輸出新增的 Todo ID 與標題，失敗直接回傳 error；`list` 仍透過 `printTodos(w, service.List())` 輸出。
- `main()` 在有命令列參數時先 `parseCommand(os.Args[1:])`，成功後才建立 `storage.NewJSONRepository("todos.json")` 和 `todo.NewServiceWithRepository(repo)`，再呼叫 `runCLI`；無參數時仍進入原有的 `runInteractive()`。
- `main_test.go` 已有 `TestPrintTodos`、`TestParseCommand`、`TestParseCommand_Failed`、`TestRunCLI_List`、`TestRunCLI_Add`、`TestRunCLI_Add_SaveFailed`，並使用 `fakeRepository` 與 `bytes.Buffer`。
- **上一輪 Code Review 結論：Task 27 核心功能驗收通過。** 該輪回覆表示 `go test -count=1 ./...`、`go vet ./...` 均通過，且額外診斷涵蓋新增輸出、儲存失敗傳遞、狀態維持及 Tab／換行空白驗證。此為歷史驗收紀錄，並非每次續課都已重新執行。
- **非阻擋測試補強，建議在 Task 28 一併補上**：`TestRunCLI_Add` 目前只驗證 `service.List()`，還應確認 `fakeRepo.todos` 與成功輸出 `buf.String()`；`TestRunCLI_Add_SaveFailed` 已檢查有 Error 和 Service 狀態不變，還可用 `errors.Is` 驗證錯誤、檢查 Repository 狀態不變以及無成功訊息。這些補強不影響 Task 27 已通過的結論。

### 當時預定的開課位置（歷史紀錄，已完成）：第 5 堂 — Task 28：stderr、Exit Code 與可測試的 main()

**不要重講 Task 01–27，也不要提前進入 `done`／`delete`／`search`、SQLite 或 REST API。** 按照「問題 → 原因 → Go 概念 → 小任務 → 學員提交 → Code Review」延續課程。

**起始思考題（已在上次課程布置、尚待學員回答）**：

> 為什麼選擇讓 `runApp()` 回傳 `int`，而不是直接在 `runApp()` 呼叫 `os.Exit()`？除了方便測試以外，還有什麼程式設計上的好處？

**已布置的 Task 28 實作要求：**

1. 新增 `func runApp(args []string, stdout io.Writer, stderr io.Writer) int`，集中處理命令解析、Repository／Service 初始化與呼叫 `runCLI()`。
2. 任何解析、初始化或執行錯誤要寫入 `stderr` 並回傳 `1`；成功時資料寫入 `stdout`、回傳 `0`。
3. `main()` 保留「沒有命令列參數就呼叫 `runInteractive()`」的現有行為；有參數時呼叫 `runApp()`，並由 `main()` 負責 `os.Exit()`。**不要在 `runApp()` 或 `runCLI()` 中直接呼叫 `os.Exit()`**，以免終止測試程序。
4. 保持目前 `runCLI(command CLICommand, service *todo.Service, w io.Writer) error` 的函式簽章；不引入第三方 CLI Framework。
5. 新增 `TestRunApp`，使用 `bytes.Buffer` 分別收集 stdout 與 stderr，至少測試：`["unknown"]` → stdout 空／stderr 有 `unknown command`／exit 1；`["add"]` → stdout 空／stderr 有 `add missing title`／exit 1；`["list", "extra"]` → stdout 空／stderr 有 `too many arguments`／exit 1。這些案例均在 Repository 初始化前失敗，不需要 Fake Repository。
6. 追加前述 Task 27 的非阻擋 Assertions，讓測試完整涵蓋 Repository 持久化、成功輸出及儲存失敗的錯誤與狀態。

**教學提示**：要說清楚 `return` 離開 `main()` 不等於非零 Exit Code；stdout 是正常資料，stderr 是錯誤／診斷輸出；`os.Exit` 屬於進程控制，應只留在程式最外層。提供函式簽章與必要提示即可，**不要直接替學員寫出整份 Task 28 解答**。等學員回答思考題並提交程式後再進行驗收。

**檔案版本提醒**：本檔已更新至 Task 28，但同一份舊 ZIP 中的 `main.go`、`main_test.go` 可能仍是 Task 25 左右的版本。再次開課時應優先使用 2026-11-10 另行上傳的最新版，若無法取得則請學員提供最新檔案，不要據舊 ZIP 否定已驗收的 Task 26–27。


---

## 最新續課更新（2026-11-12）— 第 5 堂 Task 28 已完成，下次從 Task 29 開始

> **此節為目前有效的最新續課定位；上方日期較早的 Task 28 布置與「下次開始」文字僅屬歷史紀錄，不應再當作未完成作業。**

### Task 28 已完成並驗收 — stderr、Exit Code 與可測試的 main()

- 學員實作 `func runApp(args []string, stdout io.Writer, stderr io.Writer) int`：先 `parseCommand(args)`；合法才建立 `storage.NewJSONRepository("todos.json")` 與 `todo.NewServiceWithRepository(repo)`，再呼叫 `runCLI(command, service, stdout)`；錯誤使用 `fmt.Fprintln(stderr, err)`，回傳 Exit Code `1`；成功回傳 `0`。
- `main()` 保留「沒有命令列參數時進入 `runInteractive()`」；有參數時執行 `exitCode := runApp(os.Args[1:], os.Stdout, os.Stderr)` 並 `os.Exit(exitCode)`。目前無參數分支使用 `os.Exit(0)`；老師指出改成 `return` 可讓 `main()` 中的 `defer` 正常執行，**這是可選風格改善，非 Task 28 阻擋項**。
- `main_test.go` 新增 `TestRunAppCommand_Error` table-driven tests：未知命令、`add` 缺少標題、`list` 多餘參數；分別檢查 Exit Code `1`、stdout 空、stderr 包含對應錯誤；每個 `t.Run` 內自行建立 `bytes.Buffer`，沒有共用污染問題。
- `TestRunCLI_Add` 補上成功輸出（包含 `2. Sleep 添加成功`）與 `fakeRepo.todos` 等於期望 Todo 清單的檢查；`TestRunCLI_Add_SaveFailed` 補上 `errors.Is(err, saveError)`、Fake Repository 狀態不變、沒有輸出「添加成功」；原本 Service 狀態不變的檢查保留。
- 學員正確理解：`os.Exit` 會直接終止 Process、不會執行當前或其他函式中尚未執行的 `defer`；若 `runApp()` 已經正常 `return`，則其 `defer` 已先完成。`%w` 保留錯誤鏈，而 `%v` 不保留可供 `errors.Is` 遍歷的原始錯誤關聯。
- 已討論 Unit Test、Integration Test、End-to-End Test 的測試範圍：`runApp()` 單元測試通過不代表 `main()` 一定有正確將 Exit Code 傳給作業系統；學員修正後有補上外部執行檢驗。
- 已釐清 **Service 管理應用程式中的 Todo 狀態；Repository 負責載入／儲存**。`fakeRepository` 只是記憶體中的模擬實作，檢查 `fakeRepo.todos` 不能視為真實 JSON 檔案持久化測試。`AddAndSave` 是先 Save 候選資料，成功才更新 Service 狀態；失敗時維持舊狀態。
- **歷史驗收結論**：上次對話回覆記錄：使用單獨提交的最新版 `main.go`／`main_test.go` 在暫存副本執行 `go test -count=1 ./...`、`go vet ./...`、`go build` 均通過，CLI 未知命令／缺參數／多參數的 Exit Code 為 `1`。環境為 Go 1.23.2，僅暫時調整副本的 `go.mod` 宣告，沒有修改原始專案。這是前次對話的驗收紀錄，**不是本次更新 `class.md` 時重新跑的測試**。
- **Task 28 正式結束**；不要重新要求同一組測試、不要重講 stdout／stderr、`os.Exit` 或 Fake Repository 概念，除非學員後續實作有新問題。

### 下次正式開課：第 5 堂 — Task 29：新增 `done` 子命令

**功能目標**：逐步支援 `todo done 3`，將正整數 ID 對應的 Todo 標記完成。第一步先在 `parseCommand()` 中解析與驗證 `done`，**暫時不改 `runCLI()`**；待學員完成第一階段後再處理執行／儲存流程。

**起始設計思考題（已布置，尚未得到學員回答）：**

- **A**：`parseCommand()` 以 `strconv.Atoi()` 驗證 ID，但 `CLICommand.Args` 保留字串 `"3"`，`runCLI()` 使用時再解析。
- **B**：調整 `CLICommand` 結構，使其能儲存已解析的整數 ID，避免重複轉換。
- **下次開課先請學員選擇 A 或 B，說明原因。** 先討論解析責任、型別安全、結構複雜度與擴充性，不要預先宣告唯一答案；沿用以學員自行思考實作為主的課程模式。

**Task 29 第一階段預期案例**（不在續課檔直接提供完成程式碼）：

| 輸入參數 | 預期行為 |
| --- | --- |
| `[]string{"done", "3"}` | 合法命令，成功解析 |
| `[]string{"done"}` | 缺少 ID 錯誤 |
| `[]string{"done", "abc"}` | 非數字 ID 錯誤 |
| `[]string{"done", "0"}` | 非法 ID 錯誤 |
| `[]string{"done", "-1"}` | 非法 ID 錯誤 |
| `[]string{"done", "3", "extra"}` | 參數過多錯誤 |

**教學原則**：先確認最新源碼，延續「問題 → 原因 → Go 概念 → 小任務 → 學員提交 → Code Review → 驗收」；暫不進入 `delete`、`search`、SQLite、REST API 或完整多命令改寫。沿用標準函式庫與既有專案結構。

**文件與程式碼版本提醒**：本次僅更新 `class.md` 的教學進度。歷史 ZIP 不是已更新到 Task 28 的整合版；下次若需要完整專案，請優先使用學員當次提供的最新檔案，避免拿舊 ZIP 快照當成最新成果。
