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
	r, job_channel, worker_channel := routers.InitRouter()

	go jobqueue.ManageQueue(job_channel, worker_channel)

	// run and serve
	r.Run("0.0.0.0:8080")
}
