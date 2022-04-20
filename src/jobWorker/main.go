package main

import (
	// import worker lib

	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/JosephS11723/CooPIR/src/jobWorker/instance/jobs"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/worker"
)

func ctrlCExit() {
	// return if user presses ctrl+c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func main() {
	// create an instance of the worker for each cpu core
	worker := worker.NewJobWorker(runtime.NumCPU())

	// add a job that can be done
	worker.AddJobWithFunction("Determine-MimeType", jobs.DetermineMimeType)

	// start the worker
	worker.Start()

	// stop worker if user presses ctrl+c
	ctrlCExit()

	// stop the worker
	worker.Stop()
}
