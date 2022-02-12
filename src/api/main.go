package main

import (
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/jobqueue"
	"github.com/JosephS11723/CooPIR/src/api/routers"
)

func main() {
	// set log to print line numbers
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// initialize router with handlers
	r, queue_chan := routers.InitRouter()

	go jobqueue.ManageQueue(queue_chan)

	// run and serve
	r.Run("0.0.0.0:8080")
}
