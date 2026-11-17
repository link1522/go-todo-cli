# Go Todo 專案實作課程｜教學規則與續課紀錄

> 最後整理：2026-11-17  
> **目前進度：第 5 堂「CLI Application」進行中；Task 01–29 已完成並通過核心驗收。**  
> **下次起點：第 5 堂 — Task 30「新增 `delete` 子命令」，從 Task 30.1「Extract Function Refactoring」開始。**

本文件是跨對話續課的**唯一有效進度索引**。請以「下次上課」段落為最新指引；下方「已完成任務」是學習履歷，不是待辦清單。除非當次程式碼出現新的問題，**不得要求重新完成 Task 01–29**。

---

## 一、課程目標與使用方式

使用者以自己的 `todo-cli` 專案學習 Go：從 CLI、測試、Service 與 Repository，循序演進至 SQLite、REST API、Concurrency 和工程化的小型 Go Backend。目標不只是學語法，而是理解實際專案中**資料流、設計取捨、錯誤契約、測試與 Go 慣例**。

每次開課：

1. 先閱讀本文件的「下次上課」及「專案現況」，再檢查**本次提供的最新原始碼**，不得依過時 ZIP 猜測現況。
2. 若單獨上傳的程式檔比 ZIP 新，以**最新提交的檔案**為準；也可使用已整理的整合版 ZIP 作為基準，但仍須核對新檔案。
3. 每次只教一個合理大小的階段，按照「**問題 → 原因 → Go 概念 → 設計討論 → 學員實作 → Code Review → 驗收**」推進。
4. 先說明為什麼需要修改，再給任務、提示、函式簽章或部分 Skeleton；**不要直接交付整個任務的完整答案**。
5. 每階段完成後簡述學到了什麼、專案變化與下一步。已回答的思考題不要重問。

## 二、下次上課：Task 30（唯一有效的續課入口）

### Task 30 — 新增 `delete` CLI 子命令

目標逐步支援：

```bash
todo delete 1
```

目前 `runInteractive()` 雖已有互動式 `delete`，但 **CLI 子命令** 的 `parseCommand()`／`runCLI()` 尚未支援 `delete`。不得將兩者混為一談。

### Task 30.1 — Extract Function Refactoring（先完成這一階段）

**問題**：目前 `done` 的 `parseCommand()` 分支使用 `strconv.Atoi(args[1])`，再檢查 `id > 0`。新 `delete` 命令有相同的 ID 解析規則；直接複製會造成重複驗證邏輯。

**請交付給學員的小任務**：

1. 新增 `func parsePositiveID(raw string) (int, error)`，負責字串轉 `int`、驗證 ID 大於零，以及適當回傳錯誤。
2. 讓既有 `done` 和新的 `delete` 分支共用 `parsePositiveID()`，但維持 `list`、`add`、`done` 原有行為不變。
3. 暫時沿用 `type CLICommand struct { Name string; Args []string; ID int }`，**不要此時重構為每個命令一個 Struct**。
4. 在 `main_test.go` 擴充成功與失敗的 table-driven tests；**本階段只修改解析及測試，不要修改 `runCLI()`**。

**驗收案例**：

| 傳入 `parseCommand()` 的參數       | 預期                              |
| ---------------------------------- | --------------------------------- |
| `[]string{"delete", "3"}`          | `CLICommand{Name:"delete", ID:3}` |
| `[]string{"delete"}`               | `delete missing ID`               |
| `[]string{"delete", "abc"}`        | ID 必須是整數                     |
| `[]string{"delete", "0"}`          | ID 必須大於 0                     |
| `[]string{"delete", "-1"}`         | ID 必須大於 0                     |
| `[]string{"delete", "3", "extra"}` | 參數過多                          |
| `[]string{"done", "3"}`            | 既有行為維持不變                  |

