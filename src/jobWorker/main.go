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
	"github.com/JosephS11723/CooPIR/src/jobWorker/instance/jobs/regexurls"
	"github.com/JosephS11723/CooPIR/src/jobWorker/instance/jobs/untar"
	"github.com/JosephS11723/CooPIR/src/jobWorker/instance/jobs/unzip"
)

func ctrlCExit() {
	// return if user presses ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func main() {
	// create an instance of the worker for each cpu core
	worker := worker.NewJobWorker(runtime.NumCPU())

	// add a job that can be done
	worker.AddJobWithFunction("Determine-MimeType", detectMime.DetermineMimeType)
	worker.AddJobWithFunction("ParseURLs", regexurls.ParseUrls)
	worker.AddJobWithFunction("Unzip", unzip.Unzip)
	worker.AddJobWithFunction("Untar", untar.Untar)

	// start the worker
	worker.Start()

	// print press ctrl+c to exit
	fmt.Println("Press ctrl+c to exit")

	// stop worker if user presses ctrl+c
	ctrlCExit()

	// stop the worker
	worker.Stop()
}
