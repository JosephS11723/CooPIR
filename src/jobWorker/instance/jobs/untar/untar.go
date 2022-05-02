package untar

import (
	"archive/tar"
	"errors"
	"io"
	"log"
	"os"

	"github.com/JosephS11723/CooPIR/src/jobWorker/config"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/resultTypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/worker"
)

// Unzip attempts to unzip a file and upload its artifacts to seaweed
func Untar(job *dbtypes.Job, resultChan chan worker.ResultContainer, returnChan chan string) error {
	// get information
	caseUUID := job.CaseUUID
	fileUUID := job.Files[0]

	// open file
	reader, err := os.Open(config.WorkDir + "/" + caseUUID + "/" + fileUUID)
	if err != nil {
		log.Println("Error opening file:", err)
		return errors.New("error opening file")
	}
	defer reader.Close()

	// create tarreader
	tarReader := tar.NewReader(reader)

	// for each file in the tar
	for {
		f, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return err
		}

		filename := f.Name

		log.Println("Untar-ing file: ", filename)

		if f.FileInfo().IsDir() {
			//TODO: handle directories
			//filename = filepath.Join(sourceFilename, filename)
			// Make Folder
			//os.MkdirAll(fpath, os.ModePerm)
			continue
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
			FileReader: tarReader,
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
			Tags:       []string{"compressed", "tar"},
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
	}

	return nil
}