**已布置、尚未回答的思考題**：`parsePositiveID()` 應只負責「字串轉整數與正整數驗證」，還是連「缺少參數、參數過多」也一起處理？請讓學員從 **Single Responsibility Principle** 與函式抽象化邊界提出理由，不要代替回答。

**提交與驗收**：請學員實作後提交最新 `main.go`、`main_test.go` 和思考題答案。Review 檢查行為、錯誤訊息、回歸測試及函式職責；Task 30.1 完成後，再循序設計 `DeleteAndSave()`、`runCLI()` 中的 `delete` 分支與儲存失敗一致性測試。不要提前完成後續階段，也不要跳往 SQLite／REST API。

---

## 三、目前專案現況（Task 29 最終提交）

### 檔案與模組

- Module：`todo-cli`。歷史整合版 `go.mod` 宣告 `go 1.25.5`；請依最新實際檔案核對。
- `main.go`、`main_test.go`：命令解析、CLI 執行、錯誤輸出與測試。
- `todo/todo.go`、`todo/todo_test.go`：Todo Domain、資料操作、自訂錯誤。
- `todo/service.go`、`todo/service_test.go`：Service、Repository 注入、狀態一致性。
- `todo/repository.go`、`todo/stats.go`、`todo/stats_test.go`：介面與統計服務。
- `storage/JSONRepository.go`、對應測試及既有 `storage.go`：JSON 持久化。

### CLI 實作狀態

| 命令          | CLI 子命令   | 備註                                 |
| ------------- | ------------ | ------------------------------------ |
| `list`        | 已完成       | 經 `printTodos(io.Writer, ...)` 輸出 |
| `add "Title"` | 已完成       | `Service.AddAndSave()`               |
| `done 1`      | 已完成       | `Service.CompleteAndSave()`          |
| `delete 1`    | **尚未完成** | Task 30，先處理參數解析              |
| `search "Go"` | 尚未完成     | 後續 CLI 任務                        |

核心型別與函式：

```go
type CLICommand struct {
    Name string
    Args []string
    ID   int
}

func parseCommand(args []string) (CLICommand, error)
func runCLI(command CLICommand, service *todo.Service, w io.Writer) error
func runApp(args []string, stdout io.Writer, stderr io.Writer) int
```

- `parseCommand()` 先驗證命令與參數，無效命令不應啟動 Repository Load。
- `runApp()` 建立 `storage.NewJSONRepository("todos.json")`、`todo.NewServiceWithRepository(repo)`，呼叫 `runCLI()`；錯誤寫入 `stderr` 並回傳 Exit Code `1`，成功回傳 `0`。
- `main()` 有命令列參數時執行 `runApp(os.Args[1:], os.Stdout, os.Stderr)` 與 `os.Exit(exitCode)`；無參數時保留原有 `runInteractive()`。
- `runCLI()` 只執行 Service 動作、透過注入的 `io.Writer` 輸出成功訊息並回傳錯誤；不負責建立 Repository、直接寫 `os.Stderr` 或呼叫 `os.Exit()`。
- `done` 解析已在 `parseCommand()` 將字串 ID 轉成正整數，儲存在 `CLICommand.ID`；**尚未抽出 `parsePositiveID()`**。
- `runInteractive()` 仍是舊的 Scanner／Domain／storage 路徑，與逐步開發中的 CLI 子命令路徑不同；除非當前任務需要，不要順手大改。

### 既有 Service／Repository 契約

```go
type TodoRepository interface {
    Load() ([]Todo, error)
    Save(todos []Todo) error
}
```

