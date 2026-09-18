我正在使用一個 `todo-cli` 專案學習 Go。請把我上傳的專案當成整門課的核心教材，先閱讀目前專案內容，再依照現有程式碼循序漸進教我 Go。

我的目標不是只學語法，而是透過同一個 Todo 專案逐步學會 Go 的實務開發方式，最後把它從簡單 CLI 演進成具有測試、Service、Repository、SQLite、REST API、Concurrency 與工程化結構的小型 Go Backend 專案。

目前預計的課程主線如下：

1. 完整單元測試 + Slice 深入

   * testing
   * table-driven test
   * edge cases
   * array / slice
   * len / cap
   * backing array
   * append
   * slice mutation

2. Go Error Handling

   * error
   * errors.Is
   * errors.New
   * fmt.Errorf("%w")
   * 正確處理檔案不存在、JSON 損壞、IO error

3. Struct 與 Method

   * method
   * pointer receiver
   * value receiver
   * constructor pattern
   * zero value
   * 將 Todo function 演進為 Service

4. Interface 與 Dependency Injection

   * interface
   * implicit implementation
   * small interface
   * dependency injection
   * Repository Pattern
   * TodoRepository
   * JSONRepository

5. CLI Application

   * os.Args
   * flag / flag.FlagSet
   * stdin / stdout
   * exit code
   * 改成：
     todo add "Learn Go"
     todo list
     todo done 1
     todo delete 1
     todo search "Go"

6. 擴充 Todo Domain

   * time.Time
   * custom type
   * const / iota
   * validation
   * sorting
   * priority
   * due date
   * created time
   * filtering

7. SQLite

   * database/sql
   * Query / QueryRow / Exec
   * Scan
   * JSONRepository → SQLiteRepository

8. SQL Transaction 與 Relation

   * transaction
   * migration
   * foreign key
   * one-to-many / many-to-many
   * tags

9. REST API

   * net/http
   * Handler
   * ServeMux
   * JSON request / response
   * HTTP method
   * status code
   * REST API

10. HTTP Testing 與 Middleware

* httptest
* middleware
* logging
* request ID
* recovery
* API error handling

11. Goroutine / Channel

* goroutine
* channel
* select
* sync.WaitGroup
* sync.Mutex
* context.Context
* cancellation
* graceful shutdown

12. Go Project 工程化

* cmd/
* internal/
* package design
* go fmt
* go vet
* go test
* go test -race
* config
* environment variables
* structured logging
* Makefile
* Dockerfile
* GitHub Actions

教學方式請遵守以下原則：

* 每一堂課都直接基於我目前的專案修改，不要另外建立與 Todo 無關的玩具範例。
* 先檢查我目前的程式碼，不要假設專案狀態。
* 每次只進行一個合理大小的學習階段，不要一次把整堂課的最終答案全部寫給我。
* 教學順序以「問題 → 原因 → Go 概念 → 實作」為主。
* 在叫我改程式前，先說明為什麼這樣設計。
* 優先使用 Go standard library；例如 HTTP 先使用 `net/http`，資料庫先使用 `database/sql`，不要一開始就使用 Gin、Echo、GORM 等 framework。
* 不要只是告訴我答案，要讓我自己完成關鍵部分。
* 可以先給我任務、提示、函式介面或部分 skeleton code，讓我實作後再幫我 review。
* 如果我的寫法可以執行但不符合 Go 慣例，請明確指出 idiomatic Go 的寫法與原因。
* Code review 時請特別注意：

  * correctness
  * error handling
  * ownership / mutation
  * pointer 使用
  * package responsibility
  * naming
  * testability
  * Go idioms
* 如果某個概念在目前專案還沒有自然的使用情境，不要為了教而硬塞進來。
* 對 pointer、interface、slice、goroutine、channel 等比較容易誤解的概念，要結合目前專案實際資料流說明。
* 每完成一個階段，請簡短整理：

  1. 我這次學到了什麼
  2. 專案結構發生什麼變化
  3. 下一步會解決什麼問題

目前如果我沒有另外說明進度，請從：

「第 1 堂：完整單元測試 + Slice 深入」

開始。

第一步先閱讀我上傳的專案，告訴我：

1. 目前專案結構
2. 已經使用到哪些 Go 概念
3. 第 1 堂預計會修改哪些地方
4. 第一個要交給我實作的小任務

不要直接把第 1 堂所有完成版程式碼一次給我。
