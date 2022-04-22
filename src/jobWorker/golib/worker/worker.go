package worker

// import when library is public
/*import (
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
)*/

// temporary local import
import (
	"log"
	"sync"

	"github.com/JosephS11723/CooPIR/src/jobWorker/golib/dbtypes"
)

// JobWorker is the struct that wraps the basic JobWorker methods and variables.
type JobWorker struct {
	// JobQueue is the queue of jobs to be done
	JobQueue chan dbtypes.Job

	// JobResultQueue is the queue of job results to be submitted to the api.
	JobResultQueue chan JobResult

	// MaxWorkers is the maximum number of workers routines to be spawned
	MaxWorkers int `json:"maxworkers default:2"`

	// curJobs is the current number of in-progress jobs
	curJobs int `default:"0"`

	// curJobsLock is the mutex for curJobs
	curJobsLock sync.Mutex

	// jobList is the map of jobstrings that can be performed by this worker by an associated function
	jobList map[string]JobFunc

	// wg is the waitgroup for all the worker routines
	wg sync.WaitGroup

	// done is a bool that is set to true when the worker is set to exit
	done bool `default:"false"`

	// doneChan is a channel that allows blocked workers to unblock obtaining a job
	doneChan chan bool
}

// JobResult is the struct that wraps the information sent to the api
type JobResult struct {
	// JobUUID is the uuid of the job
	JobUUID string `json:"jobuuid"`

	// resultType is the type of result
	ResultType string `json:"resulttype"`

	// tags are the tags associated with the result
	Tags []string `json:"tags"`

	// Relations are the relations associated with the result
	Relations []string `json:"relations"`

	// name is the name of the result
	Name string `json:"name"`

	// done is a bool that indicates if the result is the final result
	Done bool `json:"done"`
}

type JobFunc func(*dbtypes.Job, chan JobResult)

// NewJobWorker creates a new JobWorker
func NewJobWorker(maxWorkers int) *JobWorker {
	// create a new JobWorker
	j := JobWorker{}

	// create the job queue
	j.JobQueue = make(chan dbtypes.Job)

	// create the job result queue
	j.JobResultQueue = make(chan JobResult, 100)

	// create the job list
	j.jobList = make(map[string]JobFunc)

	// make the waitgroup
	j.wg = sync.WaitGroup{}

	// set the max workers
	j.MaxWorkers = maxWorkers

	// return the new JobWorker
	return &j
}

// getWorkLoop is a routine that gets jobs from the api and queues them
func (j *JobWorker) getWorkLoop() {
	for {
		// if curJobs is less than max workers
		if j.curJobs < j.MaxWorkers {
			// ask server for work
			job, err := j.GetJob()

			// if there is an error, log it and continue
			if err != nil {
				log.Println("Error getting job:", err)
				continue
			}

			// increment curJobs
			j.curJobsLock.Lock()
			j.curJobs++
			j.curJobsLock.Unlock()

			// add the job to the queue
			j.JobQueue <- job
		}
	}
}

// submitWorkLoop is a routine that submits job results to the api
func (j *JobWorker) submitWorkLoop() {
	for {
		// wait for a job result
		jobResult := <-j.JobResultQueue

		// submit the job result
		err := j.SubmitJob(jobResult)

		// if there is an error, log it and continue
		if err != nil {
			log.Println("Error submitting job result:", err)
			// add job back into the finished queue
			j.JobResultQueue <- jobResult
			continue
		}

		// decrement curJobs
		j.curJobsLock.Lock()
		j.curJobs--
		j.curJobsLock.Unlock()
	}
}

// Start starts the JobWorker
func (j *JobWorker) Start() {
	go func() {
		// spawn workers
		for i := 0; i < j.MaxWorkers; i++ {
			// add a worker to the waitgroup
			j.wg.Add(1)

			// spawn a worker
			go j.workerRoutine()
		}

		// spawn the get work loop
		go j.getWorkLoop()

		// spawn the submit work loop
		go j.submitWorkLoop()
	}()
}

// Stop stops the JobWorker.
func (j *JobWorker) Stop() error {
	// set the done bool to true for all the workers
	j.done = true

	// wait for the routines and the waitgroup to finish
	j.wg.Wait()

	return nil
}

// SubmitJob submits a job result to the api
func (j *JobWorker) SubmitJob(result JobResult) error {
	return nil
}

// GetJob gets a job from the api
func (j *JobWorker) GetJob() (dbtypes.Job, error) {
	return dbtypes.Job{}, nil
}

// AddJobWithFunction adds a function to the list of possible jobs this worker can perform
func (j *JobWorker) AddJobWithFunction(jobTypeName string, jobFunction JobFunc) {
	// add job to job list
	j.jobList[jobTypeName] = jobFunction
}

// workerRoutine is the main routine spawned by the worker struct that creates the "workers" who take in jobs from the queue and perform them
func (j *JobWorker) workerRoutine() {
	// infinite loop
	for {
		// select between getting a job or exiting
		select {
		case job := <-j.JobQueue:
			// if the job is not in the list of jobs this worker can perform, log it and continue
			if _, ok := j.jobList[job.JobType]; !ok {
				log.Println("[ERROR workerRoutine]: Job type not supported:", job.JobType)
				continue
			}

			// perform the job
			j.jobList[job.JobType](&job, j.JobResultQueue)

		case <-j.doneChan:
			// decrement the waitgroup
			j.wg.Done()

			// exit
			return
		}
	}
}
