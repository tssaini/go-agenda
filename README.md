# go-agenda
Persistent job scheduling for Golang.

[![Go Report Card](https://goreportcard.com/badge/github.com/tssaini/go-agenda)](https://goreportcard.com/report/github.com/tssaini/go-agenda)

## Usage

```go
db, _ := sql.Open("mysql", "user:test@/goagenda")

agenda, _ := agenda.New(db)
agenda.Define("print hello", func() error {
    fmt.Println("Hello world")
    return nil
})

agenda.RepeatEvery("print hello", "* * * * *")
agenda.Start()
```

See cmd/main.go for a working example.