- `Service` 包含 `todos []Todo` 與 `repo TodoRepository`；具有 `NewService`、`NewServiceWithRepository`、`Add`、`List`、`Complete`、`Delete`、`Search`、`Save`、`AddAndSave`、`CompleteAndSave`。
- `AddAndSave()`／`CompleteAndSave()` 使用「**複製候選狀態 → 修改候選 → Save 複製資料 → Save 成功後提交到 Service**」流程；存檔失敗時不更新 `s.todos`，並以 `%w` 保留錯誤鏈。
- `CompleteAndSave()` 找不到 ID 時回傳 `&TodoNotFoundError{ID:id}`；其 `Unwrap()` 連到 `ErrTodoNotFound`，讓 `errors.Is()` 與 `errors.As()` 都能使用。
- Service 在 Load、List、Save 和變更候選值時利用 `slices.Clone()` 隔離 Slice 底層陣列。**`slices.Clone` 是淺複製**：目前 `Todo` 只有 `int`、`string`、`bool` 等值型別欄位，足以滿足目前的元素隔離需求；未來如加入 Map、Slice、Pointer 欄位須重新評估。
- `DeleteTodo()` 仍允許修改輸入 Slice 的底層陣列；需要保留原始資料的呼叫端自行 Clone。`SearchTodos("")` 的既定規格是回傳全部 Todo。
- `fakeRepository.Load()` 直接回傳 `r.todos`，`Save()` 成功時直接保存傳入 Slice；這是為了刻意檢查 Service 的 Ownership 防護，不必擅自改成 Clone。Fake 的記憶體狀態測試**不是**真正的磁碟持久化測試。

### 最近一次驗收（歷史紀錄，非本次重新執行）

- Task 29.3 的 `main.go`、`main_test.go` 與 Task 29.2 的 `service.go`、`service_test.go` 曾整合驗證，當次教學紀錄為 `go test -count=1 ./...`、`go vet ./...` 通過。
- 曾在暫存目錄執行 `add` → `done` → `list`，觀察 JSON 儲存、`[x]` 輸出，以及找不到 ID 時的 stderr／Exit Code。
- 過往執行環境為 Go 1.23.2，因此**僅在測試暫存副本**調整 `go.mod` 的版本宣告；不得誤認為原專案已降級，也不可把歷史測試結果當成每次新提交都已驗證。
- **小型非阻擋改善**：`TestRunCLI_Done_NotFound` 可補「stdout 應為空」的 assertion。Task 29 已通過，不得因此阻止 Task 30。

---

## 四、已完成任務與學習成果

> 以下是去重後的知識與驗收紀錄。相同概念僅保留首次學習或關鍵深化的部分，避免把已完成的歷史任務誤當成新作業。

### 第 1 堂｜Testing、Slice 與 Todo Domain（初期 Task 01–08）【已完成】

- **AddTodo／Testing**：`testing`、`t.Run`、table-driven tests；驗證新增 Todo 的長度、ID、標題、預設 `Done=false` 與原資料不變。
- **Slice 基礎**：Slice Header、`len`／`cap`、`append`、Backing Array，共享資料與 Mutation；理解 `backup := original` 只複製 Header，`slices.Clone` 會建立另一份元素儲存空間。
- **DeleteTodo**：測試刪第一、中、末、找不到 ID、空 Slice；以 `slices.Equal` 檢查結果；建立 `TestDeleteTodo_MutatesInput` 與 Clone 後原資料不變的測試，修正 `found` 與 Mutation 驗證混在一起可能產生的假陽性。保留「`DeleteTodo` 可能修改原輸入 Backing Array」的契約。
- **CompleteTodo**：測試存在、不存在與重複完成；理解 `range` 複製 Struct 元素，修改 `item.Done` 不會修改 `todos[i].Done`。
- **GetNextID**：測試空、連續、不連續、亂序 ID；理解 `maxID + 1` 在刪除、資料庫歷史關聯或並行情境可能面臨 ID 管理限制。
- **SearchTodos**：完整／部分匹配、大小寫不敏感、找不到、空字串與順序；定義空關鍵字回傳全部，並測試搜尋不修改輸入。當前 `Todo` 值型別欄位可避免元素引用共享，未來有巢狀參考型欄位則須重新檢查。

### 第 2 堂｜Error Handling（本堂 Task 01–09）【已完成】

