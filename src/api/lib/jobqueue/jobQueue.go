package jobqueue

type JobQueue struct {
	NewJobs []NewJob
}

func (jq JobQueue) push(job_request NewJob) {

	jq.NewJobs = append(jq.NewJobs, job_request)

}

func (jq JobQueue) pop() {

}

func ManageQueue(receiver <-chan interface{}) {

}
