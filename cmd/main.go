package main

import (
	"fmt"
	"time"

	"github.com/tssaini/go-agenda"
)

func main() {
	agenda, err := agenda.New()
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
		return nil
	})
	agenda.RepeatEvery("print bad", "* * * * *")
	time.Sleep(5000 * time.Minute)
}
