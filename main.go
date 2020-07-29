package main

import (
	"fmt"
	"time"
)

func main() {
	agenda := New()

	agenda.Define("print hello", func() error {
		fmt.Println("Hello world")
		return nil
	})

	// agenda.Now("print hello")

	//every hour
	agenda.RepeatEvery("print hello", "15 * * * *")
	agenda.Start()
	time.Sleep(1 * time.Second)
	agenda.Define("print bad", func() error {
		fmt.Println("BAD")
		return nil
	})
	agenda.RepeatEvery("print bad", "* * * * *")
	time.Sleep(5000 * time.Minute)
}
