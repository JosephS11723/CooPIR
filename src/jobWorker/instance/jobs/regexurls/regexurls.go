package regexurls

import (
	"bytes"
	"log"
	"os"
	reg "regexp"
	"time"

	"github.com/JosephS11723/CooPIR/src/jobWorker/config"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/regex"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/resultTypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/worker"
	"github.com/asaskevich/govalidator"
)

var urlRegex = reg.MustCompile("(?:\"|')(((?:[a-zA-Z]{1,10}://|//)[^\"'/]{1,}\\.[a-zA-Z]{2,}[^\"']{0,})|((?:/|\\.\\./|\\./)[^\"'><,;| *()(%%$^/\\\\[\\]][^\"'><,;|()]{1,})|([a-zA-Z0-9_\\-/]{1,}/[a-zA-Z0-9_\\-/]{1,}\\.(?:[a-zA-Z]{1,4}|action)(?:[\\?|/][^\"|']{0,}|))|([a-zA-Z0-9_\\-]{1,}\\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml|exe)(?:\\?[^\"|']{0,}|)))(?:\"|')")

// ParseURLs parses the urls in a file
func ParseUrls(job *dbtypes.Job, resultChan chan worker.ResultContainer, returnChan chan string) {
	// get information
	caseUUID := job.CaseUUID
	fileUUID := job.Files[0]

	// get reader for file
	file, err := os.OpenFile(config.WorkDir+"/"+caseUUID+"/"+fileUUID, os.O_RDONLY, 0755)

	for err != nil {
		// log error
		log.Println(err)

		// sleep for 5 seconds
		time.Sleep(time.Duration(5) * time.Second)

		// get reader for file
		file, err = os.OpenFile(config.WorkDir+"./"+caseUUID+"/"+fileUUID, os.O_RDONLY, 0755)
	}

	defer file.Close()

	// get the urls from the file
	urls := regex.FindByReader(urlRegex, file, 100)

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
