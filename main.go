package main

import (
	"fmt"
)

func main() {
	agenda := New()

	agenda.Define("print hello", func() error {
		fmt.Println("Hello world")
		return nil
	})

	agenda.Now("print hello")

	agenda.RepeatEvery("print hello", "0 8 * * *")
}
