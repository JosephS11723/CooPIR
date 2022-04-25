package main

import (
	// import worker lib

	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/worker"
	"github.com/JosephS11723/CooPIR/src/jobWorker/instance/jobs/detectMime"
	"github.com/JosephS11723/CooPIR/src/jobWorker/instance/jobs/unzip"
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
	worker.AddJobWithFunction("Determine-MimeType", detectMime.DetermineMimeType)
	worker.AddJobWithFunction("Unzip", unzip.Unzip)

	// start the worker
	worker.Start()

	// print press ctrl+c to exit
	fmt.Println("Press ctrl+c to exit")

	// stop worker if user presses ctrl+c
	ctrlCExit()

	// stop the worker
	worker.Stop()
}
