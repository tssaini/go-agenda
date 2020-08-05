package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tssaini/go-agenda"
)

func main() {
	db, err := sql.Open("mysql", "user:test@/goagenda")
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	agenda, err := agenda.New(db)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	agenda.Define("print hello", func() error {
		fmt.Println("Hello world")
		return nil
	})

	// agenda.Now("print hello")

	//every hour
	agenda.RepeatEvery("print hello", "@every 2m")
	agenda.Start()
	agenda.Define("print bad", func() error {
		fmt.Println("BAD")
		time.Sleep(10 * time.Second)
		return errors.New("Unable to run Print bad")
	})
	agenda.RepeatEvery("print bad", "* * * * *")
	time.Sleep(5000 * time.Minute)
}
