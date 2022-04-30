package regexssn

import (
	"bytes"
	"errors"
	"log"
	"os"
	reg "regexp"

	"github.com/JosephS11723/CooPIR/src/jobWorker/config"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/regex"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/resultTypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/worker"
	"github.com/asaskevich/govalidator"
	"github.com/mingrammer/commonregex"
)

var regExpression *reg.Regexp = commonregex.SSNRegex

// ParseURLs parses the urls in a file
func RegexSSN(job *dbtypes.Job, resultChan chan worker.ResultContainer, returnChan chan string) error {
	// get information
	caseUUID := job.CaseUUID
	fileUUID := job.Files[0]

	// get reader for file
	file, err := os.OpenFile(config.WorkDir+"/"+caseUUID+"/"+fileUUID, os.O_RDONLY, 0755)

	if err != nil {
		log.Println("Error opening file:", err)
		return errors.New("could not open file")
	}

	defer file.Close()

	// get the urls from the file
	urls := regex.FindByReader(regExpression, file, 1024)

	log.Println(urls)

	for _, url := range urls {
		filename := url

		// assert is url with library
		if !govalidator.IsURL(url) {
			continue
		}

		// convert url to bytes
		urlBytes := []byte(url)

		// create reader for bytes
		reader := bytes.NewReader(urlBytes)

		log.Println("Uploading url: ", filename)

		// create result struct
		jobResult := worker.JobResult{
			ResultType: resultTypes.CreateFile,
			JobUUID:    job.JobUUID,
			CaseUUID:   job.CaseUUID,
			Tags:       []string{"url"},
			Name:       filename,
			Done:       false,
			FileUUID:   filename,
		}

		// create container
		container := worker.ResultContainer{
			JobResult:  jobResult,
			FileReader: reader,
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
			Tags:       []string{},
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
