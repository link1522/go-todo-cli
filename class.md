# Go Todo 專案實作課程｜教學規則與續課紀錄

> 最後整理：2026-11-20  
> **目前進度：第 5 堂「CLI Application」進行中；Task 01–31 已完成，Task 31.1、31.2、31.3（含 31.3.1）均已驗收通過。**  
> **下次起點：第 5 堂 — Task 32.1「CLI Help Command」，先回答 help 與 Repository 初始化位置的架構思考題，再逐步實作。**

本文件是跨對話續課的**唯一有效進度索引**。請以「下次上課」段落為最新指引；下方「已完成任務」是學習履歷，不是待辦清單。除非當次程式碼出現新的問題，**不得要求重新完成 Task 01–30**。

---

## 一、課程目標與使用方式

使用者以自己的 `todo-cli` 專案學習 Go：從 CLI、測試、Service 與 Repository，循序演進至 SQLite、REST API、Concurrency 和工程化的小型 Go Backend。目標不只是學語法，而是理解實際專案中**資料流、設計取捨、錯誤契約、測試與 Go 慣例**。

每次開課：

1. 先閱讀本文件的「下次上課」及「專案現況」，再檢查**本次提供的最新原始碼**，不得依過時 ZIP 猜測現況。
2. 若單獨上傳的程式檔比 ZIP 新，以**最新提交的檔案**為準；也可使用已整理的整合版 ZIP 作為基準，但仍須核對新檔案。
3. 每次只教一個合理大小的階段，按照「**問題 → 原因 → Go 概念 → 設計討論 → 學員實作 → Code Review → 驗收**」推進。
4. 先說明為什麼需要修改，再給任務、提示、函式簽章或部分 Skeleton；**不要直接交付整個任務的完整答案**。
5. 每階段完成後簡述學到了什麼、專案變化與下一步。已回答的思考題不要重問。

## 二、下次上課：Task 32.1 — CLI Help Command（唯一有效的續課入口）

### 背景與目標

目前 `list`、`add`、`done`、`delete`、`search` 的 CLI 子命令都已完成；執行 `todo --help` 仍會得到 `unknown command`。下一步為 CLI 加入簡明的使用說明，讓使用者知道可執行哪些命令。

**目標命令**：

```bash
todo --help
```

**預期 stdout（格式視為本階段的起始規格）**：

```text
Usage: todo <command> [arguments]

Commands:
  list               List all todos
  add <title>        Add a new todo
  done <id>          Mark a todo as done
  delete <id>        Delete a todo
  search <keyword>   Search todos
```

### 當前架構與需要討論的取捨

目前命令執行路徑是：

```text
main()
  -> runApp(args, stdout, stderr)
  -> runAppWithPath(args, stdout, stderr, path)
       -> parseCommand(args)
       -> storage.NewJSONRepository(path)
       -> todo.NewServiceWithRepository(repo)  // 載入資料
       -> runCLI(command, service, stdout)
```

`runApp()` 只是提供預設路徑 `"todos.json"` 的 Wrapper，`runAppWithPath()` 才是正式與整合測試共用的執行流程。`--help` 是不需要 Todo 資料的說明性操作，因此應先討論如何避免為 help 載入 JSON。**目前尚未實作 help，也尚未決定分支應放在哪一層；不要先假定已完成。**

### 下次開場的思考題（尚未回答）

> 如果 `todos.json` 的格式損壞，執行 `todo --help` 是否應該因此失敗？如果不應該，你認為應該在現有架構的哪一個位置處理 `--help`，才能避免不必要的 Repository 初始化？

先讓學員回答，再評估其理由；不要提前提供完整實作。

### 思考題之後的預定任務與驗收方向

1. 讓命令解析可辨識 `--help`，不破壞現有 `list`／`add`／`done`／`delete`／`search` 的解析契約。是否同時支援 `help`、額外參數如何處理，應先明確決定規格，**不能擅自擴充**。
2. 在合適的邊界輸出上述說明到可注入的 `io.Writer`，成功返回 Exit Code `0`、stderr 為空；help 不應依賴資料檔案正常與否。
3. 新增 table-driven 或整合測試，至少涵蓋 help 的完整輸出、exit code、stderr，及 JSON 遺失／損壞時仍可顯示 help。
4. 不順手重構 `runInteractive()`，不引入 Cobra 等 CLI Framework；維持小階段「討論 → 學員實作 → Code Review → 驗收」。

