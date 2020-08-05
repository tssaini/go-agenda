package main

import (
	"errors"
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
		time.Sleep(10 * time.Second)
		return errors.New("Unable to run Print bad")
	})
	agenda.RepeatEvery("print bad", "* * * * *")
	time.Sleep(5000 * time.Minute)
}
