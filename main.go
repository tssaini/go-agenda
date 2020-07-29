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
	agenda.RepeatEvery("print hello", "59 * * * *")
	agenda.Start()
	time.Sleep(5000 * time.Minute)
}