**提交方式**：下次先回答思考題；若已完成實作，再附最新 `main.go`、`main_test.go`。通過 Task 32.1 後才規劃下一階段，不重做 Task 31。

---

## 三、目前專案現況（Task 31.3 最新提交）

### 檔案與模組

- Module：`todo-cli`。原專案 `go.mod` 宣告 `go 1.25.5`；先前測試環境為 Go 1.23.2，僅在**暫存副本**調整版本以執行驗證，未修改原始 `go.mod`。
- `main.go`、`main_test.go`：CLI 命令解析、輸出、路徑注入及整合測試。**2026-11-20 最新獨立提交版本優先於 11/18 原始 ZIP。**
- `todo/todo.go`、`todo/todo_test.go`：Todo Domain、資料操作、自訂錯誤。
- `todo/service.go`、`todo/service_test.go`：Service、Repository 注入、狀態一致性。
- `todo/repository.go`、`todo/stats.go`、`todo/stats_test.go`：介面與統計服務。
- `storage/JSONRepository.go`、對應測試及 `storage.go`：JSON 持久化。

### CLI 實作狀態

| 命令 | CLI 子命令狀態 | 備註 |
| --- | --- | --- |
| `list` | 已完成 | `printTodos(io.Writer, ...)` 輸出 |
| `add "Title"` | 已完成 | `Service.AddAndSave()` |
| `done 1` | 已完成 | `Service.CompleteAndSave()` |
| `delete 1` | 已完成 | `Service.DeleteAndSave()` |
| `search "Go"` | **已完成** | Task 31.1–31.3：解析、`Service.Search()`、`printTodos()`、JSON 整合測試 |
| `--help` | **尚未開始** | 下一步 Task 32.1；當前會回傳 `unknown command` |

核心型別與函式（以最新 `main.go` 為準）：

```go
type CLICommand struct {
    Name string
    Args []string
    ID   int
}

func parseCommand(args []string) (CLICommand, error)
func runCLI(command CLICommand, service *todo.Service, w io.Writer) error
func runApp(args []string, stdout io.Writer, stderr io.Writer) int
func runAppWithPath(args []string, stdout io.Writer, stderr io.Writer, path string) int
```

- `runApp()` 保持原簽章，轉呼叫 `runAppWithPath(args, stdout, stderr, "todos.json")`；正式環境與整合測試走同一段邏輯，以不同路徑隔離資料。
- `runAppWithPath()` 先 `parseCommand()`，再建立 `storage.NewJSONRepository(path)`、`todo.NewServiceWithRepository(repo)`、執行 `runCLI()`。解析、讀檔或執行錯誤寫入 stderr，回傳 Exit Code `1`；成功回傳 `0`。
- `main()` 有參數時呼叫 `runApp(os.Args[1:], os.Stdout, os.Stderr)` 並 `os.Exit(exitCode)`；無參數時仍執行原有 `runInteractive()`。
- `runCLI()` 不建立 Repository、不呼叫 `os.Exit()`，只執行 Service 動作並向注入的 `io.Writer` 寫入結果。
- `parseCommand()` 的 `search` 只接受**恰好一個**關鍵字參數；使用 `strings.TrimSpace()` 驗證非空白，但成功時保留原始關鍵字。缺少、多餘、空白參數分別回報 `search missing keyword`、`too many arguments`、`search empty keyword`。
- `runCLI()` 的 `search` 分支呼叫 `service.Search(command.Args[0])`；有結果時共用 `printTodos()`，無結果輸出 `"找不到任何項目\n"`；不呼叫 `Save()`，因為搜尋是唯讀操作。已移除多餘的「搜尋成功」與額外空白行。
- `done`／`delete` 仍共用 `parsePositiveID()` 的轉換與正整數驗證，並在 Service 中執行相應的 *AndSave 方法。
- `runInteractive()` 仍採 Scanner／Domain／storage 的舊路徑。**互動模式與 `os.Args` 子命令是不同流程**；除非任務需要，不要順手改動。

### 最新測試涵蓋範圍