- `storage.LoadTodos() ([]todo.Todo, error)`：檔案不存在（`errors.Is(err, os.ErrNotExist)`）視為首次使用；其他 I/O／JSON 解析錯誤往上回傳，避免以空資料覆寫損壞檔。
- `fmt.Errorf("...: %w", err)` 的 `%w` 保留錯誤鏈；`%v` 只有字串內容。練習 `errors.Is`、`errors.As`、`Unwrap()`。
- 理解 `var pathErr *os.PathError` 與 `&pathErr` 的 `*T`／`**T` 關係，以及 `errors.As()` 如何填入目標變數。
- 建立 `TodoNotFoundError{ID int}`（Pointer Receiver `Error()`）與 Sentinel `ErrTodoNotFound`，透過 `Unwrap()` 連接。`DeleteTodo()` 改回傳 `([]Todo, error)`；測試錯誤語意、具體型別、ID 與包裝後的行為。
- 學會在 `errors.As()` 失敗、後續需解參照指標時使用 `t.Fatal`／`t.Fatalf` 而非 `t.Errorf`。

### 第 3 堂｜Struct、Method 與 Service（Task 10–15）【已完成】

- **Task 10**：`(*Todo).MarkDone()`，以 Pointer Receiver 更新原始物件；`CompleteTodo` 使用 `todos[i].MarkDone()`。
- **Task 11**：`NewTodo(id, title)` Constructor 與 Go Zero Value；理解 Constructor 是設計選項而非強制。
- **Task 12**：`(Todo).Matches(keyword)` Value Receiver，使用 `strings.ToLower`／`strings.Contains`；保留空關鍵字比對全部，討論 Method Set。
- **Task 13**：`Service{todos []Todo}`、`NewService`、`Add`、`List`；在建構與回傳時 Clone，保障呼叫端不能修改 Service 內部 Slice。
- **Task 14**：`Service.Complete()` 與自訂錯誤傳遞；測試存在、不存在、已完成等情境，以及 `errors.Is`／`errors.As`。
- **Task 15**：`Service.Delete()`、`Service.Search()`，沿用 Domain 邏輯、錯誤契約與搜尋行為；已有成功／失敗與狀態不變測試。早期曾在暫存測試中檢查額外搜尋 Ownership 情境，是否寫入正式測試須以最新檔案為準。

### 第 4 堂｜Interface、Dependency Injection 與資料一致性（Task 16–21）【已完成】

- **Task 16 — `TodoReader`**：定義 `List() []Todo` 小介面與 `CountPending()`；以真實 Service、Fake Reader、空集合測試。理解 Go 隱式實作，及 `List` 是 Pointer Receiver 時 `*Service` 符合介面、`Service` 值不符合的 Method Set 差異。
- **Task 17 — `StatsService`**：Constructor Injection 與 `PendingCount()`；介面降低依賴耦合，避免為假想需求過度抽象。
- **Task 18 — `TodoRepository`**：定義 `Load() ([]Todo,error)`、`Save([]Todo) error`；釐清 Package Function 與型別 Method 的差異。
- **Task 19 — `storage.JSONRepository`**：Constructor 注入檔案路徑、`os.ReadFile`／`os.WriteFile` 與 JSON Marshal／Unmarshal；`t.TempDir()` 測試檔案隔離、合法資料、缺檔及錯誤。可選編譯期介面斷言 `var _ todo.TodoRepository = (*JSONRepository)(nil)`。
- **Task 20 — Repository 注入 Service**：`NewServiceWithRepository()`、`Save()`、Fake Repository；處理 `nil`、Load Error、Save Error，並在 Load／Save 邊界 Clone 防止共用 Backing Array。認知「Save 過」不代表記憶體自動隔離。
- **Task 21 — `AddAndSave()`**：改成候選資料先持久化、成功才提交的模式；儲存失敗時保留原始狀態，錯誤使用 `%w`。學員正確解釋 Slice Ownership、`candidate` 與 Repository 傳入 Clone 的差異。早期版本曾在最末提交時額外 Clone，屬安全但冗餘；**目前最新程式碼是 `s.todos = candidate`**。這套模式不等於 Mutex 並發安全、檔案原子寫入或 SQL Transaction。

