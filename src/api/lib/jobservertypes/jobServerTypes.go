package jobservertypes

/*
type NewJobQueue struct {
	NewJobs []NewJob
}

//push and pop functions
func (nj *NewJobQueue) push(job_request NewJob) {

	nj.NewJobs = append(nj.NewJobs, job_request)

}

func (nj *NewJobQueue) pop() NewJob {

	var item NewJob

	item, nj.NewJobs = nj.NewJobs[0], nj.NewJobs[1:]

	return item
}

type Workers struct {
	Info []WorkerInfo
}*/

type NewJob struct {
	ID      int32
	Case    string
	JobType string
	File    string
}

type WorkerInfo struct {
	Name           string
	Host           string
	Responsibility string
}