- `TestParseCommand`／`TestParseCommand_Failed`：搜尋合法與非法引數、保留關鍵字前後空白等。
- `TestRunCLI_Search`：多筆符合、無符合、精確 stdout 比對、Service 與 Fake Repository 的資料內容不變。
- `TestRunApp_Search`：`t.TempDir()` 建立獨立 JSON 檔案，透過 `runAppWithPath()` 測試真實 JSON Repository、stdout、stderr、Exit Code 及搜尋後資料相同。
- `TestRunApp_Search_Failed`：缺少、空字串、全空白、過多關鍵字引數，驗證 Exit Code `1`、stdout 空、精確 stderr、資料內容不變。
- **可選小改善**：以 `bytes.Equal(before, after)` 比較搜尋前後原始檔案，從「內容不變」加強到「檔案位元組不變」；此前列為非阻擋，不得要求重做 Task 31。

### 既有 Service／Repository 契約

```go
type TodoRepository interface {
    Load() ([]Todo, error)
    Save(todos []Todo) error
}
```

- `Service` 包含 `todos []Todo` 與 `repo TodoRepository`；具有 `NewService`、`NewServiceWithRepository`、`Add`、`List`、`Complete`、`Delete`、`Search`、`Save`、`AddAndSave`、`CompleteAndSave`、`DeleteAndSave`。
- `AddAndSave()`／`CompleteAndSave()`／`DeleteAndSave()` 使用「**複製候選狀態 → 修改候選 → Save 複製資料 → Save 成功後提交到 Service**」流程；存檔失敗時不更新 `s.todos`，並以 `%w` 保留錯誤鏈。
- `CompleteAndSave()`／`DeleteAndSave()` 找不到 ID 時保留 `TodoNotFoundError{ID:id}` 的錯誤鏈；其 `Unwrap()` 連到 `ErrTodoNotFound`，讓 `errors.Is()` 與 `errors.As()` 都能使用。
- Service 在 Load、List、Save 和變更候選值時利用 `slices.Clone()` 隔離 Slice 底層陣列。**`slices.Clone` 是淺複製**：目前 `Todo` 只有 `int`、`string`、`bool` 等值型別欄位，足以滿足目前的元素隔離需求；未來如加入 Map、Slice、Pointer 欄位須重新評估。
- `DeleteTodo()` 仍允許修改輸入 Slice 的底層陣列；需要保留原始資料的呼叫端自行 Clone。`SearchTodos("")` 的既定規格是回傳全部 Todo。
- `fakeRepository.Load()` 直接回傳 `r.todos`，`Save()` 成功時直接保存傳入 Slice；這是為了刻意檢查 Service 的 Ownership 防護，不必擅自改成 Clone。Fake 的記憶體狀態測試**不是**真正的磁碟持久化測試。

### 最近一次驗收（歷史紀錄；新提交需重新驗證）

- **2026-11-20，Task 31.3.1 + 搜尋整合測試通過**：最新 `main.go`、`main_test.go` 審查確認 `runAppWithPath()`、`t.TempDir()`、JSON 載入搜尋、stdout／stderr／Exit Code、查詢前後資料內容比對；先前回覆記錄 `go test -count=1 ./...`、`go vet ./...`、`gofmt` 均通過。這是**當時的驗收紀錄**，不是本次修改課程檔時重新測試的結果。
- **Task 31.2** 初稿自行格式化搜尋結果，包含多餘 `"搜尋成功"` 與空白行；因使用 `strings.Contains()` 測試產生 False Positive。學員改用 `printTodos()` 與完整字串比對，修正後驗收通過。
- **Task 31.1** 搜尋引數解析及 table-driven tests 通過；學員已正確理解 Shell 引號、CLI 參數契約，及未來 flags 不應被 `strings.Join()` 混入 keyword。
- **Task 31.3 思考題**：學員指出共用真實 `todos.json` 會導致初始狀態不確定、平行測試相互干擾。**Task 31.3.1 思考題**：學員指出 Wrapper 能共用正式／測試的核心邏輯；已討論向後相容性。
- **Task 31.2 CQS 思考題**：學員原先說 `search` 與持久層無關；已補充唯讀操作不需要 `Save()`，但未來可透過 Repository 查詢資料庫。**不要重問。**
- Task 30.1–30.3 的解析、`DeleteAndSave()`、`delete` CLI 均已驗收通過。之前測試環境 Go 1.23.2，僅在暫存副本調整 `go.mod` 的版本宣告；不要把此視為原專案已改版。

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

### 第 5 堂｜CLI Application（Task 22–31）【已完成項目；Task 32 待開始】

