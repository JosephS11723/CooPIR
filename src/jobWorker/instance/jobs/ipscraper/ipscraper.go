package ipscraper

import (
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/worker"
)

func IPScraper(job *dbtypes.Job, resultChan chan worker.ResultContainer, returnChan chan string) {

	caseUUID := job.CaseUUID
	fileUUIDs := job.Files

	for _, file := range fileUUIDs {

	}

}