### 第 5 堂｜CLI Application（Task 22–29）【目前已完成部分】

- **Task 22–24 — CLI 與輸出測試**：`os.Args`、`list` 子命令、空參數防護；`printTodos(w io.Writer, ...)` 搭配 `fmt.Fprintf`，在測試使用 `bytes.Buffer`。理解 `*bytes.Buffer` 實作 `io.Writer`，並用 table-driven tests 比對空清單、未完成、已完成與多筆資料的精確輸出。
- **Task 25 — Fake Repository**：使用 Fake 提供可預期的 Todo 資料，避免 CLI 單元測試依賴真實 JSON 檔案。
- **Task 26 — 解析與執行分離**：`parseCommand()` 先驗證，`runCLI()` 改為接收 `*todo.Service` 而非 Repository；依賴建立置於外層。此設計是學員自行選擇，當前沒有必要額外建立專屬 CLI Interface。
- **Task 27 — `add` 子命令**：`CLICommand{Name, Args}`，驗證非空白標題、精確參數數量，使用 `Service.AddAndSave`；測試成功輸出、Service 與 Fake Repository 狀態、Save 失敗錯誤傳遞及無成功訊息。
- **Task 28 — `runApp()`、stderr／Exit Code**：合法命令才初始化依賴；`runApp()` 集中處理錯誤輸出並回傳狀態碼，`main()` 才能 `os.Exit`。學員理解 `os.Exit` 不執行仍待執行的 `defer`，且 `return`／`os.Exit(0)` 有不同生命週期語意；測試 `unknown`、缺標題、多餘參數的 Exit Code、stdout 與 stderr。區分 Unit／Integration／E2E，Fake Repository 測試不代表磁碟測試。
- **Task 29.1 — `done` 解析**：選擇在 `parseCommand()` 轉成 `int`，避免多處重複解析；將 `ID int` 加入 `CLICommand`，測試正整數、缺參數、非整數、零、負數、多餘參數。學員提出「每個 Command 個別 Struct」設計；已討論 Interface／Type Switch，但現階段不重構。即便 `DoneCommand`、`DeleteCommand` 欄位相同，學員主張不同業務語意應保持分離，已獲認可。
- **Task 29.2 — `CompleteAndSave()`**：使用 Clone Before Mutation、Save Before Commit；找不到時回傳可供 `errors.Is`／`errors.As` 使用的 `TodoNotFoundError`，Save 失敗用 `%w` 包裝並保持狀態不變。涵蓋成功、ID 不存在、已完成、無 Repository、Save 失敗、Ownership 測試。學員已理解 Repository 不可任意修改 Service 候選狀態。
- **Task 29.3 — `done` 執行**：`runCLI()` 呼叫 `service.CompleteAndSave(command.ID)`，成功才寫 `"%d 已完成\n"`；錯誤交給 `runApp()`，使職責分離且容易用 `bytes.Buffer` 測試。成功／NotFound／SaveFailed 測試已涵蓋 Service、Fake Repository 狀態、錯誤型別、`t.Fatalf` 防止 nil 解參照；已完成整合與核心驗收。

---

## 五、教學與 Code Review 約定

