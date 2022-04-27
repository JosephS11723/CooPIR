package unzip

import (
	"archive/zip"
	"log"

	"github.com/JosephS11723/CooPIR/src/jobWorker/config"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/resultTypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/worker"
)

// Unzip attempts to unzip a file and upload its artifacts to seaweed
func AuthLog(job *dbtypes.Job, resultChan chan worker.ResultContainer, returnChan chan string) {
	// get information
	caseUUID := job.CaseUUID
	fileUUID := job.Files[0]

	// open reader for
	r, err := zip.OpenReader(config.WorkDir + "/" + caseUUID + "/" + fileUUID)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Close()

	for _, f := range r.File {
		filename := f.Name

		log.Println("Unzipping file: ", filename)

		if f.FileInfo().IsDir() {
			//TODO: handle directories
			//filename = filepath.Join(sourceFilename, filename)
			// Make Folder
			//os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// attempt to "open" the file
		rc, err := f.Open()
		if err != nil {
			log.Println("Could not open file: ", err)
		}

		// create result struct
		jobResult := worker.JobResult{
			ResultType: resultTypes.CreateFile,
			JobUUID:    job.JobUUID,
			CaseUUID:   job.CaseUUID,
			Tags:       []string{"uncompressed"},
			Name:       filename,
			Done:       false,
			FileUUID:   filename,
		}

		// create container
		container := worker.ResultContainer{
			JobResult:  jobResult,
			FileReader: rc,
		}

		// send to result channel
		resultChan <- container

		// output chan
		submittedFileuuid := <-returnChan

		// modify parent to have relation to child
		// create job result for moditying zip file "contains"
		jobResult = worker.JobResult{
			ResultType: resultTypes.ModifyFile,
			JobUUID:    job.JobUUID,
			CaseUUID:   job.CaseUUID,
			Tags:       []string{"compressed"},
			Name:       filename,
			Relations:  []string{submittedFileuuid + ":contains"},
			Done:       false,
			FileUUID:   fileUUID,
		}

		// create container
		container = worker.ResultContainer{
			JobResult:  jobResult,
			FileReader: nil,
		}

		// send to result channel
		resultChan <- container

		// void output
		<-returnChan

		// Close the file without defer to close before next iteration of loop
		rc.Close()
	}

	// send empty update and close
	/*jobResult := worker.JobResult{
		ResultType: resultTypes.ModifyFile,
		JobUUID:    job.JobUUID,
		CaseUUID:   job.CaseUUID,
		Tags:       []string{},
		Name:       "",
		Relations:  []string{},
		Done:       true,
		FileUUID:   fileUUID,
	}

	// create container
	container := worker.ResultContainer{
		JobResult:  jobResult,
		FileReader: nil,
	}

	// send to result channel
	resultChan <- container

	// void output
	<-returnChan*/
}
