package main

import (
	"fmt"

	"github.com/tssaini/go-agenda/job"
)

func main() {
	job.Define("print hello", func() error {
		fmt.Println("Hello world")
		return nil
	})

	job.Now("print hello")
}
