package worker

// import when library is public
/*import (
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
)*/

// temporary local import
import (
	"github.com/JosephS11723/CooPIR/src/jobWorker/lib/dbtypes"
)

// JobWorker is the struct that wraps the basic JobWorker methods and variables.
type JobWorker struct {
	// JobQueue is the queue of jobs to be processed.
	JobQueue chan dbtypes.Job
	// JobResultQueue is the queue of job results to be submitted to the api.
	JobResultQueue chan dbtypes.JobResult
	//
}

// Start starts the JobWorker
func (j *JobWorker) Start() error

// Stop stops the JobWorker.
func (j *JobWorker) Stop() error

// SubmitJob submits a job result to the api
func (j *JobWorker) SubmitJob(job dbtypes.Job) error

// GetJob gets a job from the api
func (j *JobWorker) GetJob() (dbtypes.Job, error)