- **Task 22–24 — CLI 與輸出測試**：`os.Args`、`list` 子命令、空參數防護；`printTodos(w io.Writer, ...)` 搭配 `fmt.Fprintf`，在測試使用 `bytes.Buffer`。理解 `*bytes.Buffer` 實作 `io.Writer`，並用 table-driven tests 比對空清單、未完成、已完成與多筆資料的精確輸出。
- **Task 25 — Fake Repository**：使用 Fake 提供可預期的 Todo 資料，避免 CLI 單元測試依賴真實 JSON 檔案。
- **Task 26 — 解析與執行分離**：`parseCommand()` 先驗證，`runCLI()` 改為接收 `*todo.Service` 而非 Repository；依賴建立置於外層。此設計是學員自行選擇，當前沒有必要額外建立專屬 CLI Interface。
- **Task 27 — `add` 子命令**：`CLICommand{Name, Args}`，驗證非空白標題、精確參數數量，使用 `Service.AddAndSave`；測試成功輸出、Service 與 Fake Repository 狀態、Save 失敗錯誤傳遞及無成功訊息。
- **Task 28 — `runApp()`、stderr／Exit Code**：合法命令才初始化依賴；`runApp()` 集中處理錯誤輸出並回傳狀態碼，`main()` 才能 `os.Exit`。學員理解 `os.Exit` 不執行仍待執行的 `defer`，且 `return`／`os.Exit(0)` 有不同生命週期語意；測試 `unknown`、缺標題、多餘參數的 Exit Code、stdout 與 stderr。區分 Unit／Integration／E2E，Fake Repository 測試不代表磁碟測試。
- **Task 29.1 — `done` 解析**：選擇在 `parseCommand()` 轉成 `int`，避免多處重複解析；將 `ID int` 加入 `CLICommand`，測試正整數、缺參數、非整數、零、負數、多餘參數。學員提出「每個 Command 個別 Struct」設計；已討論 Interface／Type Switch，但現階段不重構。即便 `DoneCommand`、`DeleteCommand` 欄位相同，學員主張不同業務語意應保持分離，已獲認可。
- **Task 29.2 — `CompleteAndSave()`**：使用 Clone Before Mutation、Save Before Commit；找不到時回傳可供 `errors.Is`／`errors.As` 使用的 `TodoNotFoundError`，Save 失敗用 `%w` 包裝並保持狀態不變。涵蓋成功、ID 不存在、已完成、無 Repository、Save 失敗、Ownership 測試。學員已理解 Repository 不可任意修改 Service 候選狀態。
- **Task 29.3 — `done` 執行**：`runCLI()` 呼叫 `service.CompleteAndSave(command.ID)`，成功才寫 `"%d 已完成\n"`；錯誤交給 `runApp()`，使職責分離且容易用 `bytes.Buffer` 測試。成功／NotFound／SaveFailed 測試已涵蓋 Service、Fake Repository 狀態、錯誤型別、`t.Fatalf` 防止 nil 解參照；已完成整合與核心驗收。
- **Task 30.1 — Extract Function / `delete` 解析**：將 ID 轉換與正整數驗證提取為 `parsePositiveID(raw string)`，供 `done`、`delete` 共用；缺參數與參數過多仍由 `parseCommand()` 負責。學員選擇方案 A，理由是函式命名與職責一致，不應混合 CLI 參數數量驗證；新增成功與失敗 table-driven tests 並修正一次測試資料誤植，驗收通過。
- **Task 30.2 — `DeleteAndSave()`**：學員正確說明 `slices.Clone()` 是避免 `DeleteTodo()` 的 `append()` 修改相同 Backing Array；實作 Clone 候選狀態、`DeleteTodo(candidate,id)`、`repo.Save(slices.Clone(newTodos))`、成功才提交 `s.todos`，以 `%w` 包裝 Delete/Save 錯誤。測試成功、NotFound、SaveFailed、無 Repository，驗收通過。
- **Task 30.3 — `delete` CLI 執行**：`runCLI()` 呼叫 `service.DeleteAndSave(command.ID)`，錯誤直接回傳，成功才輸出 `"%d 刪除成功\n"`；新增成功、NotFound、SaveFailed 測試並驗收通過。學員選擇錯誤直接向上傳遞（方案 B），理解 `runApp()` 統一將錯誤寫入 stderr 並決定 Exit Code；不必重問。
- **Task 31.1 — `search` 參數解析**：加入 `parseCommand()` 的 `search` 分支，驗證恰好一個非空白關鍵字，保留原始字串，加入正常、保留空白、缺參數、空白、多餘參數的 table-driven tests。學員正確說明 Shell 引號與避免 `strings.Join()` 混淆 keyword／未來旗標的原因；驗收通過。
- **Task 31.2 — `search` CLI 執行**：`runCLI()` 使用 `Service.Search()` 和既有 `printTodos()`，查無結果輸出 `"找不到任何項目\n"`；不修改資料。初稿重複格式化與使用 `strings.Contains()` 容易誤判，修正為共用輸出函式與完整字串比較後驗收通過。已補充 Command Query Separation：查詢可讀取持久層，但不需寫入。
- **Task 31.3／31.3.1 — 搜尋整合測試與可測試性重構**：學員提出真實 `todos.json` 狀態不定及平行測試衝突，並正確說明 Wrapper 共用正式與測試邏輯。保留 `runApp()` 簽章，新增可注入路徑的 `runAppWithPath()`；使用 `t.TempDir()`、JSON Marshal／Unmarshal、`os.WriteFile`／`os.ReadFile` 測試真實 Repository 的搜尋成功／失敗、精確輸出、Exit Code 與資料內容不變。程式與測試均驗收通過。原始檔案位元組比較是可選建議，不影響驗收。

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
| 5    | CLI Application             | `os.Args`、`flag/FlagSet`、stdout／stderr、Exit Code、`add/list/done/delete/search`、Help Command                         | **進行中，下一步 Task 32.1** |
| 6    | 擴充 Todo Domain            | `time.Time`、Custom Type、`const/iota`、驗證、排序、Priority、Due Date、Created Time、Filter                                | 未開始                     |
| 7    | SQLite                      | `database/sql`、Query／QueryRow／Exec／Scan、SQLiteRepository                                                               | 未開始                     |
| 8    | SQL Transaction 與 Relation | Migration、Foreign Key、Transaction、1:N／N:M、Tags                                                                         | 未開始                     |
| 9    | REST API                    | `net/http`、Handler／ServeMux、JSON、HTTP Method／Status Code                                                               | 未開始                     |
| 10   | HTTP Testing 與 Middleware  | `httptest`、Logging、Request ID、Recovery、API Error Handling                                                               | 未開始                     |
| 11   | Goroutine／Channel          | Goroutine、Channel、Select、WaitGroup、Mutex、Context、Cancellation、Graceful Shutdown                                      | 未開始                     |
| 12   | 工程化                      | `cmd/`、`internal/`、Package Design、fmt／vet／test／race、Config、Structured Logging、Makefile、Dockerfile、GitHub Actions | 未開始                     |

