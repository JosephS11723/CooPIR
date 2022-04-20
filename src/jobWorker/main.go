package main

import (
	// import worker lib

	"os"
	"os/signal"
	"syscall"

	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/worker"
)

func ctrlCExit() {
	// return if user presses ctrl+c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func main() {
	// create an instance of the worker
	worker := worker.NewJobWorker()

	// add a job that can be done

	// start the worker
	worker.Start()

	// stop worker if user presses ctrl+c
	ctrlCExit()

	// stop the worker
	worker.Stop()
}
