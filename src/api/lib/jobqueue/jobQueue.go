package jobqueue

import (
	"github.com/JosephS11723/CooPIR/src/api/lib/jobservertypes"
)

type JobQueue struct {
	NewJobs []jobservertypes.NewJob
	Workers []WorkerInfo
}

type WorkerInfo struct {
	Name           string
	Host           string
	Responsibility string
}

func (jq JobQueue) push(job_request jobservertypes.NewJob) {

	jq.NewJobs = append(jq.NewJobs, job_request)

}

func (jq JobQueue) pop() jobservertypes.NewJob {

	var item jobservertypes.NewJob

	item, jq.NewJobs = jq.NewJobs[0], jq.NewJobs[1:]

	return item
}

func ManageQueue(receiver <-chan interface{}) {

}
