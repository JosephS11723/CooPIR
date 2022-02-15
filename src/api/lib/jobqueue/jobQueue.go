package jobqueue

import (
	"github.com/JosephS11723/CooPIR/src/api/lib/jobservertypes"
)

//has a 'queue' for NewJobs and WorkerInfo items
type JobQueue struct {
	//NewJobs jobservertypes.NewJobQueue
	//Workers jobservertypes.Workers
	NewJobs []jobservertypes.NewJob
	Workers []jobservertypes.WorkerInfo
}

//push and pop functions
func (jq JobQueue) job_push(job_request jobservertypes.NewJob) {

	jq.NewJobs = append(jq.NewJobs, job_request)

}

func (jq JobQueue) job_pop() jobservertypes.NewJob {

	var item jobservertypes.NewJob

	item, jq.NewJobs = jq.NewJobs[0], jq.NewJobs[1:]

	return item
}

//push and pop functions
func (jq JobQueue) workers_push(job_request jobservertypes.WorkerInfo) {

	jq.Workers = append(jq.Workers, job_request)

}

func (jq JobQueue) workers_pop() jobservertypes.WorkerInfo {

	var item jobservertypes.WorkerInfo

	item, jq.Workers = jq.Workers[0], jq.Workers[1:]

	return item
}

/*'main' function for the queue;
this should get passed to a goroutine spawned in main*/
func ManageQueue(job_receiver <-chan jobservertypes.NewJob, worker_receiver <-chan jobservertypes.WorkerInfo) {

	queue := JobQueue{}

	queue.NewJobs = make([]jobservertypes.NewJob, 0)

	queue.Workers = make([]jobservertypes.WorkerInfo, 0)

	for {

		select {

		case data := <-job_receiver:
			queue.job_push(data)

		default:

		}

		select {

		case data := <-worker_receiver:
			queue.workers_push(data)

		default:

		}

	}

}
