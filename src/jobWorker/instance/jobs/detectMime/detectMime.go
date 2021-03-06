package detectMime

import (
	"bufio"
	"errors"
	"log"
	"os"

	"github.com/JosephS11723/CooPIR/src/jobWorker/config"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/resultTypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/worker"
	"github.com/gabriel-vasile/mimetype"
)

// DetermineMimeType determines the mime type of a file
func DetermineMimeType(job *dbtypes.Job, resultChan chan worker.ResultContainer, returnChan chan string) error {
	// create directory for the job
	//os.Mkdir(job.JobUUID, 0755)

	// defer deleting the folder with everything in it (cleanup)
	//defer os.RemoveAll(job.JobUUID)

	// get relevant information from job
	caseUUID := job.CaseUUID
	fileUUID := job.Files[0]

	// get the file from seaweed by mounting the file
	// create mount
	//mount := seaweed.CreateSWMount("./"+job.JobUUID+"/"+fileUUID, "/"+caseUUID+"/"+fileUUID)

	// mount
	//err := mount.Mount()

	// check for error
	//if err != nil {
	//	log.Println("Error mounting file: ", err)
	//}

	// defer closing the mount
	//defer mount.Unmount()

	// JOB
	// get the mime type of the file
	// get reader for file
	file, err := os.OpenFile(config.WorkDir+"/"+caseUUID+"/"+fileUUID, os.O_RDONLY, 0755)

	if err != nil {
		log.Println("Error opening file:", err)
		return errors.New("error opening file")
	}

	defer file.Close()

	// get the mime type
	mimeType, err := getMimeTypeFromReader(bufio.NewReader(file))

	for err != nil {
		log.Println("Error getting mime type:", err)
		return errors.New("error getting mime type")
	}

	// create job result
	jobResult := worker.JobResult{
		JobUUID:    job.JobUUID,
		ResultType: resultTypes.ModifyFile,
		Tags:       []string{"mimetype:" + mimeType},
		Relations:  []string{},
		Name:       job.Name,
		Done:       false,
		FileUUID:   fileUUID,
		CaseUUID:   caseUUID,
	}

	// add to job container and send job result to job result queue
	resultChan <- worker.ResultContainer{
		JobResult:  jobResult,
		FileReader: nil,
	}

	// void return value
	<-returnChan

	return nil
}

func getMimeTypeFromReader(reader *bufio.Reader) (string, error) {
	// peek the first 3072 bytes from the reader
	bytes, _ := reader.Peek(3072)

	// get the mime type
	mimeType := mimetype.Detect(bytes)

	return mimeType.String(), nil
}