1. **真實專案導向**：每堂課直接延續當前 Todo 程式碼；不另做無關的玩具範例，也不為了教概念而硬塞設計。
2. **循序漸進**：一次一個小階段、先讓學員嘗試；思考題先收集理由，再討論不同設計的取捨，不預設唯一答案。
3. **優先標準函式庫**：CLI 以 `os.Args`／`flag` 為主；日後 API 先用 `net/http`，資料庫先用 `database/sql`；不急著引入 Cobra、Gin、Echo、GORM 等 Framework。
4. **Go 慣例**：若程式可執行但不夠 idiomatic，指出具體差異及原因，不只說「比較好」。
5. **Review 重點**：Correctness、Error Handling、Ownership／Mutation、Pointer／Method Set、Package Responsibility、命名、Testability、Go Idioms、錯誤契約與測試是否存在 False Positive。
6. **測試習慣**：Table-driven tests、`t.Run`、`errors.Is`／`errors.As`、必要時 `t.Fatalf`、狀態與持久化結果分開斷言；確認 Save 失敗時不出現成功訊息。
7. **避免過度推進**：Concurrency、原子寫檔、資料庫 Transaction、複雜 CLI Framework 等主題留到專案自然需要時再深入。
8. **每次驗收依據最新檔案**：歷史「測試通過」只是紀錄；如要宣稱本次測試通過，需實際執行。若 Go toolchain 不符，明確揭露暫存副本的處理方式。

## 六、後續長期學習路線

| 階段 | 主題                        | 重點                                                                                                                        | 狀態                       |
| ---- | --------------------------- | --------------------------------------------------------------------------------------------------------------------------- | -------------------------- |
| 1    | Testing 與 Slice            | `testing`、table-driven tests、`len/cap`、Backing Array、Mutation                                                           | 完成                       |
| 2    | Error Handling              | `%w`、`errors.Is/As`、Sentinel／Custom Error、I/O 錯誤                                                                      | 完成                       |
| 3    | Struct 與 Method            | Pointer／Value Receiver、Constructor、Service                                                                               | 完成                       |
| 4    | Interface 與 DI             | Small Interface、Repository Pattern、JSONRepository                                                                         | 完成                       |
| 5    | CLI Application             | `os.Args`、`flag/FlagSet`、stdout／stderr、Exit Code、`add/list/done/delete/search`                                         | **進行中，下一步 Task 30** |
| 6    | 擴充 Todo Domain            | `time.Time`、Custom Type、`const/iota`、驗證、排序、Priority、Due Date、Created Time、Filter                                | 未開始                     |
| 7    | SQLite                      | `database/sql`、Query／QueryRow／Exec／Scan、SQLiteRepository                                                               | 未開始                     |
| 8    | SQL Transaction 與 Relation | Migration、Foreign Key、Transaction、1:N／N:M、Tags                                                                         | 未開始                     |
| 9    | REST API                    | `net/http`、Handler／ServeMux、JSON、HTTP Method／Status Code                                                               | 未開始                     |
| 10   | HTTP Testing 與 Middleware  | `httptest`、Logging、Request ID、Recovery、API Error Handling                                                               | 未開始                     |
| 11   | Goroutine／Channel          | Goroutine、Channel、Select、WaitGroup、Mutex、Context、Cancellation、Graceful Shutdown                                      | 未開始                     |
| 12   | 工程化                      | `cmd/`、`internal/`、Package Design、fmt／vet／test／race、Config、Structured Logging、Makefile、Dockerfile、GitHub Actions | 未開始                     |

---

## 七、跨次續課的重要提醒

- **唯一待辦起點**：Task 30.1；不要再使用舊版紀錄中「下次從 Task 26／28／29 開始」的歷史規劃。
- **最新檔案優先**：專案 ZIP 可能有版本落差；若新對話無法取得最新程式碼，請學員上傳後再做具體 Code Review，不得捏造已檢視或已執行的結果。
- **歷史決策不要重問**：Task 29 選方案 B、個別 Command Struct 的討論、Slice Clone 原因、Service／Repository 職責、runApp 負責錯誤輸出、`io.Writer` 的測試優點等均已學過。
- **非阻擋改善不得變成關卡**：例如 NotFound 時 stdout 空的補充測試，以及其他先前已標示非必要的重構建議。
- **未開始的後續命令**：`delete` CLI 與 `search` CLI；`runInteractive()` 的舊功能不代表 CLI 已完成。