---

## 七、跨次續課的重要提醒

- **唯一待辦起點**：Task 32.1 — CLI Help Command；先問「`todos.json` 損壞時 `--help` 是否仍應成功？應在哪一層處理以避免 Repository 初始化？」**此思考題尚未回答。**不要跳回 Task 31，也不要跳往 SQLite／REST API。
- **最新檔案優先**：2026-11-20 已上傳的 `main.go`、`main_test.go` 比原始 11/18 ZIP 更新；若新對話沒有最新檔案，要求學員上傳或使用本次重新整合的專案 ZIP，再進行 Code Review。不得根據歷史結論捏造本次測試結果。
- **歷史決策不要重問**：Task 29 Command Struct 的取捨、Task 30 `parsePositiveID`／Clone／錯誤傳遞、Task 31 Shell 參數契約、`printTodos()` 共用、CQS、`t.TempDir()` 隔離、`runAppWithPath()` Wrapper 與向後相容性均已學過。
- **Task 31 全部完成**：`search` CLI 的解析、執行、測試及整合驗證已通過，不再標記為待辦；`--help` **尚未實作**。
- **非阻擋改善不得變成新關卡**：包括搜尋前後 JSON 原始位元組比對、早期 NotFound stdout 補充斷言等。除非新提交實際發現問題，不得要求重做已驗收工作。
- **教學節奏**：從 Task 32.1 思考題開始，評估學員提出的架構方案，再以小步驟讓學員自己實作 help，給出 Code Review；不要直接交付整個方案。
