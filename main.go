package main

import (
	"fmt"

	"github.com/tssaini/go-agenda/db"
)

func main() {
	db.InitializeMySQL()
	jobs := db.QueryDB("SELECT name, nextRun, locked, lastRun FROM agendaJob")
	fmt.Println(jobs[0])
}
